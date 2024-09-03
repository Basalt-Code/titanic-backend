package main

import (
	"context"
	"errors"
	"log"
	"net/http"

	"cmd/app/main.go/internal/api"
	"cmd/app/main.go/internal/config"
	"cmd/app/main.go/internal/pkg/db"
	"cmd/app/main.go/internal/pkg/logger"
	authrepo "cmd/app/main.go/internal/repository/auth"
	userrepo "cmd/app/main.go/internal/repository/user"
	authservices "cmd/app/main.go/internal/services/auth"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}

	pool, err := db.OpenDB(context.Background(), cfg.DBConfig)
	if err != nil {
		panic(err)
	}

	logs, err := logger.New(cfg.ServerConfig.LogFilePath)
	if err != nil {
		panic(err)
	}

	r := api.New(
		logs,
		authservices.New(
			cfg.ServerConfig,
			userrepo.New(pool),
			authrepo.New(pool),
		),
	)

	srv := &http.Server{
		Addr:    ":" + cfg.ServerConfig.ServerPort,
		Handler: r,
	}

	logs.Info("Successfully started server on port " + cfg.ServerConfig.ServerPort)
	if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("listen: %s\n", err)
	}
}
