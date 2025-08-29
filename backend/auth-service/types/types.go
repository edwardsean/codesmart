package types

import "time"

// an interface value is a 2-word pair:
// •	a pointer to type information (which tells Go what concrete type it holds)
// •	a pointer to the actual data value (the thing you’re working with)
type UserStore interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
	GetOrCreateUserFromGithub(id int, email string, login string, access_token string, github_user *GithubUser) (*User, error)
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
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	GitHubID  int       `json:"github_id,omitempty" gorm:"column:github_id"`
	CreatedAt time.Time `json:"createdAt"`
}

type GithubUser struct {
	ID    int    `json:"id"`
	Login string `json:"login"`
	Email string `json:"email"`
}
