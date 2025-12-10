package handler

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"log/slog"
	"net/http"
	"wallet_controller/config"
	"wallet_controller/internal/entity"
	"wallet_controller/internal/service"
)

type WalletHandler struct {
	walletService *service.WalletService
	cfg           *config.Config
}

func NewWalletHandler(ctx context.Context, cfg *config.Config) *WalletHandler {
	return &WalletHandler{
		walletService: service.NewWalletService(cfg.Client),
		cfg:           cfg,
	}
}

func (h *WalletHandler) GetWallet(c *gin.Context) {
	walletID, err := uuid.Parse(c.Query("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no id_wallet provided"})
		return
	}

	wallet, err := h.walletService.GetWallet(c.Request.Context(), walletID)
	if err != nil {
		slog.Error("Get wallet err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, wallet)
}

func (h *WalletHandler) AddOperation(c *gin.Context) {
	var req entity.OperationRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	wallet, err := h.walletService.AddOperation(c.Request.Context(), &req)
	if err != nil {
		slog.Error("Add operation err", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"wallet": wallet})
}
