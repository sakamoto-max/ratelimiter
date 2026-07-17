package service

import (
	"context"

	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"github.com/sakamoto-max/ratelimiter/internal/pkg/jwt"
	"github.com/sakamoto-max/ratelimiter/internal/repository"
	"github.com/sakamoto-max/ratelimiter/internal/repository/cache"
)

type Token struct {
	pg    *repository.Db
	cache *cache.Cache
}

func (t *Token) GetToken(ctx context.Context, name string) (*domain.Token, error) {
	return t.pg.Token.GetToken(ctx, name)
}

func (t *Token) CreateToken(ctx context.Context, token domain.Token) (*domain.Token, error) {

	newToken, err := jwt.GenerateToken(token.OwnerName, token.ExpiresAt)
	if err != nil {
		return nil, err
	}

	token.Token = newToken
	return t.pg.Token.NewToken(ctx, token)
}

func (t *Token) DeleteToken(ctx context.Context, name string) error {

	return t.pg.Token.DeleteToken(ctx, name)
}
