package auth

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/pkg/jwt"
	"github.com/edwardsean/codesmart/backend/pkg/response"
)

type oauthService struct {
	userRepo   domain.UserRepository
	httpClient *http.Client
}

func NewOAuthService(userRepo domain.UserRepository) *oauthService {
	return &oauthService{
		userRepo:   userRepo,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *oauthService) exchangeCodeForToken(code string) (string, error) {
	resp, err := s.httpClient.PostForm("https://github.com/login/oauth/access_token",
		url.Values{
			"client_id":     {config.Envs.GithubClientID},
			"client_secret": {config.Envs.GithubSecret},
			"code":          {code},
		})

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
		return "", errors.New("no github access token found")
	}

	return ghAccessToken, nil

}

func (s *oauthService) getGithubUser(accessToken string) (*domain.GithubUser, error) {
	//get user from github
	request, err := http.NewRequest("GET", "https://api.github.com/user", nil)

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

	if err := response.ParseJson(resp.Body, &github_user); err != nil {
		//ERROR HANDLING
		// http.Redirect(w, r, "http://localhost:3000/login?error=decode_failed", http.StatusFound)
		return nil, err
	}

	return &github_user, nil
}

func (s *oauthService) HandleGithubCallback(code string) (string, error) {
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
	ghAccessToken, err := s.exchangeCodeForToken(code)
	if err != nil {
		return "", err
	}

	github_user, err := s.getGithubUser(ghAccessToken)

	//get or create user for this github account
	user, err := s.userRepo.GetOrCreateUserFromGithub(github_user.ID, github_user.Login, github_user.Email, ghAccessToken, github_user)

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

	return refresh_token, nil
}
