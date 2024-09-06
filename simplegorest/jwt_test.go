package main

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

type JWTClaim struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func TestGenerateJWTToken(t *testing.T) {
	// Define test data
	username := "testuser"
	secretKey := []byte("mysecretkey")

	// Generate a JWT token
	token, err := GenerateJWTToken(username, secretKey)
	assert.NoError(t, err, "Expected no error when generating JWT token")
	assert.NotEmpty(t, token, "Expected a non-empty JWT token")

	// Validate the generated token to ensure it's valid
	claims, err := ValidateJWTToken(token, secretKey)
	assert.NoError(t, err, "Expected no error when validating JWT token")
	// Dereference the claims pointer and access the username
	assert.Equal(t, username, (*claims)["username"], "Expected the username in claims to match the input username")
}

func TestValidateJWTToken_ExpiredToken(t *testing.T) {
	// Create an expired token manually
	secretKey := []byte("mysecretkey")
	expiredToken := jwt.NewWithClaims(jwt.SigningMethodHS256, &jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)), // Expired one hour ago
		Subject:   "testuser",
	})
	tokenString, err := expiredToken.SignedString(secretKey)
	assert.NoError(t, err, "Expected no error when signing expired token")

	// Validate the expired token
	_, err = ValidateJWTToken(tokenString, secretKey)
	assert.Error(t, err, "Expected an error when validating an expired token")
}

func TestValidateJWTToken_InvalidToken(t *testing.T) {
	// Use an invalid token
	invalidToken := "this.is.an.invalid.token"
	secretKey := []byte("mysecretkey")

	// Validate the invalid token
	_, err := ValidateJWTToken(invalidToken, secretKey)
	assert.Error(t, err, "Expected an error when validating an invalid token")
}
