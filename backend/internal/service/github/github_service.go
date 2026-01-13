package github

import (
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/password"
	"github.com/edwardsean/codesmart/backend/pkg/response"
)

type GithubService struct {
	httpClient *http.Client
}

func NewGithubService() *GithubService {
	return &GithubService{httpClient: &http.Client{Timeout: 10 * time.Second}}
}

func (s *GithubService) GetUserRepositories(user *domain.User) (*[]domain.Repository, error) {
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
	request, err := http.NewRequest("GET", "https://api.github.com/user/repos", nil)

	if err != nil {
		return nil, errors.NewError("unable to fetch github repo api", http.StatusBadRequest)
	}

	request.Header.Set("Authorization", "Bearer "+ghAccessToken)
	resp, err := s.httpClient.Do(request)
	if err != nil {
		return nil, errors.NewError(fmt.Sprintf("unable to fetch github repo api: %v", err), http.StatusBadRequest)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.NewError(fmt.Sprintf("GitHub API error: status=%d, body=%s", resp.StatusCode, string(body)), http.StatusBadRequest)
	}

	var repositories []domain.Repository

	if err := response.ParseJson(resp.Body, &repositories); err != nil {
		return nil, errors.NewError("unable to parse github repos", http.StatusBadRequest)
	}

	return &repositories, nil
}
