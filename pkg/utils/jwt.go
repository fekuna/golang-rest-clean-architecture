package utils

import (
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/fekuna/api-mc/config"
	"github.com/fekuna/api-mc/internal/models"
)

// JWT Claims struct
type Claims struct {
	Email string `json:"email"`
	ID    string `json:"id"`
	jwt.StandardClaims
}

// Generate new JWT Token
func GenerateJWTToken(user *models.User, config *config.Config) (string, error) {
	// Register the JWT claims, which includes the username and expiry time
	claims := &Claims{
		Email: user.Email,
		ID:    user.UserID.String(),
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Minute * 60).Unix(),
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Register the JWT string
	tokenString, err := token.SignedString([]byte(config.Server.JwtSecretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
