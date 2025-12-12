package storage

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"log/slog"
	"time"
	"wallet_controller/config"
)

func NewConnection(ctx context.Context, cfg *config.Config) *pgxpool.Pool {
	env := cfg.Env

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%d/%s",
		env.DbUsername,
		env.DbPassword,
		env.DbHost,
		env.DbPort,
		env.DbName,
	)

	conConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		slog.Error("Unable to parse connection string")
		log.Fatal("Unable to parse config:", err.Error())
	}

	conConfig.MaxConns = 100
	conConfig.MinConns = 5
	conConfig.MaxConnLifetime = 30 * time.Minute
	conConfig.MaxConnIdleTime = 5 * time.Minute

	conn, err := pgxpool.NewWithConfig(context.Background(), conConfig)
	if err != nil {
		slog.Error("Unable to connect to database:", err.Error(), nil)
		panic(err)
	}

	if err = Migrate(conn); err != nil {
		slog.Error("Unable to migrate database:", err.Error(), nil)
		panic(err)
	}
	if err = DataInsert(conn); err != nil {
		slog.Error("Unable to migrate data:", err.Error(), nil)
		panic(err)
	}

	slog.Info("Connected to database")

	return conn
}
