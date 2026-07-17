package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sakamoto-max/ratelimiter/internal/domain/dto"
	myErr "github.com/sakamoto-max/ratelimiter/internal/pkg/myerrors"
	"github.com/sakamoto-max/ratelimiter/internal/service"
)

type Owner struct {
	service *service.Owner
	myErr.HttpErr
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
		o.HttpErr.ErrorWriter(w, err)
		return
	}

	RespWriter(w, resp, http.StatusCreated)
}
