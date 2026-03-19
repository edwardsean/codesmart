package domain

import "time"

type Level struct {
	ID                 int       `json:"id"          gorm:"primaryKey"`
	ProjectID          int       `json:"project_id"`
	LevelNumber        int       `json:"level_number"`
	Title              string    `json:"title"`
	Description        string    `json:"description"`
	LearningObjectives []string  `json:"learning_objectives" gorm:"type:text[]"`
	Concepts           []string  `json:"concepts"            gorm:"type:text[]"`
	StarterCode        string    `json:"starter_code,omitempty"`
	SolutionCode       string    `json:"-"           gorm:"column:solution_code"` // never expose to frontend
	TestCases          []byte    `json:"test_cases,omitempty"  gorm:"type:jsonb"`
	Hints              []byte    `json:"hints,omitempty"       gorm:"type:jsonb"`
	Difficulty         string    `json:"difficulty,omitempty"`
	EstimatedTime      int       `json:"estimated_time,omitempty"`
	Prerequisites      []int     `json:"prerequisites,omitempty" gorm:"type:integer[]"`
	Status             string    `json:"status"` // "locked", "unlocked", "completed"
	CreatedAt          time.Time `json:"created_at"`
}
