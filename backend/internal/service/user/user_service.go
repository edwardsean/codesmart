package user

import (
	stdError "errors"
	"strconv"

	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/pkg/jwt"
)

type UserService struct {
	store domain.UserRepository
}

func NewUserService(store domain.UserRepository) *UserService {
	return &UserService{store: store}
}

func (s *UserService) GetUserFromToken(token string) (*domain.User, error) {
	claims, err := jwt.GetTokenClaims(token)
	if err != nil {
		return nil, err
	}

	userIDstr, ok := claims["userID"].(string)
	if !ok {
		return nil, stdError.New("invalid user ID in token")
	}

	userId, _ := strconv.Atoi(userIDstr)

	return s.store.GetUserByID(userId)
}

func (s *UserService) GetUserByEmail(email string) (*domain.User, error) {
	user, err := s.store.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}
	return user, nil
}
