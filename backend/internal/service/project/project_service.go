package project

import (
	"context"
	"net/http"

	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/internal/dto"
	"github.com/edwardsean/codesmart/backend/internal/repository"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/validator"
)

type ProjectService struct {
	projectRepo repository.ProjectRepository
}

func NewProjectService(projectRepo repository.ProjectRepository) *ProjectService {
	return &ProjectService{projectRepo: projectRepo}
}

func (s *ProjectService) GetProjects(ctx context.Context, userID int) ([]dto.ProjectListItemDTO, error) {
	projects, err := s.projectRepo.GetProjectsByUserID(ctx, userID)
	if err != nil {
		return nil, errors.NewError("failed to get projects", http.StatusInternalServerError)
	}

	result := make([]dto.ProjectListItemDTO, len(projects))
	for i, p := range projects {
		result[i] = *dto.ToProjectListItemDTO(&p)
	}

	return result, nil
}

func (s *ProjectService) GetProjectByID(ctx context.Context, id int, userID int) (*dto.ProjectResponseDTO, error) {
	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return nil, errors.NewError("project not found", http.StatusNotFound)
	}

	// ownership check, user can only access their own projects
	if project.UserID != userID {
		return nil, errors.NewError("forbidden", http.StatusForbidden)
	}

	result := dto.ToProjectResponseDTO(project)
	return result, nil
}

func (s *ProjectService) CreateProject(ctx context.Context, userID int, payload dto.CreateProjectPayload) (*dto.ProjectResponseDTO, error) {
	if err := validator.Validate.Struct(payload); err != nil {
		return nil, errors.NewError(err.Error(), http.StatusBadRequest)
	}

	// extra validation, source_url required if github
	if payload.SourceType == domain.SourceGithub && payload.SourceURL == "" {
		return nil, errors.NewError("source_url is required for github projects", http.StatusBadRequest)
	}

	project := &domain.Project{
		UserID:      userID,
		Title:       payload.Title,
		Description: payload.Description,
		Language:    payload.Language,
		Mode:        payload.Mode,
		SourceType:  payload.SourceType,
		SourceURL:   payload.SourceURL,
		Status:      domain.StatusActive,
	}

	if err := s.projectRepo.CreateProject(ctx, project); err != nil {
		return nil, errors.NewError("failed to create project", http.StatusInternalServerError)
	}

	result := dto.ToProjectResponseDTO(project)
	return result, nil
}

func (s *ProjectService) DeleteProject(ctx context.Context, id int, userID int) error {
	//check if project is there
	project, err := s.projectRepo.GetProjectByID(ctx, id)
	if err != nil {
		return errors.NewError("project not found", http.StatusNotFound)
	}

	// ownership check, user can only delete their own projects
	if project.UserID != userID {
		return errors.NewError("forbidden", http.StatusForbidden)
	}

	if err := s.projectRepo.DeleteProject(ctx, id); err != nil {
		return errors.NewError("failed to delete project", http.StatusInternalServerError)
	}

	return nil
}
