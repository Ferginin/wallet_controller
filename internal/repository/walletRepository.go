package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"wallet_controller/internal/entity"
)

type WalletRepository struct {
	db *pgxpool.Pool
}

func NewWalletRepository(db *pgxpool.Pool) *WalletRepository {
	return &WalletRepository{db: db}
}

func (r *WalletRepository) GetByID(ctx context.Context, walletID uuid.UUID) (*entity.Wallet, error) {
	query := `
		SELECT id_wallet, balance
		FROM wallets
		WHERE id_wallet = $1
	`

	wallet := entity.Wallet{}
	err := r.db.QueryRow(ctx, query, walletID).Scan(
		&wallet.ID,
		&wallet.Balance,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("wallet not found")
		}
		return nil, fmt.Errorf("failed to get wallet: %w", err)
	}

	return &wallet, err
}

func (r *WalletRepository) AddOperation(ctx context.Context, walletID uuid.UUID, operationType string, amount int) (*entity.Wallet, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	balance := 0
	err = tx.QueryRow(ctx,
		`SELECT balance FROM wallets WHERE id_wallet = $1`,
		walletID,
	).Scan(&balance)
	if err != nil {
		slog.Error("failed to get wallet balance", err)
		return nil, err
	}

	if operationType == "WITHDRAW" && balance-amount < 0 {
		slog.Warn("Not enough money on wallet:", walletID)
		return nil, errors.New("not enough money on wallet")
	} else if operationType == "WITHDRAW" {
		balance -= amount
	} else {
		balance += amount
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO wallet_operations (id_wallet, operation_type, amount)
		VALUES ($1, $2, $3)`,
		walletID,
		operationType,
		amount,
	)
	if err != nil {
		slog.Error("failed to insert wallet operation:", err)
		return nil, err
	}

	_, err = tx.Exec(ctx,
		`UPDATE wallets
			SET balance = $1, updated_at = CURRENT_TIMESTAMP
			WHERE id_wallet = $2`,
		balance,
		walletID,
	)
	if err != nil {
		slog.Error("failed to update wallet operation:", err)
		return nil, err
	}

	return &(entity.Wallet{
			ID:      walletID,
			Balance: balance,
		}),
		tx.Commit(ctx)
}
