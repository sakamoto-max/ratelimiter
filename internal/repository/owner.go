package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"github.com/sakamoto-max/ratelimiter/internal/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Owner struct {
	pg *pgxpool.Pool
}

func (o *Owner) NewOwner(ctx context.Context, owner domain.Owner) (*domain.Owner, error) {

	query := `
		INSERT INTO 
		OWNERS
			(
				name, 
				email, 
				password
			)
		VALUES
			(
				@name, 
				@email, 
				@password
			)
		RETURNING id, name, email, created_at
	`

	trnx, err := o.pg.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction : %w", err)
	}

	defer trnx.Rollback(ctx)

	var ownerId string
	var name string
	var email string
	var createdAt time.Time

	err = trnx.QueryRow(ctx, query, pgx.NamedArgs{
		"name":     owner.Name,
		"email":    owner.Email,
		"password": owner.Password,
	}).Scan(&ownerId, &name, &email, &createdAt)
	if err != nil {
		return nil, fmt.Errorf("failed to create owner : %w", err)
	}

	query = `
		INSERT INTO
			TOKENS
				(
					name,
					token,
					expires_at
				)
			VALUES
				(
					@name,
					@token,
					@expires_at					
				)
	`

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{
		"name":       "default",
		"token":      owner.RatelimiterDefaultToken,
		"expires_at": utils.DefaultExpiresAt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert token : %w", err)
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction : %w", err)
	}

	return &domain.Owner{
		Id:        ownerId,
		Name:      name,
		Email:     email,
		CreatedAt: createdAt,
		RatelimiterDefaultToken: owner.RatelimiterDefaultToken,
		HttpReqToken: owner.HttpReqToken,
	}, nil
}