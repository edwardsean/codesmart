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
	CreateUser(ctx context.Context, user *domain.User) error
	GetOrCreateUserFromGithub(ctx context.Context, id int, email string, login string, access_token string, github_user *domain.GithubUser) (*domain.User, error)
}

type TokenRepository interface {
	BlacklistToken(ctx context.Context, token string, ttl time.Duration) error
	IsBlacklisted(ctx context.Context, token string) (bool, error)
	StoreOAuthCode(ctx context.Context, code string, data *redis.OAuthCodeData) error
	ExchangeOAuthCode(ctx context.Context, code string) (*redis.OAuthCodeData, error)
}
