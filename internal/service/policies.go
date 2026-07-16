package service

import (
	"context"
	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"github.com/sakamoto-max/ratelimiter/internal/repository"
	"github.com/sakamoto-max/ratelimiter/internal/repository/cache"
)

type Policy struct {
	pg    *repository.Db
	cache *cache.Cache
}

func (p *Policy) GetPolicies(ctx context.Context, ownerName string) (*[]domain.Policy, error) {
	return p.pg.Policy.GetPolicies(ctx, ownerName)
}

func (p *Policy) GetPolicy(ctx context.Context, policy domain.Policy) (domain.Policy, error) {

	Policy, err := p.cache.Policy.GetPolicy(ctx, policy)
	if err != nil || Policy.ResourceName == "" {
		Policy, err = p.pg.Policy.GetPolicy(ctx, policy)
		if err != nil {
			return domain.Policy{}, err
		}

		p.cache.Policy.SetPolicy(ctx, domain.Policy{
			OwnerName:         Policy.OwnerName,
			ResourceName:      Policy.ResourceName,
			BucketSize:        Policy.BucketSize,
			IntervalInSeconds: Policy.IntervalInSeconds,
			RefillPerSecond:   Policy.RefillPerSecond,
			CreatedAt:         Policy.CreatedAt,
			UpdatedAt:         Policy.UpdatedAt,
		})
	}

	return Policy, nil
}

func (p *Policy) DeletePolicy(ctx context.Context, policy domain.Policy) error {
	return p.pg.Policy.DeletePolicy(ctx, policy)
}

func (p *Policy) NewPolicy(ctx context.Context,policy domain.Policy) (*domain.Policy, error) {
	return p.pg.Policy.AddPolicy(ctx, policy)
}
