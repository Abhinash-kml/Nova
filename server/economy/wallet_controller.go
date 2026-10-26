package economy

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type WalletController struct {
	service WalletService
	config  *config.WalletConfig
	logger  *zap.Logger
}

func NewWalletController(service WalletService, c *config.WalletConfig, l *zap.Logger) *WalletController {
	return &WalletController{
		service: service,
		config:  c,
		logger:  l,
	}
}

func (c *WalletController) CreateWalletOfNewPlayer(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.wallet.CreateWalletOfNewPlayer")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to create wallet of new player", zap.Error(err))
		}
		span.End()
	}()

	var req CreateNewPlayerWalletRequest

	err = ctx.ShouldBindQuery(&req.WalletUserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.CreateWalletOfNewPlayer(sctx, req)
	if err != nil {
		c.logger.Error("Failed to create wallet of new player",
			zap.String("user_id", req.UserID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *WalletController) GetWalletOfPlayer(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.wallet.GetWalletOfPlayer")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to get wallet of player", zap.Error(err))
		}
		span.End()
	}()

	var req GetWalletOfPlayerRequest

	err = ctx.ShouldBindQuery(&req.WalletUserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	wallet, err := c.service.GetWalletOfPlayer(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get wallet of player",
			zap.String("user_id", req.UserID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, wallet)
}

func (c *WalletController) UpdateWalletOfPlayer(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.wallet.UpdateWalletOfPlayer")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to update wallet of player", zap.Error(err))
		}
		span.End()
	}()

	var req UpdateWalletOfPlayerRequest

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	updatedWallet, err := c.service.UpdateWalletOfPlayer(sctx, req)
	if err != nil {
		c.logger.Error("Failed to update wallet of player",
			zap.String("user_id", req.UserID.String()),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, updatedWallet)
}

func (c *WalletController) DeleteWalletOfPlayer(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.wallet.DeleteWalletOfPlayer")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to delete wallet of player", zap.Error(err))
		}
		span.End()
	}()

	var req DeleteWalletOfPlayerRequest

	if err := ctx.ShouldBindQuery(&req.WalletUserID); err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.DeleteWalletOfPlayer(sctx, req)
	if err != nil {
		c.logger.Error("Failed to delete wallet of player",
			zap.String("user_id", req.UserID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
