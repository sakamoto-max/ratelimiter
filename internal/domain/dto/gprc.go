package dto

import (
	"errors"

	rateLimiterPb "github.com/sakamoto-max/rate_limiter_proto/shared/rate_limiter"
	myErrs "github.com/sakamoto-max/ratelimiter/internal/pkg/myerrors"
)

type CheckRequest struct {
	ClientIp     string
	ResourceName string
	OwnerName    string
}

func ValidateGrpcReq(req *rateLimiterPb.CheckRequest) error {

	if req.ClientIp == "" {
		return myErrs.WrapErr(errors.New("client ip is required"), myErrs.BadRequestErr)
	}

	if req.ResourceName == "" {
		return myErrs.WrapErr(errors.New("resource name is required"), myErrs.BadRequestErr)
	}

	return nil
}

func PbtoLocal(req *rateLimiterPb.CheckRequest, ownerName string) CheckRequest {
	return CheckRequest{
		ClientIp:     req.ClientIp,
		ResourceName: req.ResourceName,
		OwnerName:    ownerName,
	}
}

type CheckResponse struct {
	Allowed        bool
	TokensLeft     int64
	BucketCapacity int64
	TryAfter       int64
}

func LocalToPb(data *CheckResponse) *rateLimiterPb.CheckResponse {
	return &rateLimiterPb.CheckResponse{
		Allowed:        data.Allowed,
		TokensLeft:     data.TokensLeft,
		BucketCapacity: data.BucketCapacity,
		TryAfter:       data.TryAfter,
	}
}
