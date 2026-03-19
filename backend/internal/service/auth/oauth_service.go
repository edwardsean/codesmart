package auth

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/internal/dto"
	"github.com/edwardsean/codesmart/backend/internal/repository"
	"github.com/edwardsean/codesmart/backend/internal/repository/redis"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/jwt"
	"github.com/edwardsean/codesmart/backend/pkg/password"
	"github.com/edwardsean/codesmart/backend/pkg/utils"
	"github.com/google/uuid"
)

type oauthService struct {
	userRepo   repository.UserRepository
	tokenRepo  repository.TokenRepository
	httpClient *http.Client
}

func NewOAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) *oauthService {
	return &oauthService{
		userRepo:   userRepo,
		tokenRepo:  tokenRepo,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *oauthService) exchangeCodeForToken(ctx context.Context, code string) (string, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		"https://github.com/login/oauth/access_token",
		nil,
	)
	//http.DefaultClient has no timeout — if GitHub never responds, your goroutine hangs forever. Fine for quick scripts, dangerous in a server.
	//you control the timeout
	//httpClient: &http.Client{Timeout: 10 * time.Second}
	//you can't use package-level functions with context
	//context support — requires building request manually + client.Do() req, _ := http.NewRequestWithContext(ctx, http.MethodPost, url, nil) s.httpClient.Do(req)

	if err != nil {
		return "", err
	}

	q := url.Values{
		"client_id":     {config.Envs.GithubClientID},
		"client_secret": {config.Envs.GithubSecret},
		"code":          {code},
	}
	req.URL.RawQuery = q.Encode()
	req.Header.Set("Accept", "application/x-www-form-urlencoded")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close() //must close it when done, otherwise it leaks connections. defer schedules it to run at the end of the function, so we dont have to remember to close it manually later.

	body, err := io.ReadAll(resp.Body)
	log.Printf("GitHub access token raw body: %s", string(body))

	values, err := url.ParseQuery(string(body))
	if err != nil {
		return "", err
	}

	ghAccessToken := values.Get("access_token")
	log.Printf("Parsed ghAccessToken: %s", ghAccessToken)

	if ghAccessToken == "" {
		//ERROR HANDLING
		// http.Redirect(w, r, "http://localhost:3000/login?error=no_github_token", http.StatusFound)
		return "", errors.NewError("no github access token found", http.StatusInternalServerError)
	}

	return ghAccessToken, nil

}

func (s *oauthService) getGithubUser(ctx context.Context, accessToken string) (*domain.GithubUser, error) {
	//get user from github
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)

	if err != nil {
		//ERROR HANDLING
		return nil, err
	}

	request.Header.Set("Authorization", "token "+accessToken)
	resp, err := s.httpClient.Do(request)
	if err != nil || resp.StatusCode != 200 {
		//ERROR HANDLING
		// http.Redirect(w, r, "http://localhost:3000/login?error=github_user_failed", http.StatusFound)
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch github user: status %d", resp.StatusCode)
	}

	var github_user domain.GithubUser

	if err := utils.ParseJson(resp.Body, &github_user); err != nil {
		//ERROR HANDLING
		// http.Redirect(w, r, "http://localhost:3000/login?error=decode_failed", http.StatusFound)
		return nil, err
	}

	return &github_user, nil
}

