package economy

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/config"
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
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.controller.get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto GetCurrencyDTO

	err = ctx.ShouldBindUri(&dto.CurrencyID)
	if err != nil {
		ctx.Error(err)
		return
	}

	currency, err := c.service.Get(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get currency data",
			zap.String("currency_id", dto.ID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, currency)
}

func (c *CurrencyController) GetAll(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.controller.getall")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto GetAllCurrencyDTO

	err = ctx.ShouldBindQuery(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	currencies, err := c.service.GetAll(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get all currency",
			zap.Int("limit", dto.Limit),
			zap.Int("cursor", dto.Cursor),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, currencies)
}

func (c *CurrencyController) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.controller.create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto CreateCurrencyDTO

	err = ctx.ShouldBindBodyWithJSON(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	currency, err := c.service.Create(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to create currency", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, currency)
}

func (c *CurrencyController) Update(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.controller.update")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto UpdateCurrencyDTO

	err = ctx.ShouldBindUri(&dto.CurrencyID)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = ctx.ShouldBindBodyWithJSON(&dto.UpdateData)
	if err != nil {
		ctx.Error(err)
		return
	}

	updatedCurrency, err := c.service.Update(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to update currency",
			zap.String("currency_id", dto.ID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, updatedCurrency)
}

func (c *CurrencyController) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.controller.delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto DeleetCurrencyDTO

	err = ctx.ShouldBindUri(&dto.CurrencyID)
	if err != nil {
		ctx.Error(err)
		return
	}

	_, err = c.service.Delete(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to delete currency",
			zap.String("currency_id", dto.ID),
			zap.Error(err))
		ctx.Error(err)
		return
	}
}
