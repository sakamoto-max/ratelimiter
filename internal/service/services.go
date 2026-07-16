package service

import (
	"github.com/sakamoto-max/ratelimiter/internal/repository"
	"github.com/sakamoto-max/ratelimiter/internal/repository/cache"
)

type Service struct {
	Owner  *Owner
	Grpc   *Grpc
	Policy *Policy
	Token  *Token
}

func NewService(cache *cache.Cache, repo *repository.Db) *Service {
	return &Service{
		Owner:  &Owner{cache: cache, pg: repo},
		Grpc:   &Grpc{Cache: cache, Pg: repo},
		Policy: &Policy{pg: repo, cache: cache},
		Token:  &Token{pg: repo, cache: cache},
	}
}
