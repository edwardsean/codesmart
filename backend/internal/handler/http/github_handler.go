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
	authService   service.AuthService
}

func NewGithubHandler(githubService service.GithubService, authService service.AuthService) *GithubHandler {
	return &GithubHandler{githubService: githubService, authService: authService}
}

func (h *GithubHandler) RegisterRoutes(router *mux.Router) {
	gitrouter := router.PathPrefix("/github").Subrouter()
	authMiddleware := middleware.WithJWTAuth(h.authService)
	gitrouter.HandleFunc("/repositories", authMiddleware(h.handleGetRepositories)).Methods("GET")
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
