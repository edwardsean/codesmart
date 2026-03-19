package postgres

import (
	"context"
	"errors"

	"github.com/edwardsean/codesmart/backend/internal/domain"
	"gorm.io/gorm"
)

type PostgresProjectStore struct {
	db *gorm.DB
}

func NewPostgresProjectStore(db *gorm.DB) *PostgresProjectStore {
	return &PostgresProjectStore{db: db}
}

func (s *PostgresProjectStore) GetProjectsByUserID(ctx context.Context, userID int) ([]domain.Project, error) {
	var projects []domain.Project

	err := s.db.WithContext(ctx).
		Where("user_id = ? AND status != ?", userID, "archived").
		Order("updated_at DESC").
		Find(&projects).Error

	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (s *PostgresProjectStore) GetProjectByID(ctx context.Context, id int) (*domain.Project, error) {
	var project domain.Project

	err := s.db.WithContext(ctx).
		First(&project, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}

	return &project, nil
}

func (s *PostgresProjectStore) CreateProject(ctx context.Context, project *domain.Project) error {
	return s.db.WithContext(ctx).Create(project).Error
}

func (s *PostgresProjectStore) DeleteProject(ctx context.Context, id int) error {
	// soft delete, set status to archived instead of hard deleting
	// this preserves the user's progress history
	return s.db.WithContext(ctx).
		Model(&domain.Project{}).
		Where("id = ?", id).
		Update("status", "archived").Error
}
