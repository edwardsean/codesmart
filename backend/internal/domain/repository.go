package domain

type UserRepository interface {
	GetUserByEmail(email string) (*User, error)
	GetUserByID(id int) (*User, error)
	CreateUser(User) error
	GetOrCreateUserFromGithub(id int, email string, login string, access_token string, github_user *GithubUser) (*User, error)
}
