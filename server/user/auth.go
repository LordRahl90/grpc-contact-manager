package user

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// generateToken generates the JWT token for the given user
func generateToken(userID uint32, expired bool) (string, error) {
	expiry := time.Now().Add(24 * time.Hour)
	if expired {
		expiry = time.Now().Add(-24 * time.Hour)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"nbf":     expiry,
	})

	return token.SignedString(signingSecret)
}

// validateToken this should validate the token string and return the userID
func validateToken(tokenString string) (uint32, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return signingSecret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, errInvalidToken
	}
	exp := claims["nbf"].(string)
	expiryDate, err := time.Parse(time.RFC3339, exp)
	if err != nil {
		return 0, err
	}

	if expiryDate.Before(time.Now()) {
		return 0, errTokenExpired
	}

	return uint32(claims["user_id"].(float64)), nil
}
