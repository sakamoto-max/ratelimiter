package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sakamoto-max/ratelimiter/internal/domain/dto"
	"github.com/sakamoto-max/ratelimiter/internal/middleware"
	myErr "github.com/sakamoto-max/ratelimiter/internal/pkg/myerrors"
	"github.com/sakamoto-max/ratelimiter/internal/service"
)

type Policy struct {
	service *service.Policy
	myErr.HttpErr
}

func (p *Policy) New(w http.ResponseWriter, r *http.Request) {
	var userInput dto.CreatePolicy
	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, err := userInput.Validate()
	if err != nil {
		validationErr.BadRequestErr(w)
		return
	}

	ownerName := middleware.GetOwnerName(r.Context())

	resp, err := p.service.NewPolicy(r.Context(), userInput.MapToPolicy(ownerName))
	if err != nil {

		p.HttpErr.ErrorWriter(w, err)
		return
	}

	RespWriter(w, resp, http.StatusCreated)
}

func (p *Policy) Get(w http.ResponseWriter, r *http.Request) {
	var userInput dto.GetPolicy

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, err := userInput.Validate()
	if err != nil {
		validationErr.BadRequestErr(w)
		return
	}

	ownerName := middleware.GetOwnerName(r.Context())

	resp, err := p.service.GetPolicy(r.Context(), userInput.MapToPolicy(ownerName))
	if err != nil {
		p.HttpErr.ErrorWriter(w, err)
		return
	}

	RespWriter(w, resp, http.StatusOK)
}

func (p *Policy) Patch(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)
}

func (p *Policy) Delete(w http.ResponseWriter, r *http.Request) {
	var userInput dto.GetPolicy

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, err := userInput.Validate()
	if err != nil {
		validationErr.BadRequestErr(w)
		return
	}

	ownerName := middleware.GetOwnerName(r.Context())

	err = p.service.DeletePolicy(r.Context(), userInput.MapToPolicy(ownerName))
	if err != nil {

		p.HttpErr.ErrorWriter(w, err)
		return
	}

	RespWriter(w, map[string]string{"message": "policy deleted"}, http.StatusNotFound)
}
