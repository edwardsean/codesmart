package domain

import "time"

type Achievement struct {
	ID              int       `json:"id"          gorm:"primaryKey"`
	UserID          int       `json:"user_id"`
	AchievementType string    `json:"achievement_type"`
	Title           string    `json:"title"`
	Description     string    `json:"description,omitempty"`
	IconURL         string    `json:"icon_url,omitempty"`
	XPAwarded       int       `json:"xp_awarded"`
	Metadata        []byte    `json:"metadata,omitempty" gorm:"type:jsonb"`
	EarnedAt        time.Time `json:"earned_at"`
}
