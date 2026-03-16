/**

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


**/

package middleware

import (
	"github.com/golang-jwt/jwt/v5"
	"go-ecommerce/utils"
	"net/http"
	"os"
	"strings"
)

func AuthMiddleware(next http.HandlerFunc, adminOnly bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Missing Authorization Header")
			return
		}

		// Remove "Bearer " from the start of the token string
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid or expired token")
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid token claims")
			return
		}

		// THE FIX: Use "isAdmin" (matches the JSON tag in utils/jwt.go)
		isAdmin, ok := claims["is_admin"].(bool)

		if adminOnly && (!ok || !isAdmin) {
			utils.RespondWithError(w, http.StatusForbidden, "Admin access required")
			return
		}

		next.ServeHTTP(w, r)
	}
}
