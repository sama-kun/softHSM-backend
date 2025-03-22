package middleware

import (
	"context"
	"encoding/json"
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
			http.Error(w, "invalid token", http.StatusUnauthorized)
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

func ExtractAndDecryptSessionToken(r *http.Request) (*services.ClaimsService, error) {
	var requestBody struct {
		SessionToken string `json:"sessionToken"`
	}

	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		return nil, errors.New("failed to decode request body")
	}

	if requestBody.SessionToken == "" {
		return nil, errors.New("sessionToken is missing in request body")
	}

	claimsService := services.NewClaimsService(config.MustLoad())

	claims, err := claimsService.ValidateSessionToken(requestBody.SessionToken)
	if err != nil {
		return nil, errors.New("invalid session token")
	}

	return claims, nil
}
