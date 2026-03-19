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

func (s *AuthService) Login(ctx context.Context, payload dto.LoginUserPayload) (*dto.AuthResultDTO, error) {
	if err := validator.Validate.Struct(payload); err != nil {
		return nil, stdError.New("invalid payload")
	}

	user, err := s.userRepo.GetUserByEmail(ctx, payload.Email)
	if err != nil {
		// response.WriteError(w, http.StatusBadRequest, fmt.Errorf("not found, invalid email or password"))
		return nil, errors.ErrInvalidCredentials
	}

	if !password.ComparePassword(user.Password, []byte(payload.Password)) {
		return nil, errors.ErrInvalidCredentials
	}

	secret := []byte(config.Envs.JWTSecret)
	access_token, err := jwt.CreateJWT(secret, user.ID, 15*time.Minute)

	if err != nil {
		return nil, errors.ErrTokenGeneration
	}

	refresh_token, err := jwt.CreateJWT(secret, user.ID, 7*24*time.Hour)
	if err != nil {
		return nil, errors.ErrTokenGeneration
	}

	userResponse := dto.ToUserResponseDTO(user)

	return &dto.AuthResultDTO{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		User:         userResponse,
	}, nil
}

func (s *AuthService) Register(ctx context.Context, payload dto.RegisterUserPayload) (*dto.AuthResultDTO, error) {
	// validate payload
	if err := validator.Validate.Struct(payload); err != nil {
		// errors := err.(validator.ValidationErrors)
		return nil, errors.ErrInvalidPayload
	}

	//check confirm password and password match
	if payload.ConfirmPassword != payload.Password {
		return nil, errors.NewError("passwords do not match", http.StatusBadRequest)
	}

	//check if user exists
	user, err := s.userRepo.GetUserByEmail(ctx, payload.Email)
	if err == nil {
		return nil, errors.NewError(fmt.Sprintf("user with email %s already exists", payload.Email), http.StatusBadRequest)
	}

	//if it doesnt, we create the new user
	//hashpassword
	hash_password, err := password.HashPassword(payload.Password)

	if err != nil {
		return nil, errors.NewError(err.Error(), http.StatusInternalServerError)
	}

	user, err = s.userRepo.CreateUser(ctx, &domain.User{
		Email:    payload.Email,
		Username: payload.Username,
		Password: hash_password,
	})

	if err != nil {
		return nil, errors.NewError(err.Error(), http.StatusInternalServerError)
	}

	secret := []byte(config.Envs.JWTSecret)
	access_token, err := jwt.CreateJWT(secret, user.ID, 15*time.Minute)

	if err != nil {
		return nil, errors.ErrTokenGeneration
	}

	refresh_token, err := jwt.CreateJWT(secret, user.ID, 7*24*time.Hour)
	if err != nil {
		return nil, errors.ErrTokenGeneration
	}

	userResponse := dto.ToUserResponseDTO(user)

	return &dto.AuthResultDTO{
		AccessToken:  access_token,
		RefreshToken: refresh_token,
		User:         userResponse,
	}, nil
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
