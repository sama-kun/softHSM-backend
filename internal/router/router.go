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
	userHandler "soft-hsm/internal/user/handlers"
	userRepository "soft-hsm/internal/user/repository"
	userService "soft-hsm/internal/user/services"

	blockchainKey "soft-hsm/internal/blockchain-key"
	blockchainHandler "soft-hsm/internal/blockchain-key/handlers"
	blockchainRepository "soft-hsm/internal/blockchain-key/repository"
	"soft-hsm/internal/blockchain-key/security"
	blockchainService "soft-hsm/internal/blockchain-key/services"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func SetupRouter(r *chi.Mux, cfg *config.Config, db *storage.Postgres, redis *storage.Redis) *chi.Mux {
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},                     // Доступные домены
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}, // Доступные методы
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},   // Доступные заголовки
		ExposedHeaders:   []string{"Link"},                                      // Заголовки, видимые клиенту
		AllowCredentials: true,                                                  // Разрешить отправку куки
		MaxAge:           300,
	}))
	setupAuthRoutes(r, cfg, db, redis)
	setupUserRoutes(r, db)
	setupBlockchainKeyRoutes(r, db)

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

func setupAuthRoutes(r *chi.Mux, cfg *config.Config, db *storage.Postgres, redis *storage.Redis) {
	tokenRepo := authRepository.NewTokenRepository(redis)
	userRepo := userRepository.NewUserRepository(db)

	claimsService := services.NewClaimsService(cfg)
	mailerService := mailer.NewMailer(&cfg.MailerConfig)
	passwordService := services.NewPasswordService()
	activationService := services.NewActivationService(userRepo, tokenRepo, claimsService)
	authService := services.NewAuthService(tokenRepo, claimsService, userRepo, mailerService, passwordService)

	authHandler := handlers.NewAuthHandler(authService, activationService)

	r.Mount("/v1/auth", auth.AuthRoutes(authHandler))
}

func setupBlockchainKeyRoutes(r *chi.Mux, db *storage.Postgres) {
	blockchainKeyRepo := blockchainRepository.NewBlockchainKeyRepository(db)
	passwordService := services.NewPasswordService()

	securityService := security.NewSecurityService()
	blockchainKeyService := blockchainService.NewBlockchainKeyService(blockchainKeyRepo, securityService, passwordService)

	blockchainKeyHandler := blockchainHandler.NewBlockchainKeyHandler(blockchainKeyService)

	// r.Mount("/v1/blockchain", blockchainKey.BlockchainKeyRoutes(blockchainKeyHandler))
	r.Route("/v1/blockchain", func(r chi.Router) {
		r.Use(middleware.AuthMiddleware) // Применяем middleware ко всем маршрутам
		r.Mount("/", blockchainKey.BlockchainKeyRoutes(blockchainKeyHandler))
	})
}

func setupUserRoutes(r *chi.Mux, db *storage.Postgres) {
	userRepo := userRepository.NewUserRepository(db)
	userService := userService.NewUserService(userRepo)
	userHandlers := userHandler.NewUserHandler(userService)

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Get("/v1/user/me", userHandlers.Me)
	})
}
