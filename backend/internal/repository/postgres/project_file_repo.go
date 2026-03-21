package postgres

import (
	"context"

	"github.com/edwardsean/codesmart/backend/internal/domain"
	"gorm.io/gorm"
)

type PostgresProjectFileStore struct {
	db *gorm.DB
}

func NewPostgresProjectFileStore(db *gorm.DB) *PostgresProjectFileStore {
	return &PostgresProjectFileStore{db: db}
}

func (s *PostgresProjectFileStore) GetFileByID(ctx context.Context, id int) (*domain.ProjectFile, error) {
	var file domain.ProjectFile
	err := s.db.WithContext(ctx).First(&file, id).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (s *PostgresProjectFileStore) GetFiles(ctx context.Context, projectID int) ([]domain.ProjectFile, error) {
	var files []domain.ProjectFile
	err := s.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("file_path ASC").
		Find(&files).Error
	return files, err
}

func (s *PostgresProjectFileStore) GetFileContent(ctx context.Context, projectID int, path string) (*domain.ProjectFile, error) {
	var file domain.ProjectFile
	err := s.db.WithContext(ctx).
		Where("project_id = ? AND file_path = ?", projectID, path).
		First(&file).Error
	if err != nil {
		return nil, err
	}
	return &file, nil
}

func (s *PostgresProjectFileStore) CreateFile(ctx context.Context, file *domain.ProjectFile) error {
	return s.db.WithContext(ctx).Create(file).Error
}

func (s *PostgresProjectFileStore) UpdateFile(ctx context.Context, file *domain.ProjectFile) error {
	return s.db.WithContext(ctx).Save(file).Error
}

func (s *PostgresProjectFileStore) DeleteFile(ctx context.Context, id int) error {
	return s.db.WithContext(ctx).Delete(&domain.ProjectFile{}, id).Error
}
