package service

import (
	"context"
	"math"

	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"github.com/sakamoto-max/ratelimiter/internal/domain/dto"
	"github.com/sakamoto-max/ratelimiter/internal/repository"
	"github.com/sakamoto-max/ratelimiter/internal/repository/cache"
)

type Grpc struct {
	pg    *repository.Db
	cache *cache.Cache
}

func (s *Grpc) Check(ctx context.Context, req dto.CheckRequest) (*dto.CheckResponse, error) {

	policy, err := s.cache.Policy.GetPolicy(ctx, domain.Policy{OwnerName: req.OwnerName, ResourceName: req.ResourceName})
	if err != nil {
		return nil, err
	}

	if policy.BucketSize == 0 {
		policy, err = s.pg.Policy.GetPolicy(ctx, domain.Policy{OwnerName: req.OwnerName, ResourceName: req.ResourceName})
		if err != nil {
			return nil, err
		}

		s.cache.Policy.SetPolicy(ctx, domain.Policy{
			OwnerName:         req.OwnerName,
			ResourceName:      req.ResourceName,
			BucketSize:        policy.BucketSize,
			IntervalInSeconds: policy.IntervalInSeconds,
			RefillPerSecond:   policy.RefillPerSecond,
			CreatedAt:         policy.CreatedAt,
			UpdatedAt:         policy.UpdatedAt,
		})
	}

	var allowed bool

	userBucket, err := s.cache.UserBucket.CheckLimit(ctx, req.ClientIp, policy)
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

	return &dto.CheckResponse{
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
