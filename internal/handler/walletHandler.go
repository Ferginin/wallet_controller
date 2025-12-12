package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"wallet_controller/internal/entity"
	"wallet_controller/internal/service"
)

type WalletHandler struct {
	walletService service.WalletServiceInterface
}

func NewWalletHandler(walletService service.WalletServiceInterface) *WalletHandler {
	return &WalletHandler{
		walletService: walletService,
	}
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	walletIDStr := c.Param("id")
	if walletIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wallet_id is required in path"})
		return
	}

	walletID, err := uuid.Parse(walletIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid wallet_id format"})
		return
	}

	wallet, err := h.walletService.GetWallet(c.Request.Context(), walletID)
	if err != nil {
		slog.Error("Get wallet error", "error", err.Error(), "wallet_id", walletID)

		if err.Error() == "wallet not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": "wallet not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get wallet"})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (h *WalletHandler) AddOperation(c *gin.Context) {
	var req entity.OperationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		slog.Warn("Invalid request body", "error", err.Error())
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if req.WalletID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "wallet_id is required"})
		return
	}
	if req.OperationType != "DEPOSIT" && req.OperationType != "WITHDRAW" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "operation_type must be 'DEPOSIT' or 'WITHDRAW'",
		})
		return
	}
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "amount must be positive"})
		return
	}

	wallet, err := h.walletService.AddOperation(c.Request.Context(), &req)
	if err != nil {
		slog.Error("Add operation err", err.Error(), nil)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": wallet})
}
