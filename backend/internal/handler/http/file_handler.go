package http

import (
	"log"
	"net/http"

	"github.com/edwardsean/codesmart/backend/internal/dto"
	"github.com/edwardsean/codesmart/backend/internal/handler/http/middleware"
	"github.com/edwardsean/codesmart/backend/internal/service"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/response"
	"github.com/edwardsean/codesmart/backend/pkg/utils"
	"github.com/gorilla/mux"
)

type FileHandler struct {
	fileService service.FileService
}

func NewFileHandler(fileService service.FileService) *FileHandler {
	return &FileHandler{fileService: fileService}
}

func (h *FileHandler) RegisterRoutes(router *mux.Router, authService service.AuthService) {
	r := router.PathPrefix("/projects/{projectId}/files").Subrouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(middleware.WithJWTAuth(authService)(func(w http.ResponseWriter, r *http.Request) {
			next.ServeHTTP(w, r)
		}))
	})

	r.HandleFunc("", h.handleGetFileTree).Methods("GET")
	r.HandleFunc("/content", h.handleGetFileContent).Methods("GET")
	r.HandleFunc("", h.handleCreateFile).Methods("POST")
	r.HandleFunc("/{fileId}", h.handleUpdateFile).Methods("PUT")
	r.HandleFunc("/{fileId}", h.handleDeleteFile).Methods("DELETE")
}

func (h *FileHandler) handleGetFileTree(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.WriteError(w, errors.ErrInvalidCredentials)
		return
	}

	projectId, err := parseProjectIDParam(r)
	if err != nil {
		response.WriteError(w, errors.NewError("invalid project id", http.StatusBadRequest))
		return
	}

	tree, err := h.fileService.GetFileTree(r.Context(), projectId, user.ID)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{"files": tree})

}
func (h *FileHandler) handleGetFileContent(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.WriteError(w, errors.ErrInvalidCredentials)
		return
	}

	projectId, err := parseProjectIDParam(r)
	if err != nil {
		response.WriteError(w, errors.NewError("invalid project id", http.StatusBadRequest))
		return
	}

	path := r.URL.Query().Get("path")
	if path == "" {
		response.WriteError(w, errors.NewError("path is required", http.StatusBadRequest))
		return
	}
	log.Printf("path: %s", path)

	content, err := h.fileService.GetFileContent(r.Context(), projectId, user.ID, path)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, map[string]any{"file": content})
}

func (h *FileHandler) handleCreateFile(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.WriteError(w, errors.ErrInvalidCredentials)
		return
	}

	projectId, err := parseProjectIDParam(r)
	if err != nil {
		response.WriteError(w, errors.NewError("invalid project id", http.StatusBadRequest))
		return
	}

	var payload dto.CreateFilePayload
	if err := utils.ParseJson(r.Body, &payload); err != nil {
		response.WriteError(w, errors.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	file, err := h.fileService.CreateFile(r.Context(), projectId, user.ID, payload)
	if err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusCreated, map[string]any{"file": file})
}

func (h *FileHandler) handleUpdateFile(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.WriteError(w, errors.ErrInvalidCredentials)
		return
	}

	projectId, err := parseProjectIDParam(r)
	if err != nil {
		response.WriteError(w, errors.NewError("invalid project id", http.StatusBadRequest))
		return
	}

	fileId, err := parseFileIDParam(r)
	if err != nil {
		response.WriteError(w, errors.NewError("invalid file id", http.StatusBadRequest))
		return
	}

	var payload dto.UpdateFilePayload
	if err := utils.ParseJson(r.Body, &payload); err != nil {
		response.WriteError(w, errors.NewError(err.Error(), http.StatusBadRequest))
		return
	}

	if err := h.fileService.UpdateFile(r.Context(), projectId, user.ID, fileId, payload); err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, nil)
}

func (h *FileHandler) handleDeleteFile(w http.ResponseWriter, r *http.Request) {
	user, err := middleware.GetUserFromContext(r)
	if err != nil {
		response.WriteError(w, errors.ErrInvalidCredentials)
		return
	}

	projectId, err := parseProjectIDParam(r)
	if err != nil {
		response.WriteError(w, errors.NewError("invalid project id", http.StatusBadRequest))
		return
	}

	fileId, err := parseFileIDParam(r)
	if err != nil {
		response.WriteError(w, errors.NewError("invalid file id", http.StatusBadRequest))
		return
	}

	if err := h.fileService.DeleteFile(r.Context(), projectId, user.ID, fileId); err != nil {
		response.WriteError(w, err)
		return
	}

	response.WriteJSON(w, http.StatusOK, nil)
}
