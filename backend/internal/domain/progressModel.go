package domain

import "time"

type UserLevelProgress struct {
	ID                  int        `json:"id"          gorm:"primaryKey"`
	UserID              int        `json:"user_id"`
	LevelID             int        `json:"level_id"`
	CurrentCode         string     `json:"current_code,omitempty" gorm:"type:text"`
	Status              string     `json:"status"` // "not_started", "in_progress", "completed"
	Attempts            int        `json:"attempts"`
	HintsUsed           int        `json:"hints_used"`
	HelpRequests        int        `json:"help_requests"`
	TimeSpent           int        `json:"time_spent"` // seconds
	LastExecutionStatus string     `json:"last_execution_status,omitempty"`
	LastError           string     `json:"last_error,omitempty"`
	CompletedAt         *time.Time `json:"completed_at,omitempty"` // pointer — nullable
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
