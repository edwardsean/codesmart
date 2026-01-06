package utils

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

func ParseJson(body io.Reader, payload any) error {
	if body == nil {
		return errors.New("missing body for parsing")
	}

	return json.NewDecoder(body).Decode(payload)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v) //take v, turn it into JSON, and write it directly into the HTTP response body (w)
}

func WriteError(w http.ResponseWriter, status int, err error) {
	WriteJSON(w, status, map[string]string{"message": err.Error()})
}
