package dto

import (
	"strings"

	"github.com/edwardsean/codesmart/backend/internal/domain"
)

type FileNodeDTO struct {
	Path string `json:"path"`
	Name string `json:"name"`
	Type string `json:"type"` // "file" or "dir"
}

type FileContentDTO struct {
	Path    string `json:"path"`
	Name    string `json:"name"`
	Content string `json:"content"`
}

type CreateFilePayload struct {
	Path     string          `json:"path"     validate:"required"`
	Content  string          `json:"content"`
	Language domain.Language `json:"language"`
	IsDir    bool            `json:"is_dir"`
}

type UpdateFilePayload struct {
	Path     string          `json:"path"`    // for rename
	Content  string          `json:"content"` // for save
	Language domain.Language `json:"language"`
}

type RenameFilePayload struct {
	OldPath string `json:"old_path" validate:"required"`
	NewPath string `json:"new_path" validate:"required"`
}

type FileTreeResult struct {
	Tree      []FileTreeItem `json:"tree"`
	Truncated bool           `json:"truncated"`
}

type FileTreeItem struct {
	Path string `json:"path"`
	Type string `json:"type"` // "blob" = file, "tree" = dir
	Size int    `json:"size"`
}

func ToFileContentDTO(f *domain.ProjectFile) *FileContentDTO {
	return &FileContentDTO{
		Path:    f.FilePath,
		Name:    getFileName(f.FilePath),
		Content: f.Content,
	}
}

func getFileName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
