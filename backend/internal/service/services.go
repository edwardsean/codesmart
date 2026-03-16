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
	HandleGithubCallback(ctx context.Context, code string) (string, error)
}

type AuthService interface {
	Login(ctx context.Context, payload dto.LoginUserPayload) (string, string, *dto.UserResponseDTO, error)
	Register(ctx context.Context, payload dto.RegisterUserPayload) (string, string, *dto.UserResponseDTO, error)
	Logout(ctx context.Context, refreshToken string) error
	ValidateRefreshToken(ctx context.Context, refreshToken string) (*domain.User, error)
	GetUserFromToken(ctx context.Context, token string) (*domain.User, error)
}
