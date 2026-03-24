package file

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/edwardsean/codesmart/backend/internal/clients"
	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/dto"
	"github.com/edwardsean/codesmart/backend/internal/repository"
	"github.com/edwardsean/codesmart/backend/internal/service"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
)

type FileService struct {
	projectRepo   repository.ProjectRepository
	fileRepo      repository.ProjectFileRepository
	userRepo      repository.UserRepository
	githubClient  clients.GithubClient
	workSpaceRoot string
}

func NewFileService(projectRepo repository.ProjectRepository, fileRepo repository.ProjectFileRepository, userRepo repository.UserRepository, githubClient clients.GithubClient) *FileService {
	return &FileService{
		projectRepo:   projectRepo,
		fileRepo:      fileRepo,
		userRepo:      userRepo,
		githubClient:  githubClient,
		workSpaceRoot: config.Envs.WorkSpaceRoot,
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

	workspacePath := s.getWorkspacePath(userId, projectId)

	var nodes []dto.FileNodeDTO
	err = filepath.WalkDir(workspacePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		//skip .git folder
		if d.IsDir() && d.Name() == ".git" {
			return filepath.SkipDir
		}

		//skip node modules, etc
		if service.ShouldSkipFile(d.Name()) {
			return filepath.SkipDir
		}

		relPath, _ := filepath.Rel(workspacePath, path)
		if relPath == "." {
			return nil
		}

		nodeType := "file"
		if d.IsDir() {
			nodeType = "dir"
		}

		nodes = append(nodes, dto.FileNodeDTO{
			Path: relPath,
			Name: d.Name(),
			Type: nodeType,
		})

		return nil
	})

	log.Printf("file tree: %v", nodes)
	log.Printf("file dir: %s", workspacePath)
	return nodes, err

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

	workspacePath := s.getWorkspacePath(userId, projectId)
	fullPath := filepath.Join(workspacePath, path)

	if !strings.HasPrefix(fullPath, workspacePath) { //to prevent traversal attacks in full path
		return nil, errors.NewError("invalid path", http.StatusBadRequest)
	}

	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, errors.NewError("file not found", http.StatusNotFound)
	}

	return &dto.FileContentDTO{
		Path:    path,
		Name:    filepath.Base(path),
		Content: string(content),
	}, nil
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

	workspacePath := s.getWorkspacePath(userId, projectId)
	fullPath := filepath.Join(workspacePath, payload.Path)

	if !strings.HasPrefix(fullPath, workspacePath) {
		return nil, errors.NewError("invalid path", http.StatusBadRequest)
	}

	if _, err := os.Stat(fullPath); err != nil {
		return nil, errors.NewError("file already exists", http.StatusConflict)
	}

	if payload.IsDir {
		//create directory
		if err := os.MkdirAll(fullPath, 0755); err != nil {
			return nil, errors.NewError("failed to create directory", http.StatusInternalServerError)
		}

		return &dto.FileContentDTO{
			Path:    payload.Path,
			Name:    filepath.Base(payload.Path),
			Content: "",
		}, nil
	} else {
		//create file
		//create parent directory if doesnt exist
		parentDir := filepath.Dir(fullPath)
		if err := os.MkdirAll(parentDir, 0755); err != nil {
			return nil, errors.NewError("faied to create parent directory", http.StatusInternalServerError)
		}

		//create file
		if err := os.WriteFile(fullPath, []byte(payload.Content), 0644); err != nil {
			return nil, errors.NewError("failed to create file", http.StatusInternalServerError)
		}

		return &dto.FileContentDTO{
			Path:    payload.Path,
			Name:    filepath.Base(payload.Path),
			Content: payload.Content,
		}, nil
	}

	// file := &domain.ProjectFile{
	// 	ProjectID: projectId,
	// 	FilePath:  payload.Path,
	// 	Content:   payload.Content,
	// 	Language:  domain.Language(payload.Language),
	// }

	// if err := s.fileRepo.CreateFile(ctx, file); err != nil {
	// 	return nil, errors.NewError("failed to create file", http.StatusInternalServerError)
	// }

	// return dto.ToFileContentDTO(file), nil
}

func (s *FileService) UpdateFile(ctx context.Context, projectId, userId int, payload dto.UpdateFilePayload) error {
	project, err := s.projectRepo.GetProjectByID(ctx, projectId)
	if err != nil {
		return errors.NewError("project not found", http.StatusNotFound)
	}

	//ownership
	if project.UserID != userId {
		return errors.NewError("forbidden", http.StatusForbidden)
	}

	workspacePath := s.getWorkspacePath(userId, projectId)
	fullPath := filepath.Join(workspacePath, payload.Path)

	if !strings.HasPrefix(fullPath, workspacePath) {
		return errors.NewError("invalid path", http.StatusBadRequest)
	}

	return os.WriteFile(fullPath, []byte(payload.Content), 0644)

	// existing, err := s.fileRepo.GetFileByID(ctx, fileId)
	// if err != nil {
	// 	return errors.NewError("file not found", http.StatusNotFound)
	// }

	// if existing.ProjectID != projectId {
	// 	return errors.NewError("file does not belong to this project", http.StatusForbidden)
	// }

	// if payload.Path != "" {
	// 	existing.FilePath = payload.Path
	// 	if payload.Language == "" {
	// 		existing.Language = domain.Language(service.DetectLanguage(payload.Path))
	// 	}
	// }

	// if payload.Content != "" {
	// 	existing.Content = payload.Content
	// }

	// return s.fileRepo.UpdateFile(ctx, existing)

}

func (s *FileService) DeleteFile(ctx context.Context, projectId, userId int, path string) error {
	project, err := s.projectRepo.GetProjectByID(ctx, projectId)
	if err != nil {
		return errors.NewError("project not found", http.StatusNotFound)
	}

	//ownership
	if project.UserID != userId {
		return errors.NewError("forbidden", http.StatusForbidden)
	}

	workspacePath := s.getWorkspacePath(userId, projectId)
	fullPath := filepath.Join(workspacePath, path)

	if !strings.HasPrefix(fullPath, workspacePath) {
		return errors.NewError("invalid path", http.StatusBadRequest)
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return errors.NewError("file not found", http.StatusNotFound)
	}

	if err := os.Remove(fullPath); err != nil {
		return errors.NewError("failed to delete file", http.StatusInternalServerError)
	}

	return nil
	// file, err := s.fileRepo.GetFileByID(ctx, fileId)
	// if err != nil {
	// 	return errors.NewError("file not found", http.StatusNotFound)
	// }

	// if file.ProjectID != projectId {
	// 	return errors.NewError("file does not belong to this project", http.StatusForbidden)
	// }

	// return s.fileRepo.DeleteFile(ctx, fileId)
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
			Path: file.FilePath,
			Name: getFileName(file.FilePath),
			Type: "file", //REVISE
		}
	}

	return nodes, nil
}

func (s *FileService) getWorkspacePath(userId, projectId int) string {
	return filepath.Join(s.workSpaceRoot, fmt.Sprintf("user_%d", userId), fmt.Sprintf("project_%d", projectId))
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
