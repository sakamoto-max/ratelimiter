package cache

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/redis/go-redis/v9"
)

type Cache struct {
	Policy     PolicyIface
	UserBucket BucketIface
}

func NewCache(client *redis.Client, reg *prometheus.Registry) *Cache {
	return &Cache{
		Policy:     &Policy{Client: client},
		UserBucket: newBucket(client, reg),
	}
}
