package domain

type UserService interface {
	GetUserFromToken(token string) (*User, error)
	GetUserByEmail(email string) (*User, error)
}

type GithubService interface {
	GetUserRepositories(user *User) (*[]Repository, error)
}

type OAuthService interface {
	HandleGithubCallback(code string) (string, error)
}

type AuthService interface {
	Login(payload LoginUserPayload) (string, string, error)
	Register(payload RegisterUserPayload) error
}
