package economy

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type CurrencyController struct {
	service CurrencyService
	config  *config.Config
	logger  *zap.Logger
}

func NBewCurrencyController(service CurrencyService, c *config.Config, l *zap.Logger) *CurrencyController {
	return &CurrencyController{
		service: service,
		config:  c,
		logger:  l,
	}
}

func (c *CurrencyController) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.")
	defer span.End()

	var dto GetCurrencyDTO

	if err := ctx.ShouldBindUri(&dto.CurrencyID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	currency, err := c.service.Get(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, currency)
}

func (c *CurrencyController) GetAll(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.")
	defer span.End()

	var dto GetAllCurrencyDTO

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	currencies, err := c.service.GetAll(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, currencies)
}

func (c *CurrencyController) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.")
	defer span.End()

	var dto CreateCurrencyDTO

	if err := ctx.ShouldBindBodyWithJSON(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	currency, err := c.service.Create(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, currency)
}

func (c *CurrencyController) Update(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.")
	defer span.End()

	var dto UpdateCurrencyDTO

	if err := ctx.ShouldBindUri(&dto.CurrencyID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindBodyWithJSON(&dto.UpdateData); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	updatedCurrency, err := c.service.Update(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, updatedCurrency)
}

func (c *CurrencyController) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.")
	defer span.End()

	var dto DeleetCurrencyDTO

	if err := ctx.ShouldBindUri(&dto.CurrencyID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	_, err := c.service.Delete(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}
}
