package service

import (
	"context"
	"github.com/sakamoto-max/ratelimiter/internal/utils"
	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"github.com/sakamoto-max/ratelimiter/internal/repository"
	"github.com/sakamoto-max/ratelimiter/internal/repository/cache"
)

type Owner struct {
	cache *cache.Cache
	pg    *repository.Db
}

func (o *Owner) CreateOwner(ctx context.Context, owner domain.Owner) (*domain.Owner, error) {

	token, err := utils.GenerateToken(owner.Name)
	if err != nil {
		return nil, err
	}

	owner.Token = token

	encryptedPassword, err := utils.EncryptPassword(owner.Password)
	if err != nil {
		return nil, err
	}

	owner.Password = encryptedPassword

	return o.pg.Owner.NewOwner(ctx, owner)
}
