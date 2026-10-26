package achievements

import (
	"net/http"
	"strconv"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Controller struct {
	config  *config.Config
	logger  *zap.Logger
	service Service
}

func NewController(service Service, c *config.Config, l *zap.Logger) *Controller {
	return &Controller{
		config:  c,
		service: service,
		logger:  l,
	}
}

func (c *Controller) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.Create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to create achievement", zap.Error(err))
		}
		span.End()
	}()

	var req CreateAchievementRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.Create(sctx, req)
	if err != nil {
		c.logger.Error("Failed to create achievement", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.Get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to get achievement", zap.Error(err))
		}
		span.End()
	}()

	var req GetAchievementRequest

	err = ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.Get(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get achievement",
			zap.Int("achievement_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) Update(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.Update")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to update achievement", zap.Error(err))
		}
		span.End()
	}()

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	var req UpdateAchievementRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	req.Id = id

	err = c.service.Update(sctx, req)
	if err != nil {
		c.logger.Error("Failed to update achievement",
			zap.Int("achievement_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.Delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to delete achievement", zap.Error(err))
		}
		span.End()
	}()

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	req := DeleteAchievementRequest{
		Id: id,
	}

	err = c.service.Delete(sctx, req)
	if err != nil {
		c.logger.Error("Failed to delete achievement",
			zap.Int("achievement_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
