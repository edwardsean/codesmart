package project

import (
	"context"
	"encoding/base64"
	"log"
	"net/http"

	"github.com/edwardsean/codesmart/backend/internal/clients"
	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/container"
	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/internal/dto"
	"github.com/edwardsean/codesmart/backend/internal/repository"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/password"
	"github.com/edwardsean/codesmart/backend/pkg/validator"
)

type ProjectService struct {
	projectRepo      repository.ProjectRepository
	fileRepo         repository.ProjectFileRepository
	userRepo         repository.UserRepository
	githubClient     clients.GithubClient
	containerManager *container.ContainerManager
	workSpaceRoot    string
}

func NewProjectService(projectRepo repository.ProjectRepository, fileRepo repository.ProjectFileRepository, userRepo repository.UserRepository, githubClient clients.GithubClient, cm *container.ContainerManager) *ProjectService {
	return &ProjectService{projectRepo: projectRepo, fileRepo: fileRepo, userRepo: userRepo, githubClient: githubClient, workSpaceRoot: config.Envs.WorkSpaceRoot, containerManager: cm}
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

	var result *dto.ProjectResponseDTO
	var err error
	if payload.SourceType == domain.SourceGithub && payload.SourceURL != "" {
		// go s.seedGithubFiles(context.Background(), project, userID)
		result, err = s.cloneGithubProject(ctx, userID, project)
		if err != nil {
			s.projectRepo.DeleteProject(ctx, project.ID) //rollback
			return nil, err
		}
	} else {
		result, err = s.initScratchWorkspace(ctx, userID, project)
		if err != nil {
			s.projectRepo.DeleteProject(ctx, project.ID) //rollback
			return nil, err
		}
	}

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

	//delete the container and volume if they exist
	if project.ContainerID != "" && project.VolumeName != "" {
		if err := s.containerManager.DeleteWorkspaceContainer(ctx, project.ContainerID, project.VolumeName); err != nil {
			log.Printf("Warning: failed to delete container for project %d: %v", id, err)
		}
	}

	if err := s.projectRepo.DeleteProject(ctx, id); err != nil {
		return errors.NewError("failed to delete project", http.StatusInternalServerError)
	}

	return nil
}

// func (s *ProjectService) createWorkspace(userId, projectId int) (string, error) {
// 	workspacePath := filepath.Join(s.workSpaceRoot, fmt.Sprintf("user_%d", userId), fmt.Sprintf("project_%d", projectId))

// 	if err := os.MkdirAll(workspacePath, 0755); err != nil {
// 		return "", errors.NewError("failed to create workspace", http.StatusInternalServerError)
// 	}

// 	return workspacePath, nil
// }

func (s *ProjectService) initScratchWorkspace(ctx context.Context, userId int, project *domain.Project) (*dto.ProjectResponseDTO, error) {
	workspaceInfo, err := s.containerManager.CreateWorkspaceContainer(ctx, userId, project.ID, "", "")
	if err != nil {
		return nil, errors.NewError(err.Error(), http.StatusInternalServerError)
	}

	project.ContainerID = workspaceInfo.ContainerID
	project.VolumeName = workspaceInfo.VolumeName
	if err := s.projectRepo.UpdateProject(ctx, project); err != nil {
		log.Printf("Warning: failed to update project with container info: %v", err)
	}

	result := dto.ToProjectResponseDTO(project)
	return result, nil
	// workspacePath, err := s.createWorkspace(userId, project.ID)
	// if err != nil {
	// 	return err
	// }

	// cmd := exec.CommandContext(ctx, "git", "init", workspacePath)
	// if err := cmd.Run(); err != nil {
	// 	return errors.NewError("git init failed", http.StatusInternalServerError)
	// }

	// starterPath := filepath.Join(workspacePath, "main.go")
	// os.WriteFile(starterPath, []byte("package main\n\nfunc main() {\n\n}\n"), 0644)
	// exec.CommandContext(ctx, "git", "-C", workspacePath, "add", ".").Run()
	// exec.CommandContext(ctx, "git", "-C", workspacePath,
	// 	"commit", "-m", "Initial commit").Run()

	// return nil
}

