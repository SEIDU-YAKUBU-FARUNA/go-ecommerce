/**

package utils

import (
	"time"

	"go-ecommerce/config"


	"github.com/golang-jwt/jwt/v5"
)

// GenerateToken creates a signed JWT for a user
func GenerateToken(userID string, isAdmin bool) (string, error) {

	// Create token claims (data inside token)
	claims := jwt.MapClaims{
		"user_id":  userID,
		"is_admin": isAdmin,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // expires in 24 hours
	}

	// Create token object
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Get secret from env
	secret := config.GetEnv("JWT_SECRET")

	// Sign the token
	return token.SignedString([]byte(secret))
}


**/

package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type Claims struct {
	UserID  string `json:"user_id"`
	IsAdmin bool   `json:"isAdmin"` // This is the key that goes into the JWT
	jwt.RegisteredClaims
}

func GenerateToken(userID string, isAdmin bool) (string, error) {
	claims := &Claims{
		UserID:  userID,
		IsAdmin: isAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}
