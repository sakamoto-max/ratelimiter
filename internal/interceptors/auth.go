package interceptors

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sakamoto-max/ratelimiter/internal/opentel"
	"github.com/sakamoto-max/ratelimiter/internal/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var (
	ErrTokenExpired     = errors.New("token is expired, get a new access token at /refresh")
	ErrTokenMalformed   = errors.New("token is malformed. please check the token again")
	ErrTokenInvalid     = errors.New("token is invalid")
	ErrTokenIsMissing   = errors.New("token is missing, please provide the token")
	ErrRefreshExpired   = errors.New("referesh token is expired, please login again")
	ErrSignatureInvalid = errors.New("token's signature is invalid")
)

type ctxKey string

var OwnerNameKey ctxKey = "ownername"

// var latencyKey ctxKey = "latency"

type Auth struct {
	TimeTaken prometheus.Histogram
}

func NewAuth(reg *prometheus.Registry) *Auth {
	return &Auth{
		TimeTaken: promauto.With(reg).NewHistogram(prometheus.HistogramOpts{
			Name: "time_taken_for_authentication_in_seconds",
			Help: "the total time taken for autheticaing",
			Buckets: []float64{
				0.0005,
				0.001,
				0.0025,
				0.005,
				0.01,
			},
		},
		),
	}
}



// gets the token from the metadata
// and authenticates it
// also calculates how much time it took to authenticate
func (a *Auth) UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

	ctx, span := opentel.CreateNewSpan(ctx, "auth_interceptor")
	defer span.End()

	timeStart := time.Now()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrTokenIsMissing
	}

	token := md.Get("token")

	if len(token) == 0 {
		return nil, fmt.Errorf("token is missing")
	}

	tokenStr := token[0]

	if r := recover(); r != nil {
		return nil, fmt.Errorf("please provide token")
	}

	claims, err := utils.ValidateToken(tokenStr)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token : %w", err)
	}

	a.TimeTaken.Observe(float64(time.Since(timeStart).Seconds()))

	newCtx := context.WithValue(ctx, OwnerNameKey, claims.Ownername)

	return handler(newCtx, req)
}

func GetOwnerName(ctx context.Context) string {
	val, ok := ctx.Value(OwnerNameKey).(string)
	if !ok {
		return ""
	}
	return val
}
