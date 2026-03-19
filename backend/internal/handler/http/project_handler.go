package http

import (
	"net/http"

	"github.com/edwardsean/codesmart/backend/internal/dto"
	"github.com/edwardsean/codesmart/backend/internal/handler/http/middleware"
	"github.com/edwardsean/codesmart/backend/internal/service"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/response"
	"github.com/edwardsean/codesmart/backend/pkg/sanitize"
	"github.com/edwardsean/codesmart/backend/pkg/utils"
	"github.com/gorilla/mux"
)

type ProjectHandler struct {
	projectService service.ProjectService
}

func NewProjectHandler(projectService service.ProjectService) *ProjectHandler {
	return &ProjectHandler{projectService: projectService}
}

func (h *ProjectHandler) RegisterRoutes(router *mux.Router, authService service.AuthService) {
	projectRouter := router.PathPrefix("/projects").Subrouter()

	projectRouter.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(middleware.WithJWTAuth(authService)(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		}))
	})

	projectRouter.HandleFunc("", h.handleGetProjects).Methods("GET")
	projectRouter.HandleFunc("", h.handleCreateProject).Methods("POST")
	projectRouter.HandleFunc("/{id}", h.handleGetProjectByID).Methods("GET")
	projectRouter.HandleFunc("/{id}", h.handleDeleteProjectByID).Methods("DELETE")
}

func (h *ProjectHandler) handleGetProjects(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.WriteError(w, errors.ErrInvalidCredentials)
		return
	}

	projects, err := h.projectService.GetProjects(r.Context(), user.ID)
	if err != nil {
		response.WriteError(w, errors.NewError(err.Error(), http.StatusInternalServerError))
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{"projects": projects})
}

func (h *ProjectHandler) handleCreateProject(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.WriteError(w, errors.ErrInvalidCredentials)
		return
	}

	var payload dto.CreateProjectPayload
	if err := utils.ParseJson(r.Body, &payload); err != nil {
		response.WriteError(w, errors.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	// sanitize text fields only, not mode/source_type/language (controlled values)
	payload.Title = sanitize.SanitizeString(payload.Title)
	payload.Description = sanitize.SanitizeString(payload.Description)

	project, err := h.projectService.CreateProject(r.Context(), user.ID, payload)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, project)
}

func (h *ProjectHandler) handleGetProjectByID(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.WriteError(w, errors.ErrInvalidCredentials)
		return
	}

	id, err := parseIDParam(r)
	if err != nil {
		response.WriteError(w, errors.NewError("invalid project id", http.StatusBadRequest))
		return
	}

	project, err := h.projectService.GetProjectByID(r.Context(), id, user.ID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, project)
}

func (h *ProjectHandler) handleDeleteProjectByID(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.WriteError(w, errors.ErrInvalidCredentials)
		return
	}

	id, err := parseIDParam(r)
	if err != nil {
		response.WriteError(w, errors.NewError("invalid project id", http.StatusBadRequest))
		return
	}

	if err := h.projectService.DeleteProject(r.Context(), id, user.ID); err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, nil)
}
