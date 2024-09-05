package main

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Secret key to sign JWT tokens (use environment variables in real apps)
var jwtKey = []byte("my_secret_key")

// GenerateJWTToken generates a new JWT token with an expiration time
func GenerateJWTToken(username string) (string, error) {
	// Set expiration time for the token (e.g., 1 hour)
	expirationTime := time.Now().Add(1 * time.Hour)

	// Create the JWT claims, which includes the username and expiry time
	claims := &jwt.MapClaims{
		"username": username,
		"exp":      expirationTime.Unix(),
	}

	// Create the JWT token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateJWTToken parses and validates a JWT token
func ValidateJWTToken(tokenString string) (*jwt.MapClaims, error) {
	// Parse the JWT token
	claims := &jwt.MapClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	// If token is invalid or parsing fails
	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired token")
	}

	return claims, nil
}
