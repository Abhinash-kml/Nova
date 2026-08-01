package economy

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/utils"
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
	defer span.End()

	var dto CreateNewPlayerWalletDTO

	if err := ctx.ShouldBindQuery(&dto.WalletUserID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.CreateWalletOfNewPlayer(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *WalletController) GetWalletOfPlayer(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.wallet.getwalletofplayer")
	defer span.End()

	var dto GetWalletOfPlayerDTO

	if err := ctx.ShouldBindQuery(&dto.WalletUserID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	wallet, err := c.service.GetWalletOfPlayer(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, wallet)
}

func (c *WalletController) UpdateWalletOfPlayer(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.wallet.updatewalletofplayer")
	defer span.End()

	var dto UpdateWalletOfPlayerDTO

	if err := ctx.ShouldBindBodyWithJSON(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	updatedWallet, err := c.service.UpdateWalletOfPlayer(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, updatedWallet)
}

func (c *WalletController) DeleteWalletOfPlayer(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.wallet.deletewalletofplayer")
	defer span.End()

	var dto DeleteWalletOfPlayerDTO

	if err := ctx.ShouldBindQuery(&dto.WalletUserID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.DeleteWalletOfPlayer(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
