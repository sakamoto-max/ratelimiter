package dto

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var ErrValidationFailed = errors.New("validation failed")

type fieldErr struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

type FieldErrs []fieldErr

func (f FieldErrs) BadRequestErr(w http.ResponseWriter) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(w).Encode(f)
}

func extractValidationErrors(err error) FieldErrs {

	var allFieldErrs FieldErrs

	validationErrs, ok := err.(validator.ValidationErrors)
	if !ok {
		return nil
	}

	for _, err := range validationErrs {

		field := strings.ToLower(err.Field())
		tag := err.Tag()

		allFieldErrs = append(allFieldErrs, fieldErr{
			Field:   field,
			Message: tag,
		})
	}

	return allFieldErrs
}
