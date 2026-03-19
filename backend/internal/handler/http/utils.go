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
