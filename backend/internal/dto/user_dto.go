package dto

import "time"

type UserResponseDTO struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	GithubID  int       `json:"github_id,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}
