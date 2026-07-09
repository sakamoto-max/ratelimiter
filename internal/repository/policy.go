package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Policy struct {
	pg *pgxpool.Pool
}

func (p *Policy) GetPolicies(ctx context.Context, ownerName string) (*[]domain.Policy, error) {
	query := `
		SELECT
			POLICIES.RESOURCE, 
			POLICIES.BUCKET_CAPACITY, 
			POLICIES.TIME_IN_SECONDS, 
			POLICIES.REFILL_RATE_PER_SECOND, 
			POLICIES.CREATED_AT, 
			POLICIES.UPDATED_AT 
		FROM 
			POLICIES
		INNER JOIN 
			OWNERS
		ON 
			POLICIES.OWNER_ID = OWNERS.ID
		WHERE 
			OWNERS.NAME = @ownerName
	`
	rows, err := p.pg.Query(ctx, query, pgx.NamedArgs{
		"ownerName": ownerName,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get policies : %w", err)
	}

	var allPolicies []domain.Policy

	var resource string
	var bucketCapacity int
	var timeInSeconds int
	var refillRatePerSecond float64
	var createdAt time.Time
	var updatedAt time.Time

	for rows.Next() {
		err := rows.Scan(&resource, &bucketCapacity, &timeInSeconds, &refillRatePerSecond, &createdAt, &updatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan policies : %w", err)
		}

		allPolicies = append(allPolicies, domain.Policy{
			ResourceName:      resource,
			OwnerName:         ownerName,
			BucketSize:        bucketCapacity,
			IntervalInSeconds: timeInSeconds,
			RefillPerSecond:   refillRatePerSecond,
			CreatedAt:         createdAt,
			UpdatedAt:         updatedAt,
		})
	}

	if len(allPolicies) == 0 {
		return nil, fmt.Errorf("no polices found")
	}

	return &allPolicies, nil
}

func (p *Policy) GetPolicy(ctx context.Context, data domain.Policy) (domain.Policy, error) {

	query := `
		SELECT
			BUCKET_CAPACITY, 
			TIME_IN_SECONDS, 
			REFILL_RATE_PER_SECOND, 
			POLICIES.CREATED_AT, 
			POLICIES.UPDATED_AT 
		FROM 
			POLICIES
		INNER JOIN 
			OWNERS
		ON 
			POLICIES.OWNER_ID = OWNERS.ID
		WHERE 
			OWNERS.NAME = @ownerName
		AND 
			POLICIES.RESOURCE = @resourceName
	`

	var bucketCapacity int
	var timeInSeconds int
	var refillRatePerSecond float64
	var createdAt time.Time
	var updatedAt time.Time

	err := p.pg.QueryRow(ctx, query, pgx.NamedArgs{
		"ownerName":    data.OwnerName,
		"resourceName": data.ResourceName,
	}).Scan(&bucketCapacity, &timeInSeconds, &refillRatePerSecond, &createdAt, &updatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return domain.Policy{}, fmt.Errorf("policy not found")
		}
		return domain.Policy{}, fmt.Errorf("failed to get policy : %w", err)
	}

	return domain.Policy{
		ResourceName:      data.ResourceName,
		OwnerName:         data.OwnerName,
		BucketSize:        bucketCapacity,
		IntervalInSeconds: timeInSeconds,
		RefillPerSecond:   refillRatePerSecond,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}, nil
}

func (p *Policy) DeletePolicy(ctx context.Context, policy domain.Policy) error {
	query := `
		DELETE FROM 
			POLICIES
		WHERE 
			OWNER_ID = (
				SELECT 
					ID 
				FROM 
					OWNERS
				WHERE 
					NAME = @ownerName
			)

		AND 
			RESOURCE = @resourceName
	`

	_, err := p.pg.Exec(ctx, query, pgx.NamedArgs{
		"ownerName":    policy.OwnerName,
		"resourceName": policy.ResourceName,
	})

	if err != nil {
		return fmt.Errorf("failed to delete policy : %w", err)
	}

	return nil
}

func (p *Policy) AddPolicy(ctx context.Context, policy domain.Policy) (*domain.Policy, error) {

	query := `
		SELECT 
			ID 
		FROM 
			OWNERS 
		WHERE 
			name = @name
	`

	trnx, err := p.pg.Begin(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction : %w", err)
	}

	defer trnx.Rollback(ctx)

	var ownerId string

	err = trnx.QueryRow(ctx, query, pgx.NamedArgs{
		"name": policy.OwnerName,
	}).Scan(&ownerId)
	if err != nil {
		return nil, fmt.Errorf("failed to get owner id : %w", err)
	}

	query = `
		INSERT INTO POLICIES (OWNER_ID, RESOURCE, BUCKET_CAPACITY, TIME_IN_SECONDS)
		VALUES (@ownerId, @resource, @bucketCapacity, @timeInSeconds)

		RETURNING REFILL_RATE_PER_SECOND, CREATED_AT, UPDATED_AT
	`

	var refillRatePerSecond float64
	var createdAt time.Time
	var updatedAt time.Time

	err = trnx.QueryRow(ctx, query, pgx.NamedArgs{
		"ownerId":        ownerId,
		"resource":       policy.ResourceName,
		"bucketCapacity": policy.BucketSize,
		"timeInSeconds":  policy.IntervalInSeconds,
	}).Scan(&refillRatePerSecond, &createdAt, &updatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to add policy : %w", err)
	}

	err = trnx.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to commit transaction : %w", err)
	}

	return &domain.Policy{
		OwnerName:         policy.OwnerName,
		ResourceName:      policy.ResourceName,
		BucketSize:        policy.BucketSize,
		IntervalInSeconds: policy.IntervalInSeconds,
		RefillPerSecond:   refillRatePerSecond,
		CreatedAt:         createdAt,
		UpdatedAt:         updatedAt,
	}, nil
}
