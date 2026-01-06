package auth

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/edwardsean/codesmart/backend/auth-service/config"
	"github.com/edwardsean/codesmart/backend/auth-service/types"
	"github.com/edwardsean/codesmart/backend/auth-service/utils"
	"github.com/golang-jwt/jwt"
)

type contextKey string

const (
	UserKey  contextKey = "userID"
	TokenKey contextKey = "access_token"
)

func CreateJWT(secret []byte, userID int, duration time.Duration) (string, error) {
	// expiration := time.Second * time.Duration(config.Envs.JWTExpirationInSeconds)
	expiration := time.Now().Add(duration).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiredAt": expiration,
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func WithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//get the token from the user request
		access_token := getAccessTokenFromReq(r)
		if access_token == "" {
			writeUnauthorizedError(w, errors.New("missing access token"))
			return
		}

		//validate token and fetch claims
		claims, err := GetTokenClaims(access_token)

		if err != nil {
			writeUnauthorizedError(w, errors.New(err.Error()))
			return
		}

		user, err := GetUserFromClaims(claims, store)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			writeUnauthorizedError(w, fmt.Errorf("unable to get user from claims: %v", err))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user)
		ctx = context.WithValue(ctx, TokenKey, access_token)
		r = r.WithContext((ctx)) //This is necessary because http.Request is immutable (you can’t just change its context directly)

		handlerFunc(w, r)

	}

}

func getAccessTokenFromReq(r *http.Request) string {

	authHeader := r.Header.Get("Authorization")
	log.Println("Auth header:", authHeader)
	if authHeader == "" {
		return ""
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}

	return parts[1]
}

func writeUnauthorizedError(w http.ResponseWriter, message error) {
	utils.WriteError(w, http.StatusUnauthorized, message)
}

func validateToken(tokenString string) (*jwt.Token, error) {
	//this is how jwt.Parse checks if the token’s signature matches the one generated with your JWTSecret.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(config.Envs.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return token, nil

}

func GetTokenClaims(tokenString string) (jwt.MapClaims, error) {
	token, err := validateToken(tokenString)

	if err != nil {
		return nil, fmt.Errorf("failed to validate token: %v", token)
	}

	//fetch user from db
	claims, ok := token.Claims.(jwt.MapClaims)

	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if exp, ok := claims["expiredAt"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return nil, fmt.Errorf("permission denied")
		}
	}

	return claims, nil
}

func GetUserFromClaims(claims jwt.MapClaims, store types.UserStore) (*types.User, error) {
	str, ok := claims["userID"].(string)
	if !ok {
		return nil, errors.New("user id interface is null, not string")
	}

	userId, _ := strconv.Atoi(str)

	user, err := store.GetUserByID(userId)
	if err != nil {
		return nil, err
	}

	return user, nil
}
