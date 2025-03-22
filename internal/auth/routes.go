package auth

import (
	"soft-hsm/internal/auth/handlers"
	"soft-hsm/internal/middleware"

	"github.com/go-chi/chi/v5"
)

func AuthRoutes(authHandler *handlers.AuthHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/register", authHandler.Register)
	r.Post("/login", authHandler.Login)
	r.Patch("/activate", authHandler.Activate)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Patch("/set-master", authHandler.SetMasterPassword)
		r.Post("/check-master", authHandler.CheckMasterPassword)
	})

	return r
}
