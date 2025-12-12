package tests

import (
	"context"
	"testing"
	"wallet_controller/internal/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"wallet_controller/internal/entity"
)

func setupTestDB(t *testing.T) *pgxpool.Pool {
	// Для локальной разработки используйте тестовую БД
	connString := "postgresql://postgres:postgres@localhost:5432/wallet_test"

	config, err := pgxpool.ParseConfig(connString)
	require.NoError(t, err, "Failed to parse connection config")

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	require.NoError(t, err, "Failed to create connection pool")

	ctx := context.Background()
	_, err = pool.Exec(ctx, `
		CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
		
		DROP TABLE IF EXISTS wallet_operations;
		DROP TABLE IF EXISTS wallets;
		
		CREATE TABLE wallets (
			id_wallet UUID PRIMARY KEY,
			balance BIGINT NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE TABLE wallet_operations (
			id_operation UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
			id_wallet UUID NOT NULL REFERENCES wallets(id_wallet) ON DELETE CASCADE,
			operation_type VARCHAR(16) NOT NULL CHECK (operation_type IN ('DEPOSIT', 'WITHDRAW')),
			amount BIGINT NOT NULL CHECK (amount > 0),
			created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
		
		CREATE INDEX IF NOT EXISTS idx_wallet_operations_wallet_id_created_at
		ON wallet_operations (id_wallet, created_at);
	`)
	require.NoError(t, err, "Failed to create schema")

	return pool
}

func teardownTestDB(t *testing.T, pool *pgxpool.Pool) {
	ctx := context.Background()
	_, err := pool.Exec(ctx, `
		DROP TABLE IF EXISTS wallet_operations;
		DROP TABLE IF EXISTS wallets;
	`)
	require.NoError(t, err, "Failed to cleanup database")
	pool.Close()
}

func TestGetByID_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	walletID := uuid.New()

	_, err := pool.Exec(ctx, `
		INSERT INTO wallets (id_wallet, balance)
		VALUES ($1, $2)
	`, walletID, 10000)
	require.NoError(t, err)

	repo := repository.NewWalletRepository(pool)
	wallet, err := repo.GetByID(ctx, walletID)

	assert.NoError(t, err)
	assert.NotNil(t, wallet)
	assert.Equal(t, walletID, wallet.ID)
	assert.Equal(t, 10000, wallet.Balance)
}

func TestGetByID_NotFound(t *testing.T) {
	pool := setupTestDB(t)
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	nonExistentID := uuid.New()

	repo := repository.NewWalletRepository(pool)
	wallet, err := repo.GetByID(ctx, nonExistentID)

	assert.Error(t, err)
	assert.Nil(t, wallet)
	assert.Equal(t, "wallet not found", err.Error())
}

