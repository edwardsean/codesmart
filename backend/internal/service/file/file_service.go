package file

import (
	"context"
	"net/http"
	"strings"

	"github.com/edwardsean/codesmart/backend/internal/clients"
	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/internal/dto"
	"github.com/edwardsean/codesmart/backend/internal/repository"
	"github.com/edwardsean/codesmart/backend/internal/service"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
)

type FileService struct {
	projectRepo  repository.ProjectRepository
	fileRepo     repository.ProjectFileRepository
	userRepo     repository.UserRepository
	githubClient clients.GithubClient
}

func NewFileService(projectRepo repository.ProjectRepository, fileRepo repository.ProjectFileRepository, userRepo repository.UserRepository, githubClient clients.GithubClient) *FileService {
	return &FileService{
		projectRepo:  projectRepo,
		fileRepo:     fileRepo,
		userRepo:     userRepo,
		githubClient: githubClient,
	}
}

func (s *FileService) GetFileTree(ctx context.Context, projectId int, userId int) ([]dto.FileNodeDTO, error) {
	project, err := s.projectRepo.GetProjectByID(ctx, projectId)
	if err != nil {
		return nil, errors.NewError("project not found", http.StatusNotFound)
	}

	//ownership
	if project.UserID != userId {
		return nil, errors.NewError("forbidden", http.StatusForbidden)
	}

	//since already seeded
	// if project.SourceType == domain.SourceGithub {
	// 	return s.getGithubFileTree(ctx, project)
	// }

	return s.getDBFileTree(ctx, projectId)
}

func (s *FileService) GetFileContent(ctx context.Context, projectId int, userId int, path string) (*dto.FileContentDTO, error) {
	project, err := s.projectRepo.GetProjectByID(ctx, projectId)
	if err != nil {
		return nil, errors.NewError("project not found", http.StatusNotFound)
	}

	//ownership
	if project.UserID != userId {
		return nil, errors.NewError("forbidden", http.StatusForbidden)
	}

	//since already seeded
	// if project.SourceType == domain.SourceGithub {
	// 	return s.getGithubFileContent(ctx, project, path)
	// }

	return s.getDBFileContent(ctx, projectId, path)
}

func (s *FileService) CreateFile(ctx context.Context, projectId int, userId int, payload dto.CreateFilePayload) (*dto.FileContentDTO, error) {
	project, err := s.projectRepo.GetProjectByID(ctx, projectId)
	if err != nil {
		return nil, errors.NewError("project not found", http.StatusNotFound)
	}

	//ownership
	if project.UserID != userId {
		return nil, errors.NewError("forbidden", http.StatusForbidden)
	}

	file := &domain.ProjectFile{
		ProjectID: projectId,
		FilePath:  payload.Path,
		Content:   payload.Content,
		Language:  domain.Language(payload.Language),
	}

	if err := s.fileRepo.CreateFile(ctx, file); err != nil {
		return nil, errors.NewError("failed to create file", http.StatusInternalServerError)
	}

	return dto.ToFileContentDTO(file), nil
}

func (s *FileService) UpdateFile(ctx context.Context, projectId, userId, fileId int, payload dto.UpdateFilePayload) error {
	project, err := s.projectRepo.GetProjectByID(ctx, projectId)
	if err != nil {
		return errors.NewError("project not found", http.StatusNotFound)
	}

	//ownership
	if project.UserID != userId {
		return errors.NewError("forbidden", http.StatusForbidden)
	}

	existing, err := s.fileRepo.GetFileByID(ctx, fileId)
	if err != nil {
		return errors.NewError("file not found", http.StatusNotFound)
	}

	if existing.ProjectID != projectId {
		return errors.NewError("file does not belong to this project", http.StatusForbidden)
	}

	if payload.Path != "" {
		existing.FilePath = payload.Path
		if payload.Language == "" {
			existing.Language = domain.Language(service.DetectLanguage(payload.Path))
		}
	}

	if payload.Content != "" {
		existing.Content = payload.Content
	}

	return s.fileRepo.UpdateFile(ctx, existing)

}

