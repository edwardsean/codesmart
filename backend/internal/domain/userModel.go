package domain

import "time"

type User struct {
	ID          int       `json:"id"`
	Email       string    `json:"email"`
	Username    string    `json:"username"`
	Password    string    `json:"-"`
	GitHubID    int       `json:"github_id,omitempty" gorm:"column:github_id"`
	GithubToken string    `json:"github_token,omitempty" gorm:"column:github_token"`
	CreatedAt   time.Time `json:"createdAt"`
}

type SafeUser struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	GitHubID  int       `json:"github_id,omitempty" gorm:"column:github_id"`
	CreatedAt time.Time `json:"createdAt"`
}

type GithubUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}
