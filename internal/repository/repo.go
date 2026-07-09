package repository

import (
	"context"
	"github.com/sakamoto-max/ratelimiter/internal/domain"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Db struct {
	Policy interface {
		GetPolicies(ctx context.Context, ownerName string) (*[]domain.Policy, error) 
		// needs owner_name, resource_name
		GetPolicy(ctx context.Context, data domain.Policy) (domain.Policy, error)
		DeletePolicy(ctx context.Context, policy domain.Policy) error
		AddPolicy(ctx context.Context, policy domain.Policy) (*domain.Policy, error)
	}
	Owner interface {
		NewOwner(ctx context.Context, owner domain.Owner) (*domain.Owner, error)
		// GetOwner(ctx context.Context, ownerName string) (domain.Owner, error) // todo
		// DeleteOwner(ctx context.Context, ownerName string) error // todo
		// UpdateOwner(ctx context.Context, owner domain.Owner) error // todo
	}
	Token interface { // todo 
		NewToken(ctx context.Context, token domain.Token) (domain.Token, error)
		GetToken(ctx context.Context, token string) (domain.Token, error)
		DeleteToken(ctx context.Context, token string) error
	}
}

func NewDb(pgxPool *pgxpool.Pool) *Db {
	return &Db{
		Policy: &Policy{pg: pgxPool},
		Owner:  &Owner{pg: pgxPool},
	}
}

