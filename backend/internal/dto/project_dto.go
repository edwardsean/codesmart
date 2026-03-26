package dto

import (
	"time"

	"github.com/edwardsean/codesmart/backend/internal/domain"
)

type CreateProjectPayload struct {
	Title       string             `json:"title"        validate:"required,min=1,max=255"`
	Description string             `json:"description"`
	Language    domain.Language    `json:"language"     validate:"required"` // "go", "python", "typescript"
	Mode        domain.ProjectMode `json:"mode"         validate:"required,oneof=help learn"`
	SourceType  domain.SourceType  `json:"source_type"  validate:"required,oneof=github scratch"`
	SourceURL   string             `json:"source_url"` // only required if source_type = "github"
}

type ProjectResponseDTO struct { //to get the project when click project
	ID                 int                  `json:"id"`
	UserID             int                  `json:"user_id"`
	Title              string               `json:"title"`
	Description        string               `json:"description"`
	Language           domain.Language      `json:"language"`
	Mode               domain.ProjectMode   `json:"mode"`
	Difficulty         domain.Difficulty    `json:"difficulty,omitempty"`
	SourceType         domain.SourceType    `json:"source_type"`
	SourceURL          string               `json:"source_url,omitempty"`
	TotalLevels        int                  `json:"total_levels"`
	CompletedLevels    int                  `json:"completed_levels"`
	ProgressPercentage int                  `json:"progress_percentage"`
	Status             domain.ProjectStatus `json:"status"`
	ContainerID        string               `json:"container_id"`
	VolumeName         string               `json:"volume_name"`
	CreatedAt          time.Time            `json:"created_at"`
	UpdatedAt          time.Time            `json:"updated_at"`
}

type ProjectListItemDTO struct { //for get projects in display (doesnt need description, source url, created at)
	ID                 int                  `json:"id"`
	Title              string               `json:"title"`
	Language           domain.Language      `json:"language"`
	Mode               domain.ProjectMode   `json:"mode"`
	Difficulty         domain.Difficulty    `json:"difficulty,omitempty"`
	SourceType         domain.SourceType    `json:"source_type"`
	TotalLevels        int                  `json:"total_levels"`
	CompletedLevels    int                  `json:"completed_levels"`
	ProgressPercentage int                  `json:"progress_percentage"`
	Status             domain.ProjectStatus `json:"status"`
	UpdatedAt          time.Time            `json:"updated_at"`
}

// conversion helper
func ToProjectResponseDTO(p *domain.Project) *ProjectResponseDTO {
	return &ProjectResponseDTO{
		ID:                 p.ID,
		UserID:             p.UserID,
		Title:              p.Title,
		Description:        p.Description,
		Language:           p.Language,
		Mode:               p.Mode,
		Difficulty:         p.Difficulty,
		SourceType:         p.SourceType,
		SourceURL:          p.SourceURL,
		TotalLevels:        p.TotalLevels,
		CompletedLevels:    p.CompletedLevels,
		ProgressPercentage: p.ProgressPercentage,
		Status:             p.Status,
		ContainerID:        p.ContainerID,
		VolumeName:         p.VolumeName,
		CreatedAt:          p.CreatedAt,
		UpdatedAt:          p.UpdatedAt,
	}
}

func ToProjectListItemDTO(p *domain.Project) *ProjectListItemDTO {
	return &ProjectListItemDTO{
		ID:                 p.ID,
		Title:              p.Title,
		Language:           p.Language,
		Mode:               p.Mode,
		Difficulty:         p.Difficulty,
		SourceType:         p.SourceType,
		TotalLevels:        p.TotalLevels,
		CompletedLevels:    p.CompletedLevels,
		ProgressPercentage: p.ProgressPercentage,
		Status:             p.Status,
		UpdatedAt:          p.UpdatedAt,
	}
}
