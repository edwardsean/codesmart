package github

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/utils"
)

type Client struct {
	httpClient *http.Client
}

func NewGithubClient(httpClient *http.Client) *Client {
	return &Client{httpClient: httpClient}
}

type FileTreeResult struct {
	Tree      []FileTreeItem `json:"tree"`
	Truncated bool           `json:"truncated"`
}

type FileTreeItem struct {
	Path string `json:"path"`
	Type string `json:"type"` // "blob" = file, "tree" = dir
	Size int    `json:"size"`
}

func (c *Client) GetFileContent(ctx context.Context, token, owner, repo, path string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, path)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("github API error: %d", resp.StatusCode)
	}

	var result struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
		Name     string `json:"name"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.Encoding != "base64" {
		return "", fmt.Errorf("unexpected encoding: %s", result.Encoding)
	}

	decoded, err := base64.StdEncoding.DecodeString(
		strings.ReplaceAll(result.Content, "\n", ""),
	)
	if err != nil {
		return "", err
	}

	return string(decoded), nil
}

func (c *Client) GetFileTree(ctx context.Context, token, owner, repo string) (*FileTreeResult, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/git/trees/HEAD?recursive=1", owner, repo)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Accept", "application/vnd.github+json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("github API error: %d", resp.StatusCode)
	}

	var result FileTreeResult

	if err := utils.ParseJson(resp.Body, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

func (c *Client) GetRepositories(ctx context.Context, token string) ([]domain.GithubRepository, error) {
	request, err := http.NewRequestWithContext(ctx, "GET", "https://api.github.com/user/repos", nil)

	if err != nil {
		return nil, errors.NewError("unable to fetch github repo api", http.StatusBadRequest)
	}

	request.Header.Set("Authorization", "Bearer "+token)
	resp, err := c.httpClient.Do(request)
	if err != nil {
		return nil, errors.NewError(fmt.Sprintf("unable to fetch github repo api: %v", err), http.StatusBadRequest)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.NewError(fmt.Sprintf("GitHub API error: status=%d, body=%s", resp.StatusCode, string(body)), http.StatusBadRequest)
	}

	var repositories []domain.GithubRepository

	if err := utils.ParseJson(resp.Body, &repositories); err != nil {
		return nil, errors.NewError("unable to parse github repos", http.StatusBadRequest)
	}

	return repositories, nil
}
