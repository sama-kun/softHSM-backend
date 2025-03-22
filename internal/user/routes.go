package user

import (
	"soft-hsm/internal/user/handlers"

	"github.com/go-chi/chi"
)

func UserRoutes(userHandler *handlers.UserHandler) chi.Router {
	r := chi.NewRouter()

	r.Post("/me", userHandler.Me)
	return r
}
