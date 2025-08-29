package user

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/edwardsean/codesmart/backend/auth-service/types"
	"github.com/edwardsean/codesmart/backend/auth-service/utils"
	"gorm.io/gorm"
)

type PostgreUserStore struct {
	db *gorm.DB
}

func PostgreNewStore(db *gorm.DB) *PostgreUserStore {
	return &PostgreUserStore{db: db}
}

func (s *PostgreUserStore) GetUserByEmail(email string) (*types.User, error) {
	var user types.User

	err := s.db.Raw("SELECT * FROM users WHERE email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *PostgreUserStore) GetUserByID(id int) (*types.User, error) {
	var user types.User

	err := s.db.Raw("SELECT * FROM users WHERE id = ?", id).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil

}

func (s *PostgreUserStore) CreateUser(User types.User) error {
	user := types.User{Username: User.Username, Email: User.Email, Password: User.Password}

	result := s.db.Create(&user)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *PostgreUserStore) GetOrCreateUserFromGithub(id int, username string, email string, access_token string, github_user *types.GithubUser) (*types.User, error) {
	var user types.User

	if email == "" {
		var err error
		email, err = GetGithubEmail(access_token, github_user)
		if err != nil {
			return nil, err
		}
	}

	err := s.db.Raw("SELECT * FROM users WHERE email = ?", email).First(&user).Error

	//if doesnt exists
	if errors.Is(err, gorm.ErrRecordNotFound) {
		new_user := types.User{Username: username, Email: email, Password: "", GitHubID: id}
		result := s.db.Create(&new_user)

		if result.Error != nil {
			return nil, result.Error
		}

		return &new_user, nil
	} else if err != nil {
		return nil, err
	}

	//if exists
	return &user, nil

}

func GetGithubEmail(access_token string, github_user *types.GithubUser) (string, error) {
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

	if err := utils.ParseJson(resp.Body, &emails); err != nil {
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
