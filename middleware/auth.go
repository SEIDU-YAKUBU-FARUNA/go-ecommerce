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
	"net/http"
	"os"
	"strings"

	"go-ecommerce/utils"

	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(next http.HandlerFunc, adminOnly bool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Missing Authorization Header")
			return
		}

		// 1. MUST strip the "Bearer " prefix
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("JWT_SECRET")), nil
		})

		if err != nil || !token.Valid {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid Token")
			return
		}

		// 2. Extract Claims carefully
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			utils.RespondWithError(w, http.StatusUnauthorized, "Invalid Claims")
			return
		}

		// 3. THE FIX: The key here must match your GenerateToken key exactly
		// If jwt.io showed "isAdmin", use "isAdmin" here.
		isAdmin, ok := claims["isAdmin"].(bool)

		if adminOnly && (!ok || !isAdmin) {
			utils.RespondWithError(w, http.StatusForbidden, "Admin access required")
			return
		}

		next.ServeHTTP(w, r)
	}
}
