package repository

import (
	"context"
	"time"

	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/internal/repository/redis"
)

type UserRepository interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
	GetUserByID(ctx context.Context, id int) (*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetOrCreateUserFromGithub(ctx context.Context, id int, email string, login string, hashed_token string, github_user *domain.GithubUser) (*domain.User, error)
	DeleteUser(ctx context.Context, id int) error
	LinkGithub(ctx context.Context, userId int, github_id int, hashed_token string) error
}

type TokenRepository interface {
	BlacklistToken(ctx context.Context, token string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	StoreOAuthCode(ctx context.Context, code string, data *redis.OAuthCodeData) error
	ExchangeOAuthCode(ctx context.Context, code string) (*redis.OAuthCodeData, error)
	StoreConnectCode(ctx context.Context, code string, data *redis.ConnectCodeData) error
	ExchangeConnectCode(ctx context.Context, connect_code string) (*redis.ConnectCodeData, error)
	StoreWSTicket(ctx context.Context, ticket string, userId int) error
	ExchangeWSTicket(ctx context.Context, ticket string) (int, error)
}

type ProjectRepository interface {
	GetProjectByID(ctx context.Context, projectId int) (*domain.Project, error)
	GetProjectsByUserID(ctx context.Context, userId int) ([]domain.Project, error)
	CreateProject(ctx context.Context, project *domain.Project) error
	DeleteProject(ctx context.Context, projectId int) error
}

type ProjectFileRepository interface {
	GetFiles(ctx context.Context, projectID int) ([]domain.ProjectFile, error)
	GetFileByID(ctx context.Context, fileId int) (*domain.ProjectFile, error)
	GetFileContent(ctx context.Context, projectID int, path string) (*domain.ProjectFile, error)
	CreateFile(ctx context.Context, file *domain.ProjectFile) error
	UpdateFile(ctx context.Context, file *domain.ProjectFile) error
	DeleteFile(ctx context.Context, fileId int) error
}
