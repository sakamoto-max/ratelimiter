package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/sakamoto-max/ratelimiter/internal/domain"
)

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
