package auth

import (
	"context"
	"fmt"
	"net/http"
	"time"

	stdError "errors"

	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/internal/dto"
	"github.com/edwardsean/codesmart/backend/internal/repository"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/jwt"
	"github.com/edwardsean/codesmart/backend/pkg/password"
	"github.com/edwardsean/codesmart/backend/pkg/validator"
)

type AuthService struct {
	userRepo  repository.UserRepository
	tokenRepo repository.TokenRepository
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository) *AuthService {
	return &AuthService{
		userRepo: userRepo, tokenRepo: tokenRepo,
	}
}

func (s *AuthService) Login(ctx context.Context, payload dto.LoginUserPayload) (string, string, error) {
	if err := validator.Validate.Struct(payload); err != nil {
		return "", "", stdError.New("invalid payload")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, payload.Email)
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

func (s *AuthService) Register(ctx context.Context, payload dto.RegisterUserPayload) error {
	// validate payload
	if err := validator.Validate.Struct(payload); err != nil {
		// errors := err.(validator.ValidationErrors)
		return errors.ErrInvalidPayload
	}

	//check if user exists
	user, err := s.userRepo.GetUserByEmail(ctx, payload.Email)
	if err == nil {
		return errors.NewError(fmt.Sprintf("user with email %s already exists, response: %v", payload.Email, user), http.StatusBadRequest)
	}

	//if it doesnt, we create the new user
	//hashpassword
	hash_password, err := password.HashPassword(payload.Password)

	if err != nil {
		return errors.NewError(err.Error(), http.StatusInternalServerError)
	}

	err = s.userRepo.CreateUser(ctx, &domain.User{
		Email:    payload.Email,
		Username: payload.Username,
		Password: hash_password,
	})

	if err != nil {
		return errors.NewError(err.Error(), http.StatusInternalServerError)
	}

	return nil
}

func (s *AuthService) Logout(ctx context.Context, refreshToken string) error {
	claims, err := jwt.GetTokenClaims(refreshToken)
	if err != nil {
		return nil // token already invalid, treat as successful logout
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		return errors.NewError("Invalid token claims", http.StatusBadRequest)
	}

	ttl := time.Until(time.Unix(int64(exp), 0))
	if ttl <= 0 {
		return nil //already expired no need to store
	}

	return s.tokenRepo.BlacklistToken(ctx, refreshToken, ttl)
}

func (s *AuthService) ValidateRefreshToken(ctx context.Context, token string) (*domain.User, error) {
	//check blacklist BEFORE anything else
	blacklisted, err := s.tokenRepo.IsBlacklisted(ctx, token)
	if err != nil {
		return nil, err
	}
	if blacklisted {
		return nil, errors.NewError("token has been revoked", http.StatusUnauthorized)
	}

	return s.GetUserFromToken(ctx, token)
}

func (s *AuthService) GetUserFromToken(ctx context.Context, token string) (*domain.User, error) {
	userID, err := jwt.GetUserIDFromToken(token)
	if err != nil {
		return nil, err
	}

	return s.userRepo.GetUserByID(ctx, userID)
}