func (s *ProjectService) cloneGithubProject(ctx context.Context, userId int, project *domain.Project) (*dto.ProjectResponseDTO, error) {
	user, err := s.userRepo.GetUserByID(ctx, userId)
	if err != nil {
		s.projectRepo.DeleteProject(ctx, project.ID) //rollback
		return nil, errors.NewError("failed to get user", http.StatusInternalServerError)
	}

	if user.GithubToken == "" {
		s.projectRepo.DeleteProject(ctx, project.ID)
		return nil, errors.NewError("user does not have github token", http.StatusBadRequest)
	}

	encryptionKey, err := base64.StdEncoding.DecodeString(config.Envs.EncryptionKey)
	if err != nil {
		s.projectRepo.DeleteProject(ctx, project.ID)
		return nil, errors.NewError("failed to decode encryption key", http.StatusInternalServerError)
	}

	token, err := password.Decrypt(user.GithubToken, encryptionKey)
	if err != nil {
		s.projectRepo.DeleteProject(ctx, project.ID)
		return nil, errors.NewError("failed to decrypt github token", http.StatusInternalServerError)
	}

	workspaceInfo, err := s.containerManager.CreateWorkspaceContainer(ctx, userId, project.ID, project.SourceURL, token)
	if err != nil {
		s.projectRepo.DeleteProject(ctx, project.ID)
		return nil, errors.NewError(err.Error(), http.StatusInternalServerError)
	}

	//update project container info
	project.ContainerID = workspaceInfo.ContainerID
	project.VolumeName = workspaceInfo.VolumeName
	if err := s.projectRepo.UpdateProject(ctx, project); err != nil {
		log.Printf("Warning: failed to update project with container info: %v", err)
	}

	result := dto.ToProjectResponseDTO(project)
	return result, nil

	// authURL := strings.Replace(project.SourceURL, "https://", fmt.Sprintf("https://oauth2:%s@", token), 1) //TODO: change this later maybe use GithubApp
	// workspacePath, err := s.createWorkspace(user.ID, project.ID)
	// if err != nil {
	// 	return err
	// }

	// //clone using the helper
	// cmd := exec.CommandContext(ctx, "git", "clone", authURL, workspacePath)
	// output, err := cmd.CombinedOutput()
	// if err != nil {
	// 	log.Printf("git clone failed: %s", string(output))
	// 	return errors.NewError("git clone failed", http.StatusInternalServerError)
	// }

	// return nil
}

func (s *ProjectService) Close() error {
	return s.containerManager.Close()
}

// func (s *ProjectService) seedGithubFiles(ctx context.Context, project *domain.Project, userId int) {
// 	user, err := s.userRepo.GetUserByID(ctx, userId)
// 	if err != nil {
// 		log.Printf("seedGithubFiles: failed to get user: %v", err)
// 		return
// 	}

// 	if user.GithubToken == "" {
// 		log.Printf("seedGithubFiles: user has no github token")
// 		return
// 	}

// 	encryptionKey, err := base64.StdEncoding.DecodeString(config.Envs.EncryptionKey)
// 	if err != nil {
// 		log.Printf("seedGithubFiles: failed to decode encryption key: %v", err)
// 		return
// 	}

// 	token, err := password.Decrypt(user.GithubToken, encryptionKey)
// 	if err != nil {
// 		log.Printf("seedGithubFiles: failed to decrypt token: %v", err)
// 		return
// 	}

// 	owner, repo, err := service.ParseGithubURL(project.SourceURL)
// 	if err != nil {
// 		log.Printf("seedGithubFiles: failed to parse github url: %v", err)
// 		return
// 	}

// 	result, err := s.githubClient.GetFileTree(ctx, token, owner, repo)
// 	if err != nil {
// 		log.Printf("seedGithubFiles: failed to get project tree: %v", err)
// 	}

// 	if result.Truncated {
// 		log.Printf("seedGithubFiles: tree is truncated (repo too large), partial seed")
// 	}

// 	//fetch content for each file (skip dirs and large files)
// 	const maxFileSizeBytes = 100 * 1024 // 100KB per file limit
// 	const maxFiles = 200                // don't seed more than 200 files

// 	seeded := 0
// 	for _, item := range result.Tree {
// 		if item.Type != "blob" {
// 			continue
// 		}

// 		if item.Size > maxFileSizeBytes {
// 			log.Printf("seedGithubFiles: skipping large file %s (%d bytes)", item.Path, item.Size)
// 			continue
// 		}

// 		if seeded >= maxFiles {
// 			log.Printf("seedGithubFiles: reached max file limit (%d)", maxFiles)
// 			break
// 		}

// 		if shouldSkipFile(item.Path) {
// 			continue
// 		}

// 		content, err := s.githubClient.GetFileContent(ctx, token, owner, repo, item.Path)
// 		if err != nil {
// 			log.Printf("seedGithubFiles: failed to fetch %s: %v", item.Path, err)
// 			continue
// 		}

// 		file := &domain.ProjectFile{
// 			ProjectID: project.ID,
// 			FilePath:  item.Path,
// 			Content:   content,
// 			Language:  domain.Language(service.DetectLanguage(item.Path)),
// 		}

// 		if err := s.fileRepo.CreateFile(ctx, file); err != nil {
// 			log.Printf("seedGithubFiles: failed to save %s: %v", item.Path, err)
// 			continue
// 		}

// 		seeded++
// 	}

// 	log.Printf("seedGithubFiles: seeded %d files for project %d", seeded, project.ID)

// }
