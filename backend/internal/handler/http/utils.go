package http

import (
	"net/http"
	"strconv"

	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/gorilla/mux"
)

func parseIDParam(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	return strconv.Atoi(vars["id"])
}

func parseProjectIDParam(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	return strconv.Atoi(vars["projectId"])
}

func parseFilePathParam(r *http.Request) (string, error) {
	vars := mux.Vars(r)
	path := vars["path"]
	if path == "" {
		return "", errors.NewError("path parameter is required", http.StatusBadRequest)
	}
	return path, nil
}
