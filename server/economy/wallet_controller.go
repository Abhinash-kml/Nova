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
	config  *config.Config
	logger  *zap.Logger
}

func NewWalletController(service WalletService, c *config.Config, l *zap.Logger) *WalletController {
	return &WalletController{
		service: service,
		config:  c,
		logger:  l,
	}
}

func (c *WalletController) CreateWalletOfNewPlayer(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.wallet.createwalletofnewplayer")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto CreateNewPlayerWalletDTO

	err = ctx.ShouldBindQuery(&dto.WalletUserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.CreateWalletOfNewPlayer(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to create wallet of new player",
			zap.String("user_id", dto.UserID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *WalletController) GetWalletOfPlayer(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.wallet.getwalletofplayer")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto GetWalletOfPlayerDTO

	err = ctx.ShouldBindQuery(&dto.WalletUserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	wallet, err := c.service.GetWalletOfPlayer(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get wallet of player",
			zap.String("user_id", dto.UserID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, wallet)
}

func (c *WalletController) UpdateWalletOfPlayer(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.wallet.updatewalletofplayer")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto UpdateWalletOfPlayerDTO

	err = ctx.ShouldBindBodyWithJSON(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	updatedWallet, err := c.service.UpdateWalletOfPlayer(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to update wallet of player",
			zap.String("user_id", dto.UserID.String()),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, updatedWallet)
}

func (c *WalletController) DeleteWalletOfPlayer(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.wallet.deletewalletofplayer")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto DeleteWalletOfPlayerDTO

	if err := ctx.ShouldBindQuery(&dto.WalletUserID); err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.DeleteWalletOfPlayer(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to delete wallet of player",
			zap.String("user_id", dto.UserID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
