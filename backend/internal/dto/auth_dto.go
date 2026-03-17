package dto

type RegisterUserPayload struct {
	Username        string `json:"username" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=6,max=13"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6,max=13"`
	GitHubID        int    `json:"github_id,omitempty"`
}

type LoginUserPayload struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AuthResponseDTO struct {
	AccessToken string          `json:"access_token"`
	User        UserResponseDTO `json:"user"`
}

type AuthResultDTO struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	User         UserResponseDTO `json:"user"`
}

type OAuthResultDTO struct {
	AccessToken  string          `json:"access_token"`
	RefreshToken string          `json:"refresh_token"`
	User         UserResponseDTO `json:"user"`
}
