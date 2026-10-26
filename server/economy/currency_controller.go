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
	config  *config.CurrencyConfig
	logger  *zap.Logger
}

func NewCurrencyController(service CurrencyService, c *config.CurrencyConfig, l *zap.Logger) *CurrencyController {
	return &CurrencyController{
		service: service,
		config:  c,
		logger:  l,
	}
}

func (c *CurrencyController) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.controller.Get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to get currency", zap.Error(err))
		}
		span.End()
	}()

	var req GetCurrencyRequest

	err = ctx.ShouldBindUri(&req.CurrencyID)
	if err != nil {
		ctx.Error(err)
		return
	}

	currency, err := c.service.Get(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get currency data",
			zap.String("currency_id", req.ID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, currency)
}

func (c *CurrencyController) GetAll(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.controller.GetAll")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to get all currency", zap.Error(err))
		}
		span.End()
	}()

	var req GetAllCurrencyRequest

	err = ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	currencies, err := c.service.GetAll(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get all currency",
			zap.Int("limit", req.Limit),
			zap.Int("cursor", req.Cursor),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, currencies)
}

func (c *CurrencyController) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.controller.Create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to create currency", zap.Error(err))
		}
		span.End()
	}()

	var req CreateCurrencyRequest

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	currency, err := c.service.Create(sctx, req)
	if err != nil {
		c.logger.Error("Failed to create currency", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, currency)
}

func (c *CurrencyController) Update(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.controller.Update")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to update currency", zap.Error(err))
		}
		span.End()
	}()

	var req UpdateCurrencyRequest

	err = ctx.ShouldBindUri(&req.CurrencyID)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = ctx.ShouldBindBodyWithJSON(&req.UpdateData)
	if err != nil {
		ctx.Error(err)
		return
	}

	updatedCurrency, err := c.service.Update(sctx, req)
	if err != nil {
		c.logger.Error("Failed to update currency",
			zap.String("currency_id", req.ID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, updatedCurrency)
}

func (c *CurrencyController) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "economy.currency.controller.Delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to delete currency", zap.Error(err))
		}
		span.End()
	}()

	var req DeleteCurrencyRequest

	err = ctx.ShouldBindUri(&req.CurrencyID)
	if err != nil {
		ctx.Error(err)
		return
	}

	_, err = c.service.Delete(sctx, req)
	if err != nil {
		c.logger.Error("Failed to delete currency",
			zap.String("currency_id", req.ID),
			zap.Error(err))
		ctx.Error(err)
		return
	}
}
