package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"wallet_controller/internal/handler"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"wallet_controller/internal/entity"
)

type MockWalletService struct {
	mock.Mock
}

func (m *MockWalletService) GetWallet(ctx context.Context, walletID uuid.UUID) (*entity.Wallet, error) {
	args := m.Called(ctx, walletID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*entity.Wallet), args.Error(1)
}

func (m *MockWalletService) AddOperation(ctx context.Context, operation *entity.OperationRequest) (entity.Wallet, error) {
	args := m.Called(ctx, operation)
	return args.Get(0).(entity.Wallet), args.Error(1)
}

func setupGinRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

func TestHandlerGetWallet_Success(t *testing.T) {
	mockService := new(MockWalletService)
	walletID := uuid.New()
	expectedWallet := &entity.Wallet{
		ID:      walletID,
		Balance: 5000,
	}

	mockService.On("GetWallet", mock.Anything, walletID).Return(expectedWallet, nil)

	mHandler := handler.NewWalletHandler(mockService)

	router := setupGinRouter()
	router.GET("/wallets/:id", mHandler.GetWallet)

	req := httptest.NewRequest(http.MethodGet, "/wallets/"+walletID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var wallet entity.Wallet
	err := json.Unmarshal(w.Body.Bytes(), &wallet)
	assert.NoError(t, err)
	assert.Equal(t, walletID, wallet.ID)
	assert.Equal(t, 5000, wallet.Balance)

	mockService.AssertExpectations(t)
}

func TestHandlerGetWallet_NotFound(t *testing.T) {
	mockService := new(MockWalletService)
	walletID := uuid.New()

	mockService.On("GetWallet", mock.Anything, walletID).
		Return(nil, assert.AnError)

	mHandler := handler.NewWalletHandler(mockService)

	router := setupGinRouter()
	router.GET("/wallets/:id", mHandler.GetWallet)

	req := httptest.NewRequest(http.MethodGet, "/wallets/"+walletID.String(), nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	var errResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.Contains(t, errResp["error"], "failed")

	mockService.AssertExpectations(t)
}

func TestHandlerGetWallet_InvalidUUID(t *testing.T) {
	mockService := new(MockWalletService)

	mHandler := handler.NewWalletHandler(mockService)

	router := setupGinRouter()
	router.GET("/wallets/:id", mHandler.GetWallet)

	req := httptest.NewRequest(http.MethodGet, "/wallets/invalid-uuid", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.Contains(t, errResp["error"], "invalid")
}

func TestHandlerGetWallet_MissingID(t *testing.T) {
	mockService := new(MockWalletService)

	mHandler := handler.NewWalletHandler(mockService)

	router := setupGinRouter()
	router.GET("/wallets/:id", mHandler.GetWallet)

	req := httptest.NewRequest(http.MethodGet, "/wallets/", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandlerAddOperation_Success_Deposit(t *testing.T) {
	mockService := new(MockWalletService)
	walletID := uuid.New()

	req := &entity.OperationRequest{
		WalletID:      walletID,
		OperationType: "DEPOSIT",
		Amount:        100,
	}

	expectedWallet := entity.Wallet{
		ID:      walletID,
		Balance: 10000,
	}

	mockService.On("AddOperation", mock.Anything, mock.MatchedBy(func(r *entity.OperationRequest) bool {
		return r.WalletID == walletID && r.OperationType == "DEPOSIT" && r.Amount == 100
	})).Return(expectedWallet, nil)

	mHandler := handler.NewWalletHandler(mockService)

	router := setupGinRouter()
	router.POST("/wallet", mHandler.AddOperation)

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var respBody map[string]entity.Wallet
	json.Unmarshal(w.Body.Bytes(), &respBody)
	wallet := respBody["wallet"]
	assert.Equal(t, expectedWallet, wallet)

	mockService.AssertExpectations(t)
}

func TestHandlerAddOperation_Success_Withdraw(t *testing.T) {
	mockService := new(MockWalletService)
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

	mockService.On("AddOperation", mock.Anything, mock.MatchedBy(func(r *entity.OperationRequest) bool {
		return r.WalletID == walletID && r.OperationType == "WITHDRAW" && r.Amount == 50
	})).Return(expectedWallet, nil)

	mHandler := handler.NewWalletHandler(mockService)

	router := setupGinRouter()
	router.POST("/wallet", mHandler.AddOperation)

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusOK, w.Code)

	var respBody map[string]entity.Wallet
	json.Unmarshal(w.Body.Bytes(), &respBody)
	wallet := respBody["wallet"]
	assert.Equal(t, expectedWallet, wallet)

	mockService.AssertExpectations(t)
}

func TestHandlerAddOperation_InvalidOperationType(t *testing.T) {
	mockService := new(MockWalletService)
	walletID := uuid.New()

	req := map[string]interface{}{
		"wallet_id":     walletID.String(),
		"operationType": "TRANSFER",
		"amount":        100,
	}

	mHandler := handler.NewWalletHandler(mockService)

	router := setupGinRouter()
	router.POST("/wallet", mHandler.AddOperation)

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	var errResp map[string]string
	json.Unmarshal(w.Body.Bytes(), &errResp)
	assert.Contains(t, errResp["error"], "DEPOSIT")
}

func TestHandlerAddOperation_NegativeAmount(t *testing.T) {
	mockService := new(MockWalletService)
	walletID := uuid.New()

	req := map[string]interface{}{
		"walletId":      walletID.String(),
		"operationType": "DEPOSIT",
		"amount":        -100,
	}

	mHandler := handler.NewWalletHandler(mockService)

	router := setupGinRouter()
	router.POST("/wallet", mHandler.AddOperation)

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerAddOperation_ZeroAmount(t *testing.T) {
	mockService := new(MockWalletService)
	walletID := uuid.New()

	req := map[string]interface{}{
		"walletId":      walletID.String(),
		"operationType": "DEPOSIT",
		"amount":        0,
	}

	mHandler := handler.NewWalletHandler(mockService)

	router := setupGinRouter()
	router.POST("/wallet", mHandler.AddOperation)

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerAddOperation_MissingWalletID(t *testing.T) {
	mockService := new(MockWalletService)

	req := map[string]interface{}{
		"operationType": "DEPOSIT",
		"amount":        100,
	}

	mHandler := handler.NewWalletHandler(mockService)

	router := setupGinRouter()
	router.POST("/wallet", mHandler.AddOperation)

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerAddOperation_InvalidJSON(t *testing.T) {
	mockService := new(MockWalletService)

	mHandler := handler.NewWalletHandler(mockService)

	router := setupGinRouter()
	router.POST("/wallet", mHandler.AddOperation)

	httpReq := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewReader([]byte("invalid json")))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestHandlerAddOperation_ServiceError(t *testing.T) {
	mockService := new(MockWalletService)
	walletID := uuid.New()

	req := &entity.OperationRequest{
		WalletID:      walletID,
		OperationType: "WITHDRAW",
		Amount:        1000,
	}

	mockService.On("AddOperation", mock.Anything, mock.Anything).
		Return(entity.Wallet{}, assert.AnError)

	mHandler := handler.NewWalletHandler(mockService)

	router := setupGinRouter()
	router.POST("/wallet", mHandler.AddOperation)

	body, _ := json.Marshal(req)
	httpReq := httptest.NewRequest(http.MethodPost, "/wallet", bytes.NewReader(body))
	httpReq.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, httpReq)

	assert.Equal(t, http.StatusInternalServerError, w.Code)

	mockService.AssertExpectations(t)
}
