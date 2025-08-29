package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/edwardsean/codesmart/backend/auth-service/config"
	"github.com/edwardsean/codesmart/backend/auth-service/types"
	"github.com/edwardsean/codesmart/backend/auth-service/utils"
	"github.com/golang-jwt/jwt"
)

type contextKey string

const UserKey contextKey = "userID"

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
		cookie := getAccessTokenFromReq(r)
		if cookie == nil {
			permissionDenied(w)
			return
		}

		//validate the JWT
		token, err := validateToken(cookie.Value)

		if err != nil {
			log.Printf("failed to validate token: %v", token)
			permissionDenied(w)
			return
		}

		if !token.Valid {
			log.Println("invalid token")
			permissionDenied(w)
			return
		}

		//fetch user from db
		claims := token.Claims.(jwt.MapClaims)
		str := claims["userId"].(string)

		userId, _ := strconv.Atoi(str)

		user, err := store.GetUserByID(userId)
		if err != nil {
			log.Printf("failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user.ID)
		r = r.WithContext((ctx)) //This is necessary because http.Request is immutable (you can’t just change its context directly)

		handlerFunc(w, r)

	}

}

func getAccessTokenFromReq(r *http.Request) *http.Cookie {
	token, err := r.Cookie("access_token")
	if err != nil || token.Value == "" {
		return nil
	}

	return token
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
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
