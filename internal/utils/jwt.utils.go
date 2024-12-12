package utils

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// DecodeJWT decodes and validates a JWT token using the provided secret key.
// It checks the token's signing method, validates its structure, and returns
// the claims as a map if the token is valid or an error if it is not.
func DecodeJWT(token, secretKey string) (map[string]interface{}, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok || !parsedToken.Valid {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}
