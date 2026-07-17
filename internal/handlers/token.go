package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/sakamoto-max/ratelimiter/internal/domain/dto"
	"github.com/sakamoto-max/ratelimiter/internal/middleware"
	myErr "github.com/sakamoto-max/ratelimiter/internal/pkg/myerrors"
	"github.com/sakamoto-max/ratelimiter/internal/service"
)

type Token struct {
	service *service.Token
	myErr.HttpErr
}

func (t *Token) NewToken(w http.ResponseWriter, r *http.Request) {
	var userInput dto.NewToken

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, err := userInput.Validate()
	if err != nil {
		validationErrs.BadRequestErr(w)
		return
	}

	ownerName := middleware.GetOwnerName(r.Context())

	parsedUserInput, err := userInput.ParseNdMapToToken()
	if err != nil {
		t.HttpErr.ErrorWriter(w, err)
		return
	}
	parsedUserInput.OwnerName = ownerName

	token, err := t.service.CreateToken(r.Context(), parsedUserInput)
	if err != nil {
		t.HttpErr.ErrorWriter(w, err)
		return
	}

	RespWriter(w, token, http.StatusCreated)
}

func (t *Token) GetToken(w http.ResponseWriter, r *http.Request) {
	var userInput dto.TokenName

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, err := userInput.Validate()
	if err != nil {
		validationErrs.BadRequestErr(w)
		return
	}

	domainToken, err := t.service.GetToken(r.Context(), userInput.Name)
	if err != nil {
		t.HttpErr.ErrorWriter(w, err)
		return
	}

	RespWriter(w, domainToken, http.StatusOK)
}

func (t *Token) DeleteToken(w http.ResponseWriter, r *http.Request) {
	var userInput dto.TokenName

	json.NewDecoder(r.Body).Decode(&userInput)

	validationErrs, err := userInput.Validate()
	if err != nil {
		validationErrs.BadRequestErr(w)
		return
	}

	err = t.service.DeleteToken(r.Context(), userInput.Name)
	if err != nil {
		t.HttpErr.ErrorWriter(w, err)
		return
	}

	RespWriter(w, map[string]string{"message": "token deleted"}, http.StatusNoContent)
}
