package clients

import (
	"context"

	"github.com/edwardsean/codesmart/backend/internal/clients/github"
	"github.com/edwardsean/codesmart/backend/internal/domain"
)

type GithubClient interface {
	GetFileContent(ctx context.Context, token, owner, repo, path string) (string, error)
	GetFileTree(ctx context.Context, token, owner, repo string) (*github.FileTreeResult, error)
	GetRepositories(ctx context.Context, token string) ([]domain.GithubRepository, error)
}
