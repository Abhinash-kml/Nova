package economy

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, cc *CurrencyController, wc *WalletController) {
	currency := router.Group("/private/currency")
	{
		currency.GET("", cc.GetAll)
		currency.GET("/:id", cc.Get)
		currency.POST("", cc.Create)
		currency.PATCH("/:id", cc.Update)
		currency.DELETE("/:id", cc.Delete)
	}

	walletPublic := router.Group("/public/wallets")
	{
		walletPublic.GET("", wc.GetWalletOfPlayer)
	}

	walletPrivate := router.Group("/private/wallets")
	{
		walletPrivate.POST("", wc.CreateWalletOfNewPlayer)
		walletPrivate.PATCH("", wc.UpdateWalletOfPlayer)
		walletPrivate.DELETE("", wc.DeleteWalletOfPlayer)
	}
}