func TestRepoAddOperation_Deposit_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	walletID := uuid.New()
	initialBalance := int(5000)
	depositAmount := int(2000)

	_, err := pool.Exec(ctx, `
		INSERT INTO wallets (id_wallet, balance)
		VALUES ($1, $2)
	`, walletID, initialBalance)
	require.NoError(t, err)

	repo := repository.NewWalletRepository(pool)
	wallet, err := repo.AddOperation(ctx, walletID, "DEPOSIT", depositAmount)

	assert.NoError(t, err)
	assert.Equal(t, walletID, wallet.ID)
	assert.Equal(t, initialBalance+depositAmount, wallet.Balance)

	var opCount int
	err = pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM wallet_operations
		WHERE id_wallet = $1 AND operation_type = 'DEPOSIT'
	`, walletID).Scan(&opCount)
	require.NoError(t, err)
	assert.Equal(t, 1, opCount)
}

func TestRepoAddOperation_Withdraw_Success(t *testing.T) {
	pool := setupTestDB(t)
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	walletID := uuid.New()
	initialBalance := int(5000)
	withdrawAmount := int(2000)

	_, err := pool.Exec(ctx, `
		INSERT INTO wallets (id_wallet, balance)
		VALUES ($1, $2)
	`, walletID, initialBalance)
	require.NoError(t, err)

	repo := repository.NewWalletRepository(pool)
	wallet, err := repo.AddOperation(ctx, walletID, "WITHDRAW", withdrawAmount)

	assert.NoError(t, err)
	assert.Equal(t, walletID, wallet.ID)
	assert.Equal(t, initialBalance-withdrawAmount, wallet.Balance)
}

func TestAddOperation_Withdraw_InsufficientFunds(t *testing.T) {
	pool := setupTestDB(t)
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	walletID := uuid.New()
	initialBalance := int(1000)
	withdrawAmount := int(2000)

	_, err := pool.Exec(ctx, `
		INSERT INTO wallets (id_wallet, balance)
		VALUES ($1, $2)
	`, walletID, initialBalance)
	require.NoError(t, err)

	repo := repository.NewWalletRepository(pool)
	wallet, err := repo.AddOperation(ctx, walletID, "WITHDRAW", withdrawAmount)

	assert.Error(t, err)
	assert.Equal(t, "not enough money on wallet", err.Error())
	assert.Equal(t, int64(0), int64(wallet.Balance))
}

func TestAddOperation_Withdraw_ExactAmount(t *testing.T) {
	pool := setupTestDB(t)
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	walletID := uuid.New()
	initialBalance := int(2000)
	withdrawAmount := int(2000)

	_, err := pool.Exec(ctx, `
		INSERT INTO wallets (id_wallet, balance)
		VALUES ($1, $2)
	`, walletID, initialBalance)
	require.NoError(t, err)

	repo := repository.NewWalletRepository(pool)
	wallet, err := repo.AddOperation(ctx, walletID, "WITHDRAW", withdrawAmount)

	assert.NoError(t, err)
	assert.Equal(t, 0, wallet.Balance)
}

func TestAddOperation_MultipleTransactions(t *testing.T) {
	pool := setupTestDB(t)
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	walletID := uuid.New()
	initialBalance := int(5000)

	_, err := pool.Exec(ctx, `
		INSERT INTO wallets (id_wallet, balance)
		VALUES ($1, $2)
	`, walletID, initialBalance)
	require.NoError(t, err)

	repo := repository.NewWalletRepository(pool)

	wallet, err := repo.AddOperation(ctx, walletID, "DEPOSIT", 1000)
	assert.NoError(t, err)
	assert.Equal(t, 6000, wallet.Balance)

	wallet, err = repo.AddOperation(ctx, walletID, "WITHDRAW", 1500)
	assert.NoError(t, err)
	assert.Equal(t, 4500, wallet.Balance)

	wallet, err = repo.AddOperation(ctx, walletID, "DEPOSIT", 2500)
	assert.NoError(t, err)
	assert.Equal(t, 7000, wallet.Balance)

	var opCount int
	err = pool.QueryRow(ctx, `
		SELECT COUNT(*) FROM wallet_operations
		WHERE id_wallet = $1
	`, walletID).Scan(&opCount)
	require.NoError(t, err)
	assert.Equal(t, 3, opCount)
}

func TestAddOperation_WalletNotFound(t *testing.T) {
	pool := setupTestDB(t)
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	nonExistentID := uuid.New()

	repo := repository.NewWalletRepository(pool)
	wallet, err := repo.AddOperation(ctx, nonExistentID, "DEPOSIT", 1000)

	assert.Error(t, err)
	assert.Equal(t, entity.Wallet{}, wallet)
}

func TestAddOperation_LargeAmounts(t *testing.T) {
	pool := setupTestDB(t)
	defer teardownTestDB(t, pool)

	ctx := context.Background()
	walletID := uuid.New()
	initialBalance := int(1000000000)
	depositAmount := int(9999999999)

	_, err := pool.Exec(ctx, `
		INSERT INTO wallets (id_wallet, balance)
		VALUES ($1, $2)
	`, walletID, initialBalance)
	require.NoError(t, err)

	repo := repository.NewWalletRepository(pool)
	wallet, err := repo.AddOperation(ctx, walletID, "DEPOSIT", depositAmount)

	assert.NoError(t, err)
	assert.Equal(t, initialBalance+depositAmount, wallet.Balance)
}
