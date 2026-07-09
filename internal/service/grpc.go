package service

import (
	"context"
	"math"
	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"github.com/sakamoto-max/ratelimiter/internal/interceptors"
	"github.com/sakamoto-max/ratelimiter/internal/repository"
	"github.com/sakamoto-max/ratelimiter/internal/repository/cache"

	rateLimiterPb "github.com/sakamoto-max/rate_limiter_proto/shared/rate_limiter"
)

type Grpc struct {
	rateLimiterPb.UnimplementedRateLimiterServer
	Pg       *repository.Db
	Cache    *cache.Cache
}

func (s *Grpc) Check(ctx context.Context, req *rateLimiterPb.CheckRequest) (*rateLimiterPb.CheckResponse, error) {

	ownerName := interceptors.GetOwnerName(ctx)

	policy, err := s.Cache.Policy.GetPolicy(ctx, domain.Policy{OwnerName: ownerName, ResourceName: req.ResourceName})
	if err != nil {
		return nil, err
	}

	if policy.BucketSize == 0 {
		policy, err = s.Pg.Policy.GetPolicy(ctx, domain.Policy{OwnerName: ownerName, ResourceName: req.ResourceName})
		if err != nil {
			return nil, err
		}

		s.Cache.Policy.SetPolicy(ctx, domain.Policy{
			OwnerName:         ownerName,
			ResourceName:      req.ResourceName,
			BucketSize:        policy.BucketSize,
			IntervalInSeconds: policy.IntervalInSeconds,
			RefillPerSecond:   policy.RefillPerSecond,
			CreatedAt:         policy.CreatedAt,
			UpdatedAt:         policy.UpdatedAt,
		})
	}

	var allowed bool

	userBucket, err := s.Cache.UserBucket.CheckLimit(ctx, req.ClientIp, policy)
	if err != nil {
		return nil, err
	}

	if userBucket.Allowed == 1 {
		allowed = true
	}

	var tryAfter int64

	if allowed == false {
		tryAfter = calculateRetryAfter(userBucket.TokensLeft, policy)
	}

	return &rateLimiterPb.CheckResponse{
		Allowed:        allowed,
		TokensLeft:     int64(userBucket.TokensLeft),
		BucketCapacity: int64(policy.BucketSize),
		TryAfter:       tryAfter,
	}, nil
}

func calculateRetryAfter(tokensLeft float64, policy domain.Policy) int64 {

	tokensNeeded := 1 - tokensLeft

	timeNeeded := math.Ceil(float64(tokensNeeded) / policy.RefillPerSecond)

	return int64(timeNeeded)
}
