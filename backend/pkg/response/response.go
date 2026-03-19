package response

import (
	"encoding/json"
	"net/http"

	"github.com/edwardsean/codesmart/backend/pkg/errors"
)

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(v) //take v, turn it into JSON, and write it directly into the HTTP response body (w)
}

func WriteError(w http.ResponseWriter, err error) {
	var status int

	if domainErr, ok := err.(errors.Error); ok {
		status = domainErr.StatusCode()
	} else {
		status = http.StatusInternalServerError
	}

	WriteJSON(w, status, map[string]string{"message": err.Error()})
}
