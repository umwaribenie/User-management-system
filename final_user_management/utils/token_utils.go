package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// This secret should be different (and rotated independently) from your auth JWT secret.
var passwordResetSecret = []byte("your-password-reset-secret-here")

// GeneratePasswordResetToken issues a JWT whose subject is the user's email,
// expiring in 1 hour. You embed this in the reset link.
func GeneratePasswordResetToken(email string) (string, error) {
	claims := jwt.RegisteredClaims{
		Subject:   email,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(passwordResetSecret)
}

// VerifyPasswordResetToken checks the token and returns the email if valid.
func VerifyPasswordResetToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		return passwordResetSecret, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims.Subject, nil
	}
	return "", jwt.ErrTokenInvalidClaims
}
