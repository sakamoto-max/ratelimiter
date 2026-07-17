package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sakamoto-max/ratelimiter/internal/service"
)

type Handler struct {
	Grpc   *Grpc
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
		Grpc:   &Grpc{service: service.Grpc},
	}
}

func RespWriter(w http.ResponseWriter, resp any, code int) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(resp)
}
