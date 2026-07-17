package cache

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/sakamoto-max/ratelimiter/internal/domain"
	myErr "github.com/sakamoto-max/ratelimiter/internal/pkg/myerrors"

	"github.com/redis/go-redis/v9"
)

type Policy struct {
	client *redis.Client
}

type PolicyIface interface {
	// needs owner_name, resource_name, bucket_size, refill_per_second, interval_in_seconds, created_at, updated_at
	SetPolicy(ctx context.Context, policy domain.Policy)
	// needs owner_name, resource_name
	GetPolicy(ctx context.Context, policy domain.Policy) (domain.Policy, error)
	// needs owner_name, resource_name
	DeletePolicy(ctx context.Context, policy domain.Policy) error
}

func (c *Policy) SetPolicy(ctx context.Context, policy domain.Policy) {

	mainKey := fmt.Sprintf("owner_name:%v:resource_name:%v", policy.OwnerName, policy.ResourceName)
	bucketSizeKey := "bucket_size"
	intervalInSecondsKey := "interval_in_seconds"
	refillPerSecondKey := "refill_per_second"
	createdAtKey := "created_at"
	updatedAtKey := "updated_at"

	c.client.HSet(ctx, mainKey,
		bucketSizeKey, policy.BucketSize,
		intervalInSecondsKey, policy.IntervalInSeconds,
		refillPerSecondKey, policy.RefillPerSecond,
		createdAtKey, policy.CreatedAt,
		updatedAtKey, policy.UpdatedAt,
	)
}

func (c *Policy) GetPolicy(ctx context.Context, policy domain.Policy) (domain.Policy, error) {

	mainKey := fmt.Sprintf("owner_name:%v:resource_name:%v", policy.OwnerName, policy.ResourceName)
	bucketSizeKey := "bucket_size"
	intervalInSecondsKey := "interval_in_seconds"
	refillPerSecondKey := "refill_per_second"
	createdAtKey := "created_at"
	updatedAtKey := "updated_at"

	cmd, err := c.client.HGetAll(ctx, mainKey).Result()
	if err != nil {
		return domain.Policy{}, myErr.WrapErr(fmt.Errorf("failed to get policy of resource %v and owner %v : %w", policy.ResourceName, policy.OwnerName, err), myErr.InternalServerErr)
	}

	if len(cmd) == 0 {
		return domain.Policy{}, nil
	}

	bucketSize, _ := strconv.Atoi(cmd[bucketSizeKey])
	intervalInSeconds, _ := strconv.Atoi(cmd[intervalInSecondsKey])
	refillPerSecond, _ := strconv.ParseFloat(cmd[refillPerSecondKey], 64)
	createdAt, _ := time.Parse(time.RFC3339, cmd[createdAtKey])
	updatedAt, _ := time.Parse(time.RFC3339, cmd[updatedAtKey])

	return domain.Policy{
		ResourceName:      policy.ResourceName,
		OwnerName:         policy.OwnerName,
		BucketSize:        bucketSize,
		IntervalInSeconds: intervalInSeconds,
		RefillPerSecond:   refillPerSecond,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}, nil
}

func (c *Policy) DeletePolicy(ctx context.Context, policy domain.Policy) error {
	mainKey := fmt.Sprintf("owner_name:%v:resource_name:%v", policy.OwnerName, policy.ResourceName)

	err := c.client.Del(ctx, mainKey).Err()
	if err != nil {
		return myErr.WrapErr(fmt.Errorf("failed to delete policy : %w", err), myErr.InternalServerErr)
	}

	return nil
}
