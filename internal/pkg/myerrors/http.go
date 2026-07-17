package myerrors

import (
	"encoding/json"
	"log"
	"net/http"
)

type HttpErr struct{}

func (h *HttpErr) ErrorWriter(w http.ResponseWriter, err error) {

	unwrappedErr := UnWrapErr(err)

	code := unwrappedErr.GetHttpCode()

	if code == http.StatusInternalServerError {
		log.Printf("internal server error : %v", err.Error())
		unwrappedErr.message = "internal server error"
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{
		"error": unwrappedErr.message,
	})
}
