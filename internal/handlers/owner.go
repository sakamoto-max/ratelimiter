package handlers

import (
	"encoding/json"
	"net/http"
	"github.com/sakamoto-max/ratelimiter/internal/domain/dto"
	"github.com/sakamoto-max/ratelimiter/internal/service"
)

type Owner struct {
	service *service.Owner
}

func (o *Owner) Register(w http.ResponseWriter, r *http.Request) {

	var userInput dto.Register

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErr, err := userInput.Validate()
	if err != nil {
		validationErr.BadRequestErr(w)
		return
	}

	resp, err := o.service.CreateOwner(r.Context(), userInput.MapToOwner())
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

