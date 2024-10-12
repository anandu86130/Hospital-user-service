package token

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/anandu86130/Hospital-user-service/internal/model"
	"github.com/dgrijalva/jwt-go"
)

// GenerateToken generates a token for 5 hours with given data
func GenerateToken(email string, userid uint) (string, error) {
	// Get the secret key from environment variable
	secretKey := os.Getenv("JWT_SECRET_KEY")
	if secretKey == "" {
		return "", errors.New("JWT_SECRET_KEY not set in environment")
	}

	// Create the claims
	claims := &model.UserClaims{
		UserID: userid,
		Email:  email,
		Role:   "user",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 10).Unix(), // Token expires after 10 hours
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", fmt.Errorf("failed to sign token: %w", err)
	}

	return tokenString, nil
}
