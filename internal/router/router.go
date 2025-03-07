package router

import (
	"fmt"
	"net/http"
	"soft-hsm/internal/auth"
	"soft-hsm/internal/auth/handlers"
	authRepository "soft-hsm/internal/auth/repository"
	"soft-hsm/internal/auth/services"
	"soft-hsm/internal/config"
	"soft-hsm/internal/mailer"
	"soft-hsm/internal/middleware"
	"soft-hsm/internal/storage"
	userRepository "soft-hsm/internal/user/repository"

	"github.com/go-chi/chi/v5"
)

func SetupRouter(cfg *config.Config, db *storage.Postgres, redis *storage.Redis) *chi.Mux {
	r := chi.NewRouter()

	tokenRepo := authRepository.NewTokenRepository(redis)
	claimsService := services.NewCliamsService(cfg)
	userRepo := userRepository.NewUserRepository(db)

	mailerService := mailer.NewMailer(&cfg.MailerConfig)

	passwordService := services.NewPasswordService()

	authService := services.NewAuthService(tokenRepo, claimsService, userRepo, mailerService, passwordService)
	authHandler := handlers.NewAuthHandler(authService)

	r.Mount("/api/v1/auth", auth.AuthRoutes(authHandler))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("Server is running!")); err != nil {
			http.Error(w, "Failed to write response", http.StatusInternalServerError)
		}
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		middleware.ErrorHandler(w, http.StatusNotFound, fmt.Errorf("route not found"), nil)
	})

	return r
}
