package repository

import "gorm.io/gorm"

type PostgreUserStore struct {
	db *gorm.DB
}

func PostgreNewStore(db *gorm.DB) *PostgreUserStore {
	return &PostgreUserStore{db: db}
}
