package middleware

import (
	"context"
	stdError "errors"
	"log"
	"net/http"
	"strings"

	"github.com/edwardsean/codesmart/backend/internal/domain"
	"github.com/edwardsean/codesmart/backend/internal/service"
	"github.com/edwardsean/codesmart/backend/pkg/errors"
	"github.com/edwardsean/codesmart/backend/pkg/response"
)

//http specific: http.HandlerFunc, http.ResponseWriter, http.Request

type contextKey string

const (
	UserKey  contextKey = "userID"
	TokenKey contextKey = "access_token"
)

// make it to return a handler function
func WithJWTAuth(service service.AuthService) func(http.HandlerFunc) http.HandlerFunc {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			//get the token from the user request
			access_token, err := getAccessTokenFromReq(r)
			if err != nil {
				// writeUnauthorizedError(w, stdError.New("permission denied"))
				writeUnauthorizedError(w, err)
				return
			}

			//validate token and fetch claims
			user, err := service.GetUserFromToken(r.Context(), access_token)
			if err != nil {
				writeUnauthorizedError(w, err)
				return
			}
			// claims, err := auth.GetTokenClaims(access_token)

			// if err != nil {
			// 	writeUnauthorizedError(w, errors.New(err.Error()))
			// 	return
			// }

			// // user, err := auth.GetUserFromClaims(claims, store)
			// if err != nil {
			// 	log.Printf("failed to get user by id: %v", err)
			// 	writeUnauthorizedError(w, fmt.Errorf("unable to get user from claims: %v", err))
			// 	return
			// }

			ctx := r.Context()
			ctx = context.WithValue(ctx, UserKey, user)
			ctx = context.WithValue(ctx, TokenKey, access_token)
			r = r.WithContext((ctx)) //This is necessary because http.Request is immutable (you can’t just change its context directly)

			next(w, r) //call the next handler
		}

	}
}

func getAccessTokenFromReq(r *http.Request) (string, error) {

	authHeader := r.Header.Get("Authorization")
	log.Println("Auth header:", authHeader)
	if authHeader == "" {
		return "", stdError.New("no authorization header")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", stdError.New("invalid authorization header")
	}

	return parts[1], nil
}

func GetUserFromContext(r *http.Request) (*domain.User, error) {
	user, ok := r.Context().Value(UserKey).(*domain.User)
	if !ok || user == nil {
		return nil, stdError.New("user not found in context")
	}
	return user, nil
}

func GetAccessTokenFromContext(r *http.Request) (string, error) {
	accessToken, ok := r.Context().Value(TokenKey).(string)
	if !ok || accessToken == "" {
		return "", stdError.New("access token not found in context")
	}
	return accessToken, nil
}

func writeUnauthorizedError(w http.ResponseWriter, message error) {
	// response.WriteError(w, http.StatusUnauthorized, message)
	response.WriteError(w, errors.NewError(message.Error(), http.StatusUnauthorized))
}
