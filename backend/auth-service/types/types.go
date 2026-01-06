package types

import (
	"time"
)

// an interface value is a 2-word pair:
// •	a pointer to type information (which tells Go what concrete type it holds)
// •	a pointer to the actual data value (the thing you’re working with)
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
	GetOrCreateUserFromGithub(id int, email string, login string, access_token string, github_user *GithubUser) (*User, error)
}

type RepositoryStore interface {
}
type RegisterUserPayload struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6,max=13"`
	GitHubID int    `json:"github_id,omitempty"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

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

type Repository struct {
	ID              int       `json:"id"`
	Name            string    `json:"name"`
	FullName        string    `json:"full_name"`
	Description     string    `json:"description"`
	HTMLURL         string    `json:"html_url"`
	Language        string    `json:"language"`
	StargazersCount int       `json:"stargazers_count"`
	ForksCount      int       `json:"forks_count"`
	UpdatedAt       string    `json:"updated_at"`
	Private         bool      `json:"private"`
	DefaultBranch   string    `json:"default_branch"`
	Owner           RepoOwner `json:"owner"`
}

type RepoOwner struct {
	Login     string `json:"login"`
	AvatarURL string `json:"avatar_url"`
}
