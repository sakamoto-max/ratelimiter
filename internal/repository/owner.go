package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	pgConn "github.com/jackc/pgx/v5/pgconn"
	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"github.com/sakamoto-max/ratelimiter/internal/pkg/jwt"
	myErrs "github.com/sakamoto-max/ratelimiter/internal/pkg/myerrors"

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
		return nil, myErrs.WrapErr(fmt.Errorf("failed to begin transaction : %w", err), myErrs.InternalServerErr)
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
		return nil, myErrs.WrapErr(fmt.Errorf("failed to create owner : %w", err), myErrs.InternalServerErr)
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
		"expires_at": jwt.DefaultExpiresAt,
	})
	if err != nil {
		var pgErr *pgConn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" && pgErr.ConstraintName == "tokens_name_key" {
			return nil, myErrs.WrapErr(fmt.Errorf("user already exists"), myErrs.AlreadyExistsErr)
		}
		return nil, myErrs.WrapErr(fmt.Errorf("failed to create default token : %w", err), myErrs.InternalServerErr)
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return nil, myErrs.WrapErr(fmt.Errorf("failed to commit transaction : %w", err), myErrs.InternalServerErr)
	}

	return &domain.Owner{
		Id:                      ownerId,
		Name:                    name,
		Email:                   email,
		CreatedAt:               createdAt,
		RatelimiterDefaultToken: owner.RatelimiterDefaultToken,
		HttpReqToken:            owner.HttpReqToken,
	}, nil
}
