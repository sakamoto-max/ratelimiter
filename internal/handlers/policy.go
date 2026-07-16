package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sakamoto-max/ratelimiter/internal/domain/dto"
	"github.com/sakamoto-max/ratelimiter/internal/middleware"
	"github.com/sakamoto-max/ratelimiter/internal/service"
)

type Policy struct {
	service *service.Policy
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
		resp := map[string]string{
			"error": err.Error(),
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
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
		resp := map[string]string{
			"error": err.Error(),
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(resp)
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
		resp := map[string]string{
			"error": err.Error(),
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
}
