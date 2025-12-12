package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"wallet_controller/config"
	"wallet_controller/internal/handler"
	"wallet_controller/internal/repository"
	"wallet_controller/internal/service"
)

func SetupRouter(ctx context.Context, cfg *config.Config) *gin.Engine {

	walletRepo := repository.NewWalletRepository(cfg.Client)
	walletService := service.NewWalletService(walletRepo)
	walletHandler := handler.NewWalletHandler(walletService)

	if cfg.Env.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	api.GET("/wallets/:id", walletHandler.GetWallet)
	api.POST("/wallet", walletHandler.AddOperation)

	return r
}
