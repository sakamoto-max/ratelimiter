package dto

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"github.com/sakamoto-max/ratelimiter/internal/utils"

	"github.com/go-playground/validator/v10"
)

var ErrValidationFailed = errors.New("validation failed")

type fieldErr struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type FieldErrs []fieldErr

func (f FieldErrs) BadRequestErr(w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(f)
}

// ######################################################################

type Register struct {
	Name     string `json:"name" validate:"required"`
	Email    string `json:"email" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (r *Register) Validate() (FieldErrs, error) {
	validate := validator.New()

	err := validate.Struct(r)
	if err != nil {
		return extractValidationErrors(err), ErrValidationFailed
	}

	return nil, nil
}

func (r Register) MapToOwner() domain.Owner {
	return domain.Owner{
		Name:     r.Name,
		Email:    r.Email,
		Password: r.Password,
	}
}

// ######################################################################

type CreatePolicy struct {
	Resource          string `json:"resource" validate:"required"`
	BucketSize        int    `json:"bucketSize" validate:"required"`
	IntervalInSeconds int    `json:"intervalInSeconds" validate:"required"`
}

func (p *CreatePolicy) Validate() (FieldErrs, error) {
	validate := validator.New()

	err := validate.Struct(p)
	if err != nil {
		return extractValidationErrors(err), ErrValidationFailed
	}

	return nil, nil
}

func (p CreatePolicy) MapToPolicy(ownerName string) domain.Policy {
	return domain.Policy{
		OwnerName:         ownerName,
		ResourceName:      p.Resource,
		BucketSize:        p.BucketSize,
		IntervalInSeconds: p.IntervalInSeconds,
	}
}

// ######################################################################

type GetPolicy struct {
	ResourceName string `json:"resourceName" validate:"required"`
}

func (p *GetPolicy) Validate() (FieldErrs, error) {
	validate := validator.New()

	err := validate.Struct(p)
	if err != nil {
		return extractValidationErrors(err), ErrValidationFailed
	}

	return nil, nil
}

func (p GetPolicy) MapToPolicy(ownerName string) domain.Policy {
	return domain.Policy{
		OwnerName:    ownerName,
		ResourceName: p.ResourceName,
	}
}

// ######################################################################

var (
	timeLayout = "2006-01-02 15:04:05"
)

type NewToken struct {
	Name       string `json:"name" validate:"required"`
	ExpireTime string `json:"expireTime"`
}

func (n *NewToken) Validate() (FieldErrs, error) {
	validate := validator.New()

	err := validate.Struct(n)
	if err != nil {
		return extractValidationErrors(err), ErrValidationFailed
	}

	return nil, nil
}

func (n *NewToken) ParseNdMapToToken() (domain.Token, error) {

	var ExpiresAt time.Time

	switch {
	case n.ExpireTime == "":
		ExpiresAt = utils.DefaultExpiresAt
	case n.ExpireTime != "":
		expiresAt, err := time.Parse(timeLayout, n.ExpireTime)
		if err != nil {
			return domain.Token{}, fmt.Errorf("failed to parse time :%w", err)
		}
		ExpiresAt = expiresAt
	}

	return domain.Token{
		Name:      n.Name,
		ExpiresAt: ExpiresAt,
	}, nil
	
}

// ######################################################################

type TokenName struct {
	Name string `json:"name" validate:"required"`
}

func (g *TokenName) Validate() (FieldErrs, error) {
	validate := validator.New()

	err := validate.Struct(g)
	if err != nil {
		return extractValidationErrors(err), ErrValidationFailed
	}

	return nil, nil
}

// ######################################################################

func extractValidationErrors(err error) FieldErrs {

	var allFieldErrs FieldErrs

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	for _, err := range validationErrs {

		field := strings.ToLower(err.Field())
		tag := err.Tag()

		allFieldErrs = append(allFieldErrs, fieldErr{
			Field:   field,
			Message: tag,
		})
	}

	return allFieldErrs
}
