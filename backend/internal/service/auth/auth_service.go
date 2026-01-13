package auth

import (
	"fmt"
	"net/http"
	"time"

	stdError "errors"

	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/jwt"
	"github.com/edwardsean/codesmart/backend/pkg/password"
	"github.com/edwardsean/codesmart/backend/pkg/validator"
)

type AuthService struct {
	userStore domain.UserRepository
}

func NewAuthService(userStore domain.UserRepository) *AuthService {
	return &AuthService{
		userStore: userStore,
	}
}

func (s *AuthService) Login(payload domain.LoginUserPayload) (string, string, error) {
	if err := validator.Validate.Struct(payload); err != nil {
		return "", "", stdError.New("invalid payload")
	}

	user, err := s.userStore.GetUserByEmail(payload.Email)
	if err != nil {
		// response.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return "", "", errors.ErrInvalidCredentials
	}

	if !password.ComparePassword(user.Password, []byte(payload.Password)) {
		return "", "", errors.ErrInvalidCredentials
	}

	secret := []byte(config.Envs.JWTSecret)
	access_token, err := jwt.CreateJWT(secret, user.ID, 15*time.Minute)

	if err != nil {
		return "", "", errors.ErrTokenGeneration
	}

	refresh_token, err := jwt.CreateJWT(secret, user.ID, 7*24*time.Hour)
	if err != nil {
		return "", "", errors.ErrTokenGeneration
	}

	return access_token, refresh_token, nil
}

func (s *AuthService) Register(payload domain.RegisterUserPayload) error {
	// validate payload
	if err := validator.Validate.Struct(payload); err != nil {
		// errors := err.(validator.ValidationErrors)
		return errors.ErrInvalidPayload
	}

	//check if user exists
	user, err := s.userStore.GetUserByEmail(payload.Email)
	if err == nil {
		return errors.NewError(fmt.Sprintf("user with email %s already exists, response: %v", payload.Email, user), http.StatusBadRequest)
	}

	//if it doesnt, we create the new user
	//hashpassword
	hash_password, err := password.HashPassword(payload.Password)

	if err != nil {
		return errors.NewError(err.Error(), http.StatusInternalServerError)
	}

	err = s.userStore.CreateUser(domain.User{
		Email:    payload.Email,
		Username: payload.Username,
		Password: hash_password,
	})

	if err != nil {
		return errors.NewError(err.Error(), http.StatusInternalServerError)
	}

	return nil
}
