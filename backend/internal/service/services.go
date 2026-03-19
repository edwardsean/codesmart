package service

import (
	"context"

	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/internal/dto"
)

type UserService interface {
	GetUserByEmail(ctx context.Context, email string) (*domain.User, error)
}

type GithubService interface {
	GetUserRepositories(ctx context.Context, user *domain.User) (*[]domain.GithubRepository, error)
}

type OAuthService interface {
	GithubCallback(ctx context.Context, code string) (string, error)
	ExchangeOAuthCode(ctx context.Context, code string) (*dto.OAuthResultDTO, error)
	ConnectGithub(ctx context.Context, connect_code string, code string) (string, error)
	ConnectGithubInit(ctx context.Context, userId int, redirect string) (string, error)
}

type AuthService interface {
	Login(ctx context.Context, payload dto.LoginUserPayload) (*dto.AuthResultDTO, error)
	Register(ctx context.Context, payload dto.RegisterUserPayload) (*dto.AuthResultDTO, error)
	Logout(ctx context.Context, refreshToken string) error
	ValidateRefreshToken(ctx context.Context, refreshToken string) (*domain.User, error)
	GetUserFromToken(ctx context.Context, token string) (*domain.User, error)
}

type ProjectService interface {
	GetProjects(ctx context.Context, userId int) ([]dto.ProjectListItemDTO, error)
	CreateProject(ctx context.Context, userId int, payload dto.CreateProjectPayload) (*dto.ProjectResponseDTO, error)
	GetProjectByID(ctx context.Context, projectId int, userId int) (*dto.ProjectResponseDTO, error)
	DeleteProject(ctx context.Context, projectId int, userId int) error
}
