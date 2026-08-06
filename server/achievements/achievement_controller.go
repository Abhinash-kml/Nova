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
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto CreateAchievementDTO

	err = ctx.ShouldBindJSON(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.Create(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to create achievement", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	dto := GetAchievementDTO{
		Id: id,
	}

	response, err := c.service.Get(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get achievement",
			zap.Int("achievement_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) Update(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.update")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	var dto UpdateAchievementDTO

	err = ctx.ShouldBindJSON(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	dto.Id = id

	err = c.service.Update(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to update achievement",
			zap.Int("achievement_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		ctx.Error(err)
		return
	}

	dto := DeleteAchievementDTO{
		Id: id,
	}

	err = c.service.Delete(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to delete achievement",
			zap.Int("achievement_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
