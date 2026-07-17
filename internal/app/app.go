package app

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/sakamoto-max/ratelimiter/internal/config"
	"github.com/sakamoto-max/ratelimiter/internal/database"
	"github.com/sakamoto-max/ratelimiter/internal/handlers"
	"github.com/sakamoto-max/ratelimiter/internal/interceptors"
	"github.com/sakamoto-max/ratelimiter/internal/pkg/jwt"
	"github.com/sakamoto-max/ratelimiter/internal/repository"
	"github.com/sakamoto-max/ratelimiter/internal/repository/cache"
	"github.com/sakamoto-max/ratelimiter/internal/routes"
	"github.com/sakamoto-max/ratelimiter/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
	rateLimiterPb "github.com/sakamoto-max/rate_limiter_proto/shared/rate_limiter"
	"google.golang.org/grpc"
)

type App struct {
	HttpPort    string
	redisClient *redis.Client
	pgPool      *pgxpool.Pool
	grpcServer  *grpc.Server
	httpServer  *http.Server
	service     *service.Service
	config      *config.Config
	Router      *chi.Mux
}

func New(config *config.Config) *App {

	log.Printf("creating server : STAGE : %v", config.STAGE)

	if config.STAGE != "local" {
		if err := database.Migrate(context.Background(), config); err != nil {
			log.Fatal(err)
		}
	}

	redisClient := database.NewRedisClient(config)
	log.Println("connected to redis")

	pg := database.NewPostgresPool(config)
	log.Println("connected to postgres")

	cache := cache.NewCache(redisClient)
	db := repository.NewDb(pg)

	reg := prometheus.NewRegistry()

	interceptors := interceptors.NewInterceptors(reg)

	service := service.NewService(cache, db)

	handler := handlers.NewHandler(service)

	router := routes.NewRouter(handler, reg)

	jwt.Init(config)

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			interceptors.RequestId.UnaryInterceptor,
			interceptors.Auth.UnaryInterceptor,
			interceptors.PromIncoming.UnaryInterceptor,
			interceptors.CanonicalLogger.UnaryInterceptor,
			interceptors.PromOutgoing.UnaryInterceptor,
		),
	)

	httpServer := &http.Server{
		Addr:    ":" + config.Http.Port,
		Handler: router,
	}

	rateLimiterPb.RegisterRateLimiterServer(grpcServer, handler.Grpc)

	return &App{
		HttpPort:    config.Http.Port,
		redisClient: redisClient,
		grpcServer:  grpcServer,
		httpServer:  httpServer,
		service:     service,
		config:      config,
		pgPool:      pg,
		Router:      router,
	}
}

func (a *App) Run() {

	go a.StartHttpServer()

	go a.StartGrpcServer()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	sig := <-sigChan
	log.Printf("received signal %v, shutting down", sig)

	a.Shutdown()
}

func (a *App) StartHttpServer() {
	log.Printf("http server has started on %v", a.HttpPort)
	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("failed to start http server : %v", err)
	}
}

func (a *App) StartGrpcServer() {
	lis, err := net.Listen("tcp", ":"+a.config.Grpc.Port)
	if err != nil {
		log.Fatalf("failed to listen to tcp : %v", err)
	}

	log.Printf("grpc server has started on %v", a.config.Grpc.Port)
	if err := a.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to start grpc server : %v", err)
	}
}

func (a *App) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	a.grpcServer.GracefulStop()
	log.Println("grpc server has stopped")

	a.httpServer.Shutdown(ctx)
	log.Println("http server has stopped")

	a.pgPool.Close()
	log.Println("postgres pool is closed")

	a.redisClient.Close()
	log.Println("redis client is closed")

	log.Println("server has shutdown")
}
