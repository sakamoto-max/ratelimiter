package database

import (
	"context"
	"fmt"
	"log"
	"github.com/sakamoto-max/ratelimiter/internal/config"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

func NewRedisClient(config *config.Config) *redis.Client {

	redisUrl := fmt.Sprintf("redis://%v:%v@%v:%v/%v", config.Redis.UserName, config.Redis.Password, config.Redis.Host, config.Redis.Port, config.Redis.Db)

	opts, err := redis.ParseURL(redisUrl)
	if err != nil {
		log.Fatalf("failed to create redis client : %v", err)
	}

	client := redis.NewClient(opts)
	return client
}

func NewPostgresPool(config *config.Config) *pgxpool.Pool {
	postgresUrl := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v", config.Postgres.UserName, config.Postgres.Password, config.Postgres.Host, config.Postgres.Port, config.Postgres.Db, config.Postgres.SSLmode)

	pgxConfig, err := pgxpool.ParseConfig(postgresUrl)
	if err != nil {
		log.Fatalf("failed to parse postgres config : %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), pgxConfig)
	if err != nil {
		log.Fatalf("failed to create postgres pool : %v", err)
	}

	ctxForPing, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	err = pool.Ping(ctxForPing)
	if err != nil {
		log.Fatalf("failed to ping postgres : %v", err)
	}

	return pool
}
