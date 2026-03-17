package postgres

import (
	"context"
	"encoding/base64"
	"errors"

	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/pkg/password"
	"gorm.io/gorm"
)

type PostgreUserStore struct {
	db *gorm.DB
}

func PostgreNewUserStore(db *gorm.DB) *PostgreUserStore {
	return &PostgreUserStore{db: db}
}

func (s *PostgreUserStore) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var user domain.User

	err := s.db.WithContext(ctx).Raw("SELECT * FROM users WHERE email = ?", email).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (s *PostgreUserStore) GetUserByID(ctx context.Context, id int) (*domain.User, error) {
	var user domain.User

	err := s.db.WithContext(ctx).Raw("SELECT * FROM users WHERE id = ?", id).First(&user).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &user, nil

}

func (s *PostgreUserStore) CreateUser(ctx context.Context, user *domain.User) error {
	result := s.db.WithContext(ctx).Create(user) //dont use &pointer since that would be a pointer to a pointer **domain.User

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *PostgreUserStore) GetOrCreateUserFromGithub(ctx context.Context, id int, email string, username string, access_token string, github_user *domain.GithubUser) (*domain.User, error) {
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

	var user domain.User
	err = s.db.WithContext(ctx).Raw("SELECT * FROM users WHERE email = ?", email).First(&user).Error

	//if doesnt exist
	if errors.Is(err, gorm.ErrRecordNotFound) {
		new_user := domain.User{Username: username, Email: email, Password: "", GitHubID: id, GithubToken: hashed_token}
		result := s.db.WithContext(ctx).Create(&new_user)

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
		if err := s.db.WithContext(ctx).Save(&user).Error; err != nil {
			return nil, err
		}
	}
	return &user, nil

}