func (s *FileService) DeleteFile(ctx context.Context, projectId, userId, fileId int) error {
	project, err := s.projectRepo.GetProjectByID(ctx, projectId)
	if err != nil {
		return errors.NewError("project not found", http.StatusNotFound)
	}

	//ownership
	if project.UserID != userId {
		return errors.NewError("forbidden", http.StatusForbidden)
	}

	file, err := s.fileRepo.GetFileByID(ctx, fileId)
	if err != nil {
		return errors.NewError("file not found", http.StatusNotFound)
	}

	if file.ProjectID != projectId {
		return errors.NewError("file does not belong to this project", http.StatusForbidden)
	}

	return s.fileRepo.DeleteFile(ctx, fileId)
}

// func (s *FileService) getGithubFileContent(ctx context.Context, project *domain.Project, path string) (*dto.FileContentDTO, error) {
// 	token, err := s.decryptGithubToken(ctx, project)
// 	if err != nil {
// 		return nil, err
// 	}

// 	owner, repo, err := service.ParseGithubURL(project.SourceURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// use GitHub tree API — recursive=1 gets the full tree in one call
// 	decoded, err := s.githubClient.GetFileContent(ctx, token, owner, repo, path)

// 	if err != nil {
// 		return nil, errors.NewError("failed to decode file content", http.StatusInternalServerError)
// 	}

// 	return &dto.FileContentDTO{
// 		Path:    path,
// 		Name:    getFileName(path),
// 		Content: string(decoded),
// 	}, nil
// }

func (s *FileService) getDBFileContent(ctx context.Context, projectId int, path string) (*dto.FileContentDTO, error) {
	file, err := s.fileRepo.GetFileContent(ctx, projectId, path)
	if err != nil {
		return nil, errors.NewError("failed to get files", http.StatusInternalServerError)
	}

	return dto.ToFileContentDTO(file), nil
}

// func (s *FileService) getGithubFileTree(ctx context.Context, project *domain.Project) ([]dto.FileNodeDTO, error) {
// 	token, err := s.decryptGithubToken(ctx, project)
// 	if err != nil {
// 		return nil, err
// 	}

// 	owner, repo, err := service.ParseGithubURL(project.SourceURL)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// use GitHub tree API — recursive=1 gets the full tree in one call
// 	result, err := s.githubClient.GetFileTree(ctx, token, owner, repo)

// 	if err != nil {
// 		return nil, errors.NewError(err.Error(), http.StatusInternalServerError)
// 	}

// 	nodes := make([]dto.FileNodeDTO, 0, len(result.Tree)) //length = 0, capacity = number of files (so go doesnt need to allocate space when exceed to grow capacity)
// 	for i, item := range result.Tree {
// 		nodeType := "file"
// 		if item.Type == "tree" {
// 			nodeType = "dir"
// 		}
// 		nodes = append(nodes, dto.FileNodeDTO{ //append because length is 0 at first
// 			ID:   i + 1, //REVISE
// 			Path: item.Path,
// 			Name: getFileName(item.Path),
// 			Type: nodeType,
// 		})
// 	}

// 	return nodes, nil
// }

func (s *FileService) getDBFileTree(ctx context.Context, projectId int) ([]dto.FileNodeDTO, error) {
	files, err := s.fileRepo.GetFiles(ctx, projectId)
	if err != nil {
		return nil, errors.NewError("failed to get files", http.StatusInternalServerError)
	}

	nodes := make([]dto.FileNodeDTO, len(files))
	for i, file := range files {
		nodes[i] = dto.FileNodeDTO{
			ID:   file.ID,
			Path: file.FilePath,
			Name: getFileName(file.FilePath),
			Type: "file", //REVISE
		}
	}

	return nodes, nil
}

// func (s *FileService) decryptGithubToken(ctx context.Context, project *domain.Project) (string, error) {
// 	user, err := s.userRepo.GetUserByID(ctx, project.UserID)
// 	if err != nil {
// 		return "", errors.NewError("user not found", http.StatusNotFound)
// 	}

// 	if user.GithubToken == "" {
// 		return "", errors.NewError("GitHub not connected, connect your GitHub account first", http.StatusBadRequest)
// 	}

// 	key, err := base64.StdEncoding.DecodeString(config.Envs.EncryptionKey)
// 	if err != nil {
// 		return "", errors.NewError("encryption config error", http.StatusInternalServerError)
// 	}

// 	return password.Decrypt(user.GithubToken, key)
// }

func getFileName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}
