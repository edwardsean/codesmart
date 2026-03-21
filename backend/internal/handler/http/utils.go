package http

import (
	"net/http"
	"strconv"

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
func parseFileIDParam(r *http.Request) (int, error) {
	vars := mux.Vars(r)
	return strconv.Atoi(vars["fileId"])
}
