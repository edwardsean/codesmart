package postgres

import (
	"context"
	"errors"

	"github.com/edwardsean/codesmart/backend/internal/domain"
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

func (s *PostgreUserStore) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	result := s.db.WithContext(ctx).Create(user) //dont use &pointer since that would be a pointer to a pointer **domain.User

	if result.Error != nil {
		return nil, result.Error
	}

	return user, nil
}

func (s *PostgreUserStore) GetOrCreateUserFromGithub(ctx context.Context, id int, email string, username string, hashed_token string, github_user *domain.GithubUser) (*domain.User, error) {
	//hash access token to store to database

	var user domain.User

	err := s.db.WithContext(ctx).
		Raw("SELECT * FROM users WHERE github_id = ?", id).
		First(&user).Error

	if err == nil {
		//found by github_id, then update token and return
		user.GithubToken = hashed_token
		s.db.WithContext(ctx).Save(&user)
		return &user, nil
	}

	//fall back to email, user may have registered with email before
	err = s.db.WithContext(ctx).Raw("SELECT * FROM users WHERE email = ?", email).First(&user).Error

	if err == nil {
		//found by email, thenlink their GitHub account now
		user.GitHubID = id
		user.GithubToken = hashed_token
		s.db.WithContext(ctx).Save(&user)
		return &user, nil
	}

	//if doesnt exist
	if errors.Is(err, gorm.ErrRecordNotFound) {
		new_user := domain.User{Username: username, Email: email, Password: "", GitHubID: id, GithubToken: hashed_token}
		result := s.db.WithContext(ctx).Create(&new_user)
		if result.Error != nil {
			return nil, result.Error
		}

		return &new_user, nil
	}

	return nil, err

}

func (s *PostgreUserStore) LinkGithub(ctx context.Context, userID int, githubID int, githubToken string) error {
	return s.db.WithContext(ctx).
		Model(&domain.User{}).
		Where("id = ?", userID).
		Updates(map[string]any{
			"github_id":    githubID,
			"github_token": githubToken,
		}).Error
}

func (s *PostgreUserStore) DeleteUser(ctx context.Context, id int) error {
	return s.db.WithContext(ctx).Delete(&domain.User{}, id).Error
}
