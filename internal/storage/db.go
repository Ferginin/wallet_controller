package storage

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

//go:embed 01-init.sql
var initSQL string

//go:embed 02-data.sql
var dataSQL string

func Migrate(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), initSQL)
	if err != nil {
		slog.Error("Init DB Error: ", err.Error(), nil)
		return fmt.Errorf("init sql failed: %s", err.Error())
	}

	return nil
}

func DataInsert(db *pgxpool.Pool) error {
	_, err := db.Exec(context.Background(), dataSQL)
	if err != nil {
		slog.Error("Data Insert Error: ", err.Error(), nil)
		return fmt.Errorf("data insert failed: %s", err.Error())
	}
	return nil
}
