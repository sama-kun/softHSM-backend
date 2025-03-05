package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"soft-hsm/internal/config"
	"soft-hsm/internal/lib/logger/sl"
	mw "soft-hsm/internal/middleware"
	"soft-hsm/internal/storage"
	"syscall"
	"time"

	"log/slog"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	cfg := config.MustLoad()

	slog.Info("Project running in: ", slog.String("ADDRESS", cfg.Address))

	_, err := storage.NewPostgresDB(cfg.Database)

	if err != nil {
		slog.Error("Failed to init DB", sl.Err(err))
		os.Exit(1)
	}
	redisClient, err := storage.NewRedis(cfg.RedisConfig)

	if err != nil {
		slog.Error("Ошибка подключения к Redis", "error", err)
		os.Exit(1)
	}
	defer redisClient.Close()
	slog.Info("Подключение к Redis установлено")

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	router.Use(mw.JSONResponseMiddleware)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Server is running!"))
	})

	// Запуск сервера
	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	// Горутина для запуска сервера
	go func() {
		slog.Info("Starting server...", slog.String("address", cfg.HTTPServer.Address))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("Server failed to start", sl.Err(err))
			os.Exit(1)
		}
	}()

	// Ожидание сигнала завершения (graceful shutdown)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop // Ждём сигнал

	slog.Info("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Server forced to shutdown", sl.Err(err))
	}

	slog.Info("Server exited successfully")

}
