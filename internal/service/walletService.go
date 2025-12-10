package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"wallet_controller/internal/entity"
	"wallet_controller/internal/repository"
)

type WalletService struct {
	walletRepo *repository.WalletRepository
}

func NewWalletService(db *pgxpool.Pool) *WalletService {
	return &WalletService{
		walletRepo: repository.NewWalletRepository(db),
	}
}

func (s *WalletService) GetWallet(ctx context.Context, walletID uuid.UUID) (*entity.Wallet, error) {
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *WalletService) AddOperation(ctx context.Context, operation *entity.OperationRequest) (*entity.Wallet, error) {
	Wallet, err := s.walletRepo.AddOperation(ctx, operation.WalletID, operation.OperationType, operation.Amount*100)
	if err != nil {
		slog.Error("WalletService", "AddOperation", "err", err.Error())
		return nil, err
	}

	return Wallet, nil
}
