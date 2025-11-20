package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/Nikitaannusewicz/carwash-crm/internal/middleware"
	"github.com/golang-jwt/jwt/v5"
)

// AI NOTE: AuthMiddleWare checks for a valid JWT and injects claims into the context.
// It returns a "Constructor" for the middleware because we need to inject the secret.
func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenString string

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				tokenString = parts[1]
			}
		}

		if tokenString == "" {
			tokenString = r.Header.Get("X-Auth-Token")
		}

		if tokenString == "" {
			if cookie, err := r.Cookie("auth-token"); err == nil {
				tokenString = cookie.Value
			}
		}

		if tokenString == "" {
			http.Error(w, "Authentication token required", http.StatusUnauthorized)
			return
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(s.cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid token claims", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), middleware.UserIDKey, claims["sub"])
		ctx = context.WithValue(ctx, middleware.RoleKey, claims["role"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
