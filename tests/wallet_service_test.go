package tests

import (
	"context"
	"testing"
	"wallet_controller/internal/service"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"wallet_controller/internal/entity"
)

type MockWalletRepository struct {
	mock.Mock
}

func (m *MockWalletRepository) GetByID(ctx context.Context, walletID uuid.UUID) (*entity.Wallet, error) {
	args := m.Called(ctx, walletID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Wallet), args.Error(1)
}

func (m *MockWalletRepository) AddOperation(ctx context.Context, walletID uuid.UUID, operationType string, amount int) (entity.Wallet, error) {
	args := m.Called(ctx, walletID, operationType, amount)
	return args.Get(0).(entity.Wallet), args.Error(1)
}

func TestGetWallet_Success(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	walletID := uuid.New()
	expectedWallet := &entity.Wallet{
		ID:      walletID,
		Balance: 5000,
	}

	mockRepo.On("GetByID", mock.Anything, walletID).Return(expectedWallet, nil)

	mService := service.NewWalletService(mockRepo)

	ctx := context.Background()
	wallet, err := mService.GetWallet(ctx, walletID)

	assert.NoError(t, err)
	assert.Equal(t, expectedWallet, wallet)
	mockRepo.AssertExpectations(t)
}

func TestGetWallet_NotFound(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	walletID := uuid.New()

	mockRepo.On("GetByID", mock.Anything, walletID).Return(nil, assert.AnError)

	mService := service.NewWalletService(mockRepo)

	ctx := context.Background()
	wallet, err := mService.GetWallet(ctx, walletID)

	assert.Error(t, err)
	assert.Nil(t, wallet)
	mockRepo.AssertExpectations(t)
}

func TestAddOperation_Deposit_Success(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	walletID := uuid.New()

	req := &entity.OperationRequest{
		WalletID:      walletID,
		OperationType: "DEPOSIT",
		Amount:        100, // 100 рублей = 10000 копеек
	}

	expectedWallet := entity.Wallet{
		ID:      walletID,
		Balance: 10000, // было 0, добавили 100 * 100
	}

	mockRepo.On("AddOperation", mock.Anything, walletID, "DEPOSIT", 10000).
		Return(expectedWallet, nil)

	mService := service.NewWalletService(mockRepo)

	ctx := context.Background()
	wallet, err := mService.AddOperation(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedWallet, wallet)
	mockRepo.AssertExpectations(t)
}

func TestAddOperation_Withdraw_Success(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	walletID := uuid.New()

	req := &entity.OperationRequest{
		WalletID:      walletID,
		OperationType: "WITHDRAW",
		Amount:        50,
	}

	expectedWallet := entity.Wallet{
		ID:      walletID,
		Balance: 5000,
	}

	mockRepo.On("AddOperation", mock.Anything, walletID, "WITHDRAW", 5000).
		Return(expectedWallet, nil)

	mService := service.NewWalletService(mockRepo)

	ctx := context.Background()
	wallet, err := mService.AddOperation(ctx, req)

	assert.NoError(t, err)
	assert.Equal(t, expectedWallet, wallet)
	mockRepo.AssertExpectations(t)
}

func TestAddOperation_InsufficientFunds(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	walletID := uuid.New()

	req := &entity.OperationRequest{
		WalletID:      walletID,
		OperationType: "WITHDRAW",
		Amount:        500,
	}

	mockRepo.On("AddOperation", mock.Anything, walletID, "WITHDRAW", 50000).
		Return(entity.Wallet{}, assert.AnError)

	mService := service.NewWalletService(mockRepo)

	ctx := context.Background()
	wallet, err := mService.AddOperation(ctx, req)

	assert.Error(t, err)
	assert.Equal(t, entity.Wallet{}, wallet)
	mockRepo.AssertExpectations(t)
}

func TestAddOperation_MultipleOperations(t *testing.T) {
	mockRepo := new(MockWalletRepository)
	walletID := uuid.New()

	req1 := &entity.OperationRequest{
		WalletID:      walletID,
		OperationType: "DEPOSIT",
		Amount:        100,
	}

	wallet1 := entity.Wallet{
		ID:      walletID,
		Balance: 10000,
	}

	mockRepo.On("AddOperation", mock.Anything, walletID, "DEPOSIT", 10000).
		Return(wallet1, nil)

	mService := service.NewWalletService(mockRepo)

	ctx := context.Background()
	w1, err := mService.AddOperation(ctx, req1)
	assert.NoError(t, err)
	assert.Equal(t, wallet1, w1)

	req2 := &entity.OperationRequest{
		WalletID:      walletID,
		OperationType: "WITHDRAW",
		Amount:        50,
	}

	wallet2 := entity.Wallet{
		ID:      walletID,
		Balance: 5000,
	}

	mockRepo.On("AddOperation", mock.Anything, walletID, "WITHDRAW", 5000).
		Return(wallet2, nil)

	w2, err := mService.AddOperation(ctx, req2)
	assert.NoError(t, err)
	assert.Equal(t, wallet2, w2)

	mockRepo.AssertNumberOfCalls(t, "AddOperation", 2)
}
