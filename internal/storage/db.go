package storage

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed init.sql
var initSQL string

//go:embed data.sql
var dataSQL string

func Migrate(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), initSQL)
	if err != nil {
		slog.Error("Init DB Error: ", err)
		return fmt.Errorf("init sql failed: %w", err)
	}

	return nil
}

func DataInsert(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), dataSQL)
	if err != nil {
		slog.Error("Data Insert Error: ", err)
		return fmt.Errorf("data insert failed: %w", err)
	}
	return nil
}
