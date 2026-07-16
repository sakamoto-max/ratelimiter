package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sakamoto-max/ratelimiter/internal/domain/dto"
	"github.com/sakamoto-max/ratelimiter/internal/middleware"
	"github.com/sakamoto-max/ratelimiter/internal/service"
)

type Token struct {
	service *service.Token
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
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		// todo
		return
	}
	parsedUserInput.OwnerName = ownerName

	token, err := t.service.CreateToken(r.Context(), parsedUserInput)
	if err != nil {

		resp := map[string]string{
			"error": err.Error(),
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(resp)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(token)
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
		// todo
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(domainToken)
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
		// todo
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
