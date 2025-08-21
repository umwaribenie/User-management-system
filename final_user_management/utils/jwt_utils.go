// utils/jwt_utils.go (Updated)
package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtKey []byte

// Claims defines the JWT claims, including user-specific data and standard claims.
type Claims struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// SetJWTSecret sets the JWT secret key from an external source (like an environment variable).
func SetJWTSecret(secret string) {
	jwtKey = []byte(secret)
}

// GenerateJWT creates a new JWT token for a given user.
func GenerateJWT(userID, username, role string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	claims := &Claims{
		UserID:   userID,
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ValidateJWT parses and validates a JWT token string.
func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrSignatureInvalid
	}

	return claims, nil
}
