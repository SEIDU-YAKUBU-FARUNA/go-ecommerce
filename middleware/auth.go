package middleware

import (
	"net/http"
	"strings"

	"go-ecommerce/config"

	"github.com/golang-jwt/jwt/v5"
)

// AuthMiddleware verifies JWT token
func AuthMiddleware(next http.HandlerFunc, requireAdmin bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing token", http.StatusUnauthorized)
			return
		}

		// Expect: Bearer TOKEN
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 {
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		tokenString := parts[1]
		secret := config.GetEnv("JWT_SECRET")

		// Parse token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Check admin if required
		if requireAdmin {
			claims := token.Claims.(jwt.MapClaims)
			if claims["is_admin"] != true {
				http.Error(w, "Admin access required", http.StatusForbidden)
				return
			}
		}

		next(w, r)
	}
}
