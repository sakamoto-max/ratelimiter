package interceptors

import (
	"context"
	"errors"
	"fmt"
	"github.com/sakamoto-max/ratelimiter/internal/utils"
	"time"

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

var latencyKey ctxKey = "latency"

type Auth struct {}

func (a *Auth) UnaryInterceptor(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	timeStart := time.Now()

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, ErrTokenIsMissing
	}

	token := md.Get("token")[0]

	claims, err := utils.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("failed to validate token : %w", err)
	}

	newCtx := context.WithValue(ctx, OwnerNameKey, claims.Ownername)
	newCtx = context.WithValue(newCtx, latencyKey, time.Since(timeStart))

	return handler(newCtx, req)
}

func GetOwnerName(ctx context.Context) string {
	val, ok := ctx.Value(OwnerNameKey).(string)
	if !ok {
		return ""
	}
	return val
}

