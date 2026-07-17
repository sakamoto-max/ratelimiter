package dto

import (
	"github.com/go-playground/validator/v10"
	"github.com/sakamoto-max/ratelimiter/internal/domain"
)

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
