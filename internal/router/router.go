package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"wallet_controller/config"
	"wallet_controller/internal/handler"
)

func SetupRouter(ctx context.Context, cfg *config.Config) *gin.Engine {
	if cfg.Env.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	walletHandler := handler.NewWalletHandler(ctx, cfg)

	api := r.Group("/")
	{
		wallet := api.Group("/api/v1")
		{
			wallet.POST("/wallet", walletHandler.AddOperation)
			wallet.GET("/wallets", walletHandler.GetWallet)
		}
	}

	return r
}
