package service

import (
	"context"
	"github.com/google/uuid"
	"log/slog"
	"wallet_controller/internal/entity"
	"wallet_controller/internal/repository"
)

type WalletServiceInterface interface {
	GetWallet(ctx context.Context, walletID uuid.UUID) (*entity.Wallet, error)
	AddOperation(ctx context.Context, operation *entity.OperationRequest) (entity.Wallet, error)
}

type WalletService struct {
	walletRepo repository.WalletRepositoryInterface
}

func NewWalletService(walletRepo repository.WalletRepositoryInterface) WalletServiceInterface {
	return &WalletService{
		walletRepo: walletRepo,
	}
}

func (s *WalletService) GetWallet(ctx context.Context, walletID uuid.UUID) (*entity.Wallet, error) {
	wallet, err := s.walletRepo.GetByID(ctx, walletID)
	if err != nil {
		return nil, err
	}

	return wallet, nil
}

func (s *WalletService) AddOperation(ctx context.Context, operation *entity.OperationRequest) (entity.Wallet, error) {
	Wallet, err := s.walletRepo.AddOperation(ctx, operation.WalletID, operation.OperationType, operation.Amount*100)
	if err != nil {
		slog.Error("WalletService", "AddOperation", "err", err.Error(), nil)
		return entity.Wallet{}, err
	}

	return Wallet, nil
}
