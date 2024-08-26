package db

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"cmd/app/main.go/internal/config"
)

func OpenDB(ctx context.Context, cfg config.DBConfig) (*pgxpool.Pool, error) {
	parseConfig, err := pgxpool.ParseConfig("")
	if err != nil {
		return nil, fmt.Errorf("OpenDB config parse: %w", err)
	}

	parseConfig.ConnConfig.Host = cfg.PgHost
	parseConfig.ConnConfig.Port = cfg.PgPort
	parseConfig.ConnConfig.Database = cfg.PgDatabase
	parseConfig.ConnConfig.User = cfg.PgUser
	parseConfig.ConnConfig.Password = cfg.PgPassword

	pool, err := pgxpool.NewWithConfig(ctx, parseConfig)
	if err != nil {
		return nil, fmt.Errorf("OpenDB connect: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("OpenDB ping: %w", err)
	}

	return pool, nil
}
