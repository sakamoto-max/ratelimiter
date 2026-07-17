package domain

import "time"

type UserBucket struct {
	Allowed                    int64
	TokensLeft                 float64
	LuaScriptExecutionTimeInMS int64
}

type Policy struct {
	Id                string    `db:"id" json:"id,omitempty"`
	ResourceName      string    `db:"resource" json:"resource,omitempty"`
	OwnerName         string    `db:"owner" json:"owner,omitempty"`
	BucketSize        int       `db:"bucket_capacity" json:"bucket_capacity,omitempty"`
	IntervalInSeconds int       `db:"time_in_seconds" json:"time_in_seconds,omitempty"`
	RefillPerSecond   float64   `db:"refill_rate_per_second" json:"refill_rate_per_second,omitempty"`
	CreatedAt         time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at,omitempty"`
}

type Owner struct {
	Id                      string    `db:"id" json:"id,omitempty"`
	Name                    string    `db:"name" json:"name,omitempty"`
	Email                   string    `db:"email" json:"email,omitempty"`
	Password                string    `db:"password" json:"password,omitempty"`
	CreatedAt               time.Time `db:"created_at" json:"created_at,omitempty"`
	UpdatedAt               time.Time `db:"updated_at" json:"updated_at,omitempty"`
	HttpReqToken            string    `json:"httpReqToken"`
	RatelimiterDefaultToken string    `json:"rateLimiterDefaultToken"`
}

type Token struct {
	Id        string    `db:"id" json:"id,omitempty"`
	Name      string    `db:"name" json:"name"`
	Token     string    `db:"token" json:"token"`
	OwnerName string    `json:"ownerName"`
	ExpiresAt time.Time `db:"expires_at" json:"expiresAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt,omitempty"`
}
