package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	pgConn "github.com/jackc/pgx/v5/pgconn"
	"github.com/sakamoto-max/ratelimiter/internal/domain"
)

type Token struct {
	pg *pgxpool.Pool
}

func (t *Token) NewToken(ctx context.Context, token domain.Token) (*domain.Token, error) {
	query := `	
	    INSERT INTO
		TOKENS 
			(name, token, expires_at)
		VALUES 
			(@name, @token, @expires_at)
		RETURNING created_at
	`

	var createdAt time.Time
	err := t.pg.QueryRow(ctx, query, pgx.NamedArgs{
		"name":       token.Name,
		"token":      token.Token,
		"expires_at": token.ExpiresAt,
	}).Scan(&createdAt)

	if err != nil {
		var pgErr *pgConn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, fmt.Errorf("token with name %v already exists", token.Name)
		}
		return nil, fmt.Errorf("failed to insert token : %w", err)
	}

	return &domain.Token{
		Name:      token.Name,
		Token:     token.Token,
		ExpiresAt: token.ExpiresAt,
		OwnerName: token.OwnerName,
		CreatedAt: createdAt,
	}, nil

}

func (t *Token) GetToken(ctx context.Context, name string) (*domain.Token, error) {
	query := `
		SELECT 
			token, 
			expires_at,
			created_at,
			updated_at
		FROM 
			TOKENS
		WHERE 
			name = @name
	`

	var token string
	var expiresAt time.Time
	var createdAt time.Time
	var updatedAt time.Time

	err := t.pg.QueryRow(ctx, query, pgx.NamedArgs{
		"name": name,
	}).Scan(&token, &expiresAt, &createdAt, &updatedAt)
	if err != nil {

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("token with name : %v not found", name)
		}

		return nil, fmt.Errorf("failed to get token : %w", err)
	}

	return &domain.Token{
		Name:      name,
		Token:     token,
		ExpiresAt: expiresAt,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}, nil
}

func (t *Token) DeleteToken(ctx context.Context, name string) error {
	query := `
		DELETE FROM 
			TOKENS
		WHERE 
			name = @name
	`

	_, err := t.pg.Exec(ctx, query, pgx.NamedArgs{
		"name": name,
	})
	if err != nil {
		return fmt.Errorf("failed to delete token : %w", err)
	}

	return nil
}