func (s *oauthService) GithubCallback(ctx context.Context, code string) (string, error) {
	// uri := config.Envs.GolangAPIURL + "/auth/github/callback"
	// tokenResp, err := http.PostForm("https://github.com/login/oauth/access_token", url.Values{"client_id": {config.Envs.GithubClientID}, "client_secret": {config.Envs.GithubSecret}, "code": {code}, "redirect_uri": {uri}}) //returns a tokenResp.body

	// if err != nil {
	// 	//ERROR HANDLING
	// 	// http.Redirect(w, r, "http://localhost:3000/login?error=exchange_failed", http.StatusFound)
	// 	log.Printf("error getting the access token response from github: %v", err)
	// 	return
	// }

	// defer tokenResp.Body.Close()

	// body, _ := io.ReadAll(tokenResp.Body)
	// log.Printf("GitHub access token raw body: %s", string(body))
	// values, _ := url.ParseQuery(string(body))
	// ghAccessToken := values.Get("access_token")
	// log.Printf("Parsed ghAccessToken: %s", ghAccessToken)

	// if ghAccessToken == "" {
	// 	//ERROR HANDLING

	// 	return
	// }
	ghAccessToken, err := s.exchangeCodeForToken(ctx, code)
	if err != nil {
		return "", err
	}

	github_user, err := s.getGithubUser(ctx, ghAccessToken)

	if github_user.Email == "" {
		email, err := s.getGithubEmail(ctx, ghAccessToken)
		if err != nil {
			return "", err
		}
		github_user.Email = email
	}

	secretKey, _ := base64.StdEncoding.DecodeString(config.Envs.EncryptionKey)
	hashed_token, err := password.Encrypt(ghAccessToken, secretKey)

	if err != nil {
		return "", err
	}

	//get or create user for this github account
	user, err := s.userRepo.GetOrCreateUserFromGithub(ctx, github_user.ID, github_user.Email, github_user.Login, hashed_token, github_user)

	if err != nil {
		// http.Redirect(w, r, "http://localhost:3000/login?error=user_db_fetching_failed_for_github, http.StatusFound)
		return "", err
	}

	secret := []byte(config.Envs.JWTSecret)
	// access_token, err := auth.CreateJWT(secret, user.ID, 15*time.Minute)

	// if err != nil {
	// 	// http.Redirect(w, r, "http://localhost:3000/login?error=access_token_creation_failure, http.StatusFound)
	// 	return
	// }

	refresh_token, err := jwt.CreateJWT(secret, user.ID, 7*24*time.Hour)
	if err != nil {
		// http.Redirect(w, r, "http://localhost:3000/login?error=refresh_token_creation_failure, http.StatusFound)
		return "", err
	}

	access_token, err := jwt.CreateJWT(secret, user.ID, 15*time.Minute)
	if err != nil {
		// http.Redirect(w, r, "http://localhost:3000/login?error=refresh_token_creation_failure, http.StatusFound)
		return "", err
	}

	oauthCode := uuid.New().String()
	err = s.tokenRepo.StoreOAuthCode(ctx, oauthCode, &redis.OAuthCodeData{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		User:         user,
	})
	if err != nil {
		return "", err
	}

	return oauthCode, nil
}

func (s *oauthService) ExchangeOAuthCode(ctx context.Context, code string) (*dto.OAuthResultDTO, error) {
	data, err := s.tokenRepo.ExchangeOAuthCode(ctx, code)
	if err != nil {
		return nil, err
	}

	return &dto.OAuthResultDTO{
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		User: dto.UserResponseDTO{
			ID:       data.User.ID,
			Email:    data.User.Email,
			Username: data.User.Username,
		},
	}, nil

}

func (s *oauthService) ConnectGithub(ctx context.Context, connect_code string, code string) (string, error) {
	data, err := s.tokenRepo.ExchangeConnectCode(ctx, connect_code)
	if err != nil {
		return "/dashboard", err
	}

	ghAccessToken, err := s.exchangeCodeForToken(ctx, code)
	if err != nil {
		return "/dashboard", err
	}

	githubUser, err := s.getGithubUser(ctx, ghAccessToken)
	if err != nil {
		return "/dashboard", err
	}

	if githubUser.Email == "" {
		githubUser.Email, err = s.getGithubEmail(ctx, ghAccessToken)
		if err != nil {
			return "/dashboard", err
		}
	}

	encryptionKey, _ := base64.StdEncoding.DecodeString(config.Envs.EncryptionKey)
	hashedToken, err := password.Encrypt(ghAccessToken, encryptionKey)
	if err != nil {
		return "/dashboard", err
	}

	// link github to the existing user — don't create a new user
	return data.Redirect, s.userRepo.LinkGithub(ctx, data.UserID, githubUser.ID, hashedToken)
}

func (s *oauthService) ConnectGithubInit(ctx context.Context, userId int, redirect string) (string, error) {
	connectCode := uuid.New().String()
	err := s.tokenRepo.StoreConnectCode(ctx, connectCode, &redis.ConnectCodeData{UserID: userId, Redirect: redirect})
	if err != nil {
		return "", errors.NewError("failed", http.StatusInternalServerError)
	}

	return connectCode, nil
}

func (s *oauthService) getGithubEmail(ctx context.Context, access_token string) (string, error) {
	request, _ := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user/emails", nil)
	request.Header.Set("Authorization", "token "+access_token)
	resp, err := s.httpClient.Do(request)

	if err != nil || resp.StatusCode != 200 {
		return "", fmt.Errorf("failed to fetch GitHub emails, status: %d, err: %w", resp.StatusCode, err)
	}

	defer resp.Body.Close()

	var emails []struct {
		Email      string `json:"email"`
		Verified   bool   `json:"verified"`
		Primary    bool   `json:"primary"`
		Visibility string `json:"visibility"`
	}

	if err := utils.ParseJson(resp.Body, &emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			return e.Email, nil
		}
	}

	return "", fmt.Errorf("no primary email found")

}
