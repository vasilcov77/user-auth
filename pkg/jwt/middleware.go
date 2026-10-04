package jwt

import (
	"context"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type UserIDContextKey struct{}

func Middleware(mgr *Manager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)

				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
				return mgr.AccessSecret(), nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			userID, err := uuid.Parse(claims["sub"].(string))
			if err != nil {
				http.Error(w, "Invalid user id", http.StatusUnauthorized)
			}

			ctx := context.WithValue(r.Context(), UserIDContextKey{}, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}
