package http

import (
	"fmt"
	"log"
	"net/http"

	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/internal/handler/http/middleware"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/response"

	"github.com/gorilla/mux"
)

type GithubHandler struct {
	// repository_store domain.RepositoryStore
	githubService domain.GithubService
	userService   domain.UserService
}

func NewGithubHandler(githubService domain.GithubService, userService domain.UserService) *GithubHandler {
	return &GithubHandler{githubService: githubService, userService: userService}
}

func (h *GithubHandler) RegisterRoutes(router *mux.Router) {
	gitrouter := router.PathPrefix("/github").Subrouter()
	authMiddleware := middleware.WithJWTAuth(h.userService)
	gitrouter.HandleFunc("/debug", func(w http.ResponseWriter, r *http.Request) {
		// fmt.Println("handler", r.Header)
		// w.WriteHeader(http.StatusOK)
		// w.Write([]byte("OK"))
		// json.NewEncoder(w).Encode(map[string]string{"hello": "world"})
		fmt.Fprintln(w, "OK")
	})
	gitrouter.HandleFunc("/getRepositories", authMiddleware(h.handleGetRepositories)).Methods("GET")
	// router.HandleFunc("/getRepositories", func(w http.ResponseWriter, r *http.Request) {
	// 	token, err := r.Cookie("access_token")
	// 	if err != nil || token.Value == "" {
	// 		fmt.Fprintln(w, "no token")
	// 		return
	// 	}
	// 	fmt.Fprintf(w, "token: %v", token.Value)
	// })
}

func (h *GithubHandler) handleGetRepositories(w http.ResponseWriter, r *http.Request) {
	//get the github token
	user, ok := r.Context().Value(middleware.UserKey).(*domain.User)
	if !ok || user == nil {
		// response.WriteError(w, http.StatusUnauthorized, fmt.Errorf("no valid user in context"))
		// return errors.ErrInvalidCredentials
		response.WriteError(w, errors.ErrInvalidCredentials)
		return
	}

	var repositories *[]domain.Repository
	repositories, err := h.githubService.GetUserRepositories(user)

	if err != nil {
		response.WriteError(w, err)
		return
	}

	log.Printf("successful: %v", repositories)
	response.WriteJSON(w, http.StatusOK, map[string]any{"repositories": repositories})

}
