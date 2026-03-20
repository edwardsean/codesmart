package http

import (
	"log"
	"net/http"

	"github.com/edwardsean/codesmart/backend/internal/handler/http/middleware"
	"github.com/edwardsean/codesmart/backend/internal/service"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/response"

	"github.com/gorilla/mux"
)

type GithubHandler struct {
	// repository_store domain.RepositoryStore
	githubService service.GithubService
}

func NewGithubHandler(githubService service.GithubService) *GithubHandler {
	return &GithubHandler{githubService: githubService}
}

func (h *GithubHandler) RegisterRoutes(router *mux.Router, authService service.AuthService) {
	gitRouter := router.PathPrefix("/github").Subrouter()
	gitRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(middleware.WithJWTAuth(authService)(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		}))
	})
	gitRouter.HandleFunc("/repositories", h.handleGetRepositories).Methods("GET")
}

func (h *GithubHandler) handleGetRepositories(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.WriteError(w, errors.ErrInvalidCredentials)
		return
	}

	repositories, err := h.githubService.GetUserRepositories(r.Context(), user)

	if err != nil {
		log.Printf("failed to get github repositories: %v", err)
		response.WriteError(w, err)
		return
	}

	log.Printf("successful: %v", repositories)
	response.WriteJSON(w, http.StatusOK, map[string]any{"repositories": repositories})

}
