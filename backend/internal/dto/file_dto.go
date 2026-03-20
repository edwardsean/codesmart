package dto

import (
	"strings"

	"github.com/edwardsean/codesmart/backend/internal/domain"
)

type FileNodeDTO struct {
	ID   int    `json:"id,omitempty"`
	Path string `json:"path"`
	Name string `json:"name"`
	Type string `json:"type"` // "file" or "dir"
}

type FileContentDTO struct {
	ID      int    `json:"id,omitempty"`
	Path    string `json:"path"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type CreateFilePayload struct {
	Path     string `json:"path"     validate:"required"`
	Content  string `json:"content"`
	Language string `json:"language"`
	IsDir    bool   `json:"is_dir"`
}

type UpdateFilePayload struct {
	Path    string `json:"path"`    // for rename
	Content string `json:"content"` // for save
}

func ToFileContentDTO(f *domain.ProjectFile) *FileContentDTO {
	return &FileContentDTO{
		ID:      f.ID,
		Path:    f.FilePath,
		Name:    getFileName(f.FilePath),
		Content: f.Content,
	}
}

func getFileName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
