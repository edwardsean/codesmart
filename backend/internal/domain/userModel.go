package domain

import "time"

type User struct {
	ID          int       `json:"id"          gorm:"primaryKey"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Password    string    `json:"-"           gorm:"column:password_hash"`
	GitHubID    int       `json:"github_id,omitempty" gorm:"column:github_id"`
	GithubToken string    `json:"-"           gorm:"column:github_token"` //used by github_handler to call github for the user
	XPPoints    int       `json:"xp_points"`                              //Total XP earned across all levels, drives the gamification system
	Level       int       `json:"level"`                                  //User's current level calculated from XP, shown in profile/dashboard
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"` //Last time any field changed, useful for "last active" tracking
}

type GithubUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}
