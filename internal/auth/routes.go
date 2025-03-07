package auth

import (
	"soft-hsm/internal/auth/handlers"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(authHandler *handlers.AuthHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/register", authHandler.Register)
  r.Post("/login", authHandler.Login)
	return r
}
