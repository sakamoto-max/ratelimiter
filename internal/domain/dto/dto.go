package dto

import (
	"encoding/json"
	"errors"
	"net/http"
	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"strings"

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
	OwnerName         string `json:"ownerName" validate:"required"`
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

func (p CreatePolicy) MapToPolicy() domain.Policy {
	return domain.Policy{
		OwnerName:         p.OwnerName,
		ResourceName:      p.Resource,
		BucketSize:        p.BucketSize,
		IntervalInSeconds: p.IntervalInSeconds,
	}
}

// ######################################################################

type GetPolicy struct {
	OwnerName    string `json:"ownerName" validate:"required"`
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

func (p GetPolicy) MapToPolicy() domain.Policy {
	return domain.Policy{
		OwnerName:    p.OwnerName,
		ResourceName: p.ResourceName,
	}
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
