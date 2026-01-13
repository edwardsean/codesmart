package postgres

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"

	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/pkg/password"
	"github.com/edwardsean/codesmart/backend/pkg/response"
	"gorm.io/gorm"
)

type PostgreUserStore struct {
	db *gorm.DB
}

func PostgreNewUserStore(db *gorm.DB) *PostgreUserStore {
	return &PostgreUserStore{db: db}
}

func (s *PostgreUserStore) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User

	err := s.db.Raw("SELECT * FROM users WHERE email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *PostgreUserStore) GetUserByID(id int) (*domain.User, error) {
	var user domain.User

	err := s.db.Raw("SELECT * FROM users WHERE id = ?", id).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil

}

func (s *PostgreUserStore) CreateUser(User domain.User) error {
	user := domain.User{Username: User.Username, Email: User.Email, Password: User.Password}

	result := s.db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *PostgreUserStore) GetOrCreateUserFromGithub(id int, username string, email string, access_token string, github_user *domain.GithubUser) (*domain.User, error) {
	var user domain.User

	if email == "" {
		var err error
		email, err = GetGithubEmail(access_token, github_user)
		if err != nil {
			return nil, err
		}
	}
	//hash access token to store to database
	encryption_key64 := config.Envs.EncryptionKey
	secretKey, err := base64.StdEncoding.DecodeString(encryption_key64)

	if err != nil {
		return nil, err
	}

	hashed_token, err := password.Encrypt(access_token, secretKey)

	if err != nil {
		return nil, err
	}

	err = s.db.Raw("SELECT * FROM users WHERE email = ?", email).First(&user).Error

	//if doesnt exist
	if errors.Is(err, gorm.ErrRecordNotFound) {
		new_user := domain.User{Username: username, Email: email, Password: "", GitHubID: id, GithubToken: hashed_token}
		result := s.db.Create(&new_user)

		if result.Error != nil {
			return nil, result.Error
		}

		return &new_user, nil
	} else if err != nil {
		return nil, err
	}

	//if exists
	if user.GithubToken != hashed_token {
		user.GithubToken = hashed_token
		if err := s.db.Save(&user).Error; err != nil {
			return nil, err
		}
	}
	return &user, nil

}

func GetGithubEmail(access_token string, github_user *domain.GithubUser) (string, error) {
	request, _ := http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	request.Header.Set("Authorization", "token "+access_token)
	resp, err := http.DefaultClient.Do(request)

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

	if err := response.ParseJson(resp.Body, &emails); err != nil {
		return "", err
	}

	for _, e := range emails {
		if e.Primary && e.Verified {
			github_user.Email = e.Email
			return e.Email, nil
		}
	}

	return "", fmt.Errorf("no primary email found")

}
