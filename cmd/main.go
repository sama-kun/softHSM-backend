package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"soft-hsm/internal/config"
	"soft-hsm/internal/lib/logger/sl"
	"soft-hsm/internal/storage"
	"syscall"
	"time"

	"log/slog"
	routerGroup "soft-hsm/internal/router"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()

	slog.Info("Project running in: ", slog.String("ADDRESS", cfg.Address))

	db, err := storage.NewPostgresDB(cfg.Database)
	if err != nil {
		slog.Error("Failed to init DB", sl.Err(err))
		os.Exit(1)
	}
	defer db.Close()

	redisClient, err := storage.NewRedis(cfg.RedisConfig)
	if err != nil {
		slog.Error("Ошибка подключения к Redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	slog.Info("Подключение к Redis установлено")

	r := chi.NewRouter()

	// r.Use(cors.Handler(cors.Options{
	// 	AllowedOrigins:   []string{"*"},
	// 	AllowedMethods:   []string{"*"},
	// 	AllowedHeaders:   []string{"*"},
	// 	AllowCredentials: false,
	// 	MaxAge:           300,
	// }))

	// Подключаем middleware ДО создания маршрутов
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)
	// r.Use(mw.JSONResponseMiddleware)
	

	routerGroup.SetupRouter(r,cfg, db, redisClient)

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      r,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	go func() {
		slog.Info("Starting server...", slog.String("address", cfg.HTTPServer.Address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", sl.Err(err))
			os.Exit(1)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", sl.Err(err))
	}

	slog.Info("Server exited successfully")
}

// func enableCORS(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

// 		if r.Method == http.MethodOptions {
// 			w.WriteHeader(http.StatusNoContent)
// 			return
// 		}

// 		next.ServeHTTP(w, r)
// 	})
// }