package file

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/internal/dto"
	"github.com/edwardsean/codesmart/backend/internal/repository"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/password"
	"github.com/edwardsean/codesmart/backend/pkg/utils"
)

type FileService struct {
	projectRepo repository.ProjectRepository
	fileRepo    repository.ProjectFileRepository
	userRepo    repository.UserRepository
	httpClient  *http.Client
}

func NewFileService(projectRepo repository.ProjectRepository, fileRepo repository.ProjectFileRepository, userRepo repository.UserRepository) *FileService {
	return &FileService{
		projectRepo: projectRepo,
		fileRepo:    fileRepo,
		userRepo:    userRepo,
		httpClient:  &http.Client{Timeout: 10 * time.Second},
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

	if project.SourceType == domain.SourceGithub {
		return s.getGithubFileTree(ctx, project)
	}

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

	if project.SourceType == domain.SourceGithub {
		return s.getGithubFileContent(ctx, project, path)
	}

	return s.getDBFileContent(ctx, projectId, path)
}

func (s *FileService) getGithubFileContent(ctx context.Context, project *domain.Project, path string) (*dto.FileContentDTO, error) {
	token, err := s.decryptGithubToken(ctx, project)
	if err != nil {
		return nil, err
	}

	owner, repo, err := parseGithubURL(project.SourceURL)
	if err != nil {
		return nil, err
	}

	// use GitHub tree API — recursive=1 gets the full tree in one call
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, path)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.httpClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, errors.NewError("failed to fetch GitHub file tree", http.StatusBadGateway)
	}
	defer resp.Body.Close()

	var result struct {
		Content  string `json:"content"` // base64 encoded
		Encoding string `json:"encoding"`
		Name     string `json:"name"`
	}

	if err := utils.ParseJson(resp.Body, &result); err != nil {
		return nil, errors.NewError("failed to parse file content", http.StatusInternalServerError)
	}

	//the content github returns is base64 encoded content
	decoded, err := base64.StdEncoding.DecodeString(
		strings.ReplaceAll(result.Content, "\n", ""),
	)
	if err != nil {
		return nil, errors.NewError("failed to decode file content", http.StatusInternalServerError)
	}

	return &dto.FileContentDTO{
		Path:    path,
		Name:    result.Name,
		Content: string(decoded),
	}, nil
}

func (s *FileService) getDBFileContent(ctx context.Context, projectId int, path string) (*dto.FileContentDTO, error) {
	file, err := s.fileRepo.GetFileContent(ctx, projectId, path)
	if err != nil {
		return nil, errors.NewError("failed to get files", http.StatusInternalServerError)
	}

	return dto.ToFileContentDTO(file), nil
}

func (s *FileService) getGithubFileTree(ctx context.Context, project *domain.Project) ([]dto.FileNodeDTO, error) {
	token, err := s.decryptGithubToken(ctx, project)
	if err != nil {
		return nil, err
	}

	owner, repo, err := parseGithubURL(project.SourceURL)
	if err != nil {
		return nil, err
	}

	// use GitHub tree API — recursive=1 gets the full tree in one call
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/HEAD?recursive=1", owner, repo)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := s.httpClient.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return nil, errors.NewError("failed to fetch GitHub file tree", http.StatusBadGateway)
	}
	defer resp.Body.Close()

	var result struct {
		Tree []struct {
			Path string `json:"path"`
			Type string `json:"type"` // "blob" = file, "tree" = dir
			URL  string `json:"url"`
		} `json:"tree"`
	}

	if err := utils.ParseJson(resp.Body, &result); err != nil {
		return nil, errors.NewError("failed to parse GitHub tree", http.StatusInternalServerError)
	}

	nodes := make([]dto.FileNodeDTO, 0, len(result.Tree)) //length = 0, capacity = number of files (so go doesnt need to allocate space when exceed to grow capacity)
	for i, item := range result.Tree {
		nodeType := "file"
		if item.Type == "tree" {
			nodeType = "dir"
		}
		nodes = append(nodes, dto.FileNodeDTO{ //append because length is 0 at first
			ID:   i + 1, //REVISE
			Path: item.Path,
			Name: getFileName(item.Path),
			Type: nodeType,
		})
	}

	return nodes, nil
}

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

func (s *FileService) decryptGithubToken(ctx context.Context, project *domain.Project) (string, error) {
	user, err := s.userRepo.GetUserByID(ctx, project.UserID)
	if err != nil {
		return "", errors.NewError("user not found", http.StatusNotFound)
	}

	if user.GithubToken == "" {
		return "", errors.NewError("GitHub not connected, connect your GitHub account first", http.StatusBadRequest)
	}

	key, err := base64.StdEncoding.DecodeString(config.Envs.EncryptionKey)
	if err != nil {
		return "", errors.NewError("encryption config error", http.StatusInternalServerError)
	}

	return password.Decrypt(user.GithubToken, key)

}

func getFileName(path string) string {
	parts := strings.Split(path, "/")
	return parts[len(parts)-1]
}

func parseGithubURL(url string) (owner, repo string, err error) {
	// https://github.com/owner/repo
	parts := strings.Split(strings.TrimPrefix(url, "https://github.com/"), "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid github url: %s", url)
	}
	return parts[0], parts[1], nil
}
