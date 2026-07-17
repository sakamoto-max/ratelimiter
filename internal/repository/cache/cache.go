package cache

import (
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Policy     PolicyIface
	UserBucket BucketIface
}

func NewCache(client *redis.Client) *Cache {
	return &Cache{
		Policy:     &Policy{client: client},
		UserBucket: &Bucket{client: client},
	}
}
