package domain

import "time"

type AIConversation struct {
	ID               int       `json:"id"          gorm:"primaryKey"`
	UserID           int       `json:"user_id"`
	ProjectID        *int      `json:"project_id,omitempty"` // pointer — nullable
	LevelID          *int      `json:"level_id,omitempty"`   // pointer — nullable
	ConversationType string    `json:"conversation_type"`    // "hint", "chat", "help"
	Messages         []byte    `json:"messages"    gorm:"type:jsonb"`
	Context          []byte    `json:"context,omitempty" gorm:"type:jsonb"`
	TokensUsed       int       `json:"tokens_used,omitempty"`
	ModelUsed        string    `json:"model_used,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
}

type CodeExecution struct {
	ID            int       `json:"id"          gorm:"primaryKey"`
	UserID        int       `json:"user_id"`
	LevelID       *int      `json:"level_id,omitempty"` // pointer — nullable
	Code          string    `json:"code"        gorm:"type:text"`
	Language      string    `json:"language"`
	Output        string    `json:"output,omitempty"    gorm:"type:text"`
	Error         string    `json:"error,omitempty"     gorm:"type:text"`
	ExecutionTime int       `json:"execution_time,omitempty"` // milliseconds
	MemoryUsed    int       `json:"memory_used,omitempty"`    // bytes
	Status        string    `json:"status,omitempty"`         // "success", "error", "timeout"
	TestResults   []byte    `json:"test_results,omitempty" gorm:"type:jsonb"`
	CreatedAt     time.Time `json:"created_at"`
}
