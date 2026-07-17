package handlers

import (
	"context"

	rateLimiterPb "github.com/sakamoto-max/rate_limiter_proto/shared/rate_limiter"
	"github.com/sakamoto-max/ratelimiter/internal/domain/dto"
	"github.com/sakamoto-max/ratelimiter/internal/interceptors"
	myErr "github.com/sakamoto-max/ratelimiter/internal/pkg/myerrors"
	"github.com/sakamoto-max/ratelimiter/internal/service"
)

type Grpc struct {
	service  *service.Grpc
	myErr.GrpcErr
	rateLimiterPb.UnimplementedRateLimiterServer
}

func NewGrpcHandler(service *service.Grpc, ) *Grpc {
	return &Grpc{service: service}
}

func (g *Grpc) Check(ctx context.Context, req *rateLimiterPb.CheckRequest) (*rateLimiterPb.CheckResponse, error) {

	ownerName := interceptors.GetOwnerName(ctx)

	err := dto.ValidateGrpcReq(req)
	if err != nil {
		return nil, g.GrpcErr.New(err)
	}

	data := dto.PbtoLocal(req, ownerName)

	res, err := g.service.Check(ctx, data)
	if err != nil {
		return nil, g.GrpcErr.New(err)
	}

	return dto.LocalToPb(res), nil
}
