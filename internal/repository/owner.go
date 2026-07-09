package repository

import (
	"context"
	"fmt"
	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"time"

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
					token,
					owner_id
				)
			VALUES
				(
					@token,
					@ownerId
				)
	`

	_, err = trnx.Exec(ctx, query, pgx.NamedArgs{
		"token":   owner.Token,
		"ownerId": ownerId,
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
		Token:     owner.Token,
	}, nil
}

// func (o *Owner) GetOwner(ctx context.Context, ownerName string) (domain.Owner, error)   {}
// func (o *Owner) DeleteOwner(ctx context.Context, ownerName string) error                {}
// func (o *Owner) UpdateOwner(ctx context.Context, owner domain.Owner) error              {}
