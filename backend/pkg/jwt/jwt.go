package jwt

import (
	"fmt"
	"strconv"
	"time"

	"github.com/edwardsean/codesmart/backend/internal/config"
	"github.com/golang-jwt/jwt"
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

// func GetUserFromClaims(claims jwt.MapClaims, store domain.UserStore) (*domain.User, error) {
// 	str, ok := claims["userID"].(string)
// 	if !ok {
// 		return nil, errors.New("user id interface is null, not string")
// 	}

// 	userId, _ := strconv.Atoi(str)

// 	user, err := store.GetUserByID(userId) //service layer should work with interfaces not concrete types UserStore
// 	if err != nil {
// 		return nil, err
// 	}

// 	return user, nil
// }
