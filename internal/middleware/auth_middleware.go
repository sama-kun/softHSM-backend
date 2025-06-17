package middleware

import (
	"context"
	"errors"
	"net/http"
	"soft-hsm/internal/auth/services"
	"soft-hsm/internal/config"
	"strings"
)

type ContextKey string

const UserKey ContextKey = "user"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := extractToken(r)
		if err != nil {
			ErrorHandler(w, http.StatusUnauthorized, err, "Unauthorized")
			return
		}
		claims, err := services.NewClaimsService(config.MustLoad()).ValidateToken(token)

		if err != nil {
			ErrorHandler(w, http.StatusUnauthorized, err, "invalid token")
		}

		ctx := context.WithValue(r.Context(), UserKey, claims)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
func extractToken(r *http.Request) (string, error) {
	authHeader := r.Header.Get("Authorization")

	if authHeader == "" {
		return "", errors.New("missing Auth header")
	}

	parts := strings.Split(authHeader, " ")

	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid Auth header format")
	}

	return parts[1], nil
}

func GetUserFromContext(r *http.Request) (*services.ClaimsService, error) {
	user, ok := r.Context().Value(UserKey).(*services.ClaimsService)

	if !ok {
		return nil, errors.New("user not found from context")
	}
	return user, nil
}

// func ExtractAndDecryptSessionToken(r *http.Request) (*services.ClaimsService, error) {
// 	token := r.Header.Get("X-Session-Token")
// 	if token == "" {
// 		return nil, errors.New("missing session token")
// 	}

// 	fmt.Println(token)

// 	claimsService := services.NewClaimsService(config.MustLoad())
// 	// claims, err := claimsService.ValidateSessionToken(token)
// 	if err != nil {
// 		return nil, errors.New("invalid session token")
// 	}

// 	return claims, nil
// }
