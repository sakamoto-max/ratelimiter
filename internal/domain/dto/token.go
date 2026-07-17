package dto

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/sakamoto-max/ratelimiter/internal/domain"
	"github.com/sakamoto-max/ratelimiter/internal/pkg/jwt"
	myErrs "github.com/sakamoto-max/ratelimiter/internal/pkg/myerrors"
)

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
		ExpiresAt = jwt.DefaultExpiresAt
	case n.ExpireTime != "":
		expiresAt, err := time.Parse(timeLayout, n.ExpireTime)
		if err != nil {
			return domain.Token{}, myErrs.WrapErr(fmt.Errorf("failed to parse time :%w", err), myErrs.BadRequestErr)
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
