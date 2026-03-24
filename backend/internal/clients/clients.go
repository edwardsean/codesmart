package clients

import (
	"context"

	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/internal/dto"
)

type GithubClient interface {
	GetFileContent(ctx context.Context, token, owner, repo, path string) (string, error)
	GetFileTree(ctx context.Context, token, owner, repo string) (*dto.FileTreeResult, error)
	GetRepositories(ctx context.Context, token string) ([]domain.GithubRepository, error)
}
