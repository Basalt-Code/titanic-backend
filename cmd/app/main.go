package main

import (
	"cmd/app/main.go/internal/api"
	"cmd/app/main.go/internal/config"
	store "cmd/app/main.go/internal/store/sqlstore"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq" //.
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {
	// считываем переменные и создаем с их значениями конфиг
	cfg := config.MustLoad()

	// создаем объект логгера
	log := setupLogger(cfg.Env)

	storage, err := store.New(cfg.StoragePath)
	if err != nil {
		log.Error("failed to init storage")
		os.Exit(1)
	}

	router := chi.NewRouter()
	// router.Use(middleware.RequestID)

	// Хендлеры
	router.Post("/users", func(w http.ResponseWriter, r *http.Request) {
		api.CreateUserHandler(w, r, storage)
	})

	srv := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	log.Info("Server is starting on port 8000...")
	if err := srv.ListenAndServe(); err != nil {
		log.Error("failed to start server")
	}

}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}
	return log
}
