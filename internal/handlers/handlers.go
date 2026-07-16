package handlers

import (
	"github.com/sakamoto-max/ratelimiter/internal/service"
)

type Handler struct {
	Owner  *Owner
	Policy *Policy
	Token  *Token
	Health *Health
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{
		Owner:  &Owner{service: service.Owner},
		Health: &Health{},
		Policy: &Policy{service: service.Policy},
		Token:  &Token{service: service.Token},
	}
}
