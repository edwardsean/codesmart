package github

import (
	"context"
	"encoding/base64"
	"net/http"

	"github.com/edwardsean/codesmart/backend/internal/clients"
	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/password"
)

type GithubService struct {
	githubClient clients.GithubClient
}

func NewGithubService(githubClient clients.GithubClient) *GithubService {
	return &GithubService{githubClient: githubClient}
}

func (s *GithubService) GetUserRepositories(ctx context.Context, user *domain.User) ([]domain.GithubRepository, error) {
	//if there is no github token (user didnt login using github)
	if user.GithubToken == "" {
		return nil, errors.NewError("please login using github", http.StatusBadRequest)
	}

	//decrypt the github token
	encryption_key64 := config.Envs.EncryptionKey
	secretKey, err := base64.StdEncoding.DecodeString(encryption_key64)
	if err != nil {
		return nil, errors.NewError("unable to decode github token", http.StatusBadRequest)
	}

	ghAccessToken, err := password.Decrypt(user.GithubToken, secretKey)
	if err != nil {
		return nil, errors.NewError("unable to decrypt github token", http.StatusBadRequest)
	}

	//get repos
	repositories, err := s.githubClient.GetRepositories(ctx, ghAccessToken)

	return repositories, nil
}
