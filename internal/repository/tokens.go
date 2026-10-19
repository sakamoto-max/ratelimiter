package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Token struct {
	pg *pgxpool.Pool
}

// func (t *Token) NewToken(ctx context.Context, token domain.Token) (domain.Token, error) {}
// func (t *Token) GetToken(ctx context.Context, token string) (domain.Token, error)       {}
// func (t *Token) DeleteToken(ctx context.Context, token string) error                    {}
