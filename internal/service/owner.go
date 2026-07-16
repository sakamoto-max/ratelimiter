package service

import (
	"context"

	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"github.com/sakamoto-max/ratelimiter/internal/repository"
	"github.com/sakamoto-max/ratelimiter/internal/repository/cache"
	"github.com/sakamoto-max/ratelimiter/internal/utils"
)

type Owner struct {
	cache *cache.Cache
	pg    *repository.Db
}

func (o *Owner) CreateOwner(ctx context.Context, owner domain.Owner) (*domain.Owner, error) {

	Httptoken, err := utils.GenerateToken(owner.Name, utils.DefaultExpiresAt)
	if err != nil {
		return nil, err
	}

	rlToken, err := utils.GenerateToken(owner.Name, utils.DefaultExpiresAt)
	if err != nil {
		return nil, err
	}

	owner.HttpReqToken = Httptoken
	owner.RatelimiterDefaultToken = rlToken

	encryptedPassword, err := utils.EncryptPassword(owner.Password)
	if err != nil {
		return nil, err
	}

	owner.Password = encryptedPassword

	return o.pg.Owner.NewOwner(ctx, owner)
}
