package dto

import (
	"time"

	"github.com/edwardsean/codesmart/backend/internal/domain"
)

type UserResponseDTO struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	GithubID  int       `json:"github_id,omitempty"`
	CreatedAt time.Time `json:"createdAt"`
}

func ToUserResponseDTO(user *domain.User) UserResponseDTO {
	return UserResponseDTO{
		ID:        user.ID,
		Email:     user.Email,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		GithubID:  user.GitHubID,
	}
}
