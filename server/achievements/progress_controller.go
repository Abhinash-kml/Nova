package achievements

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type ProgressController struct {
	service ProgressService
	config  *config.Config
	logger  *zap.Logger
}

func NewProgressController(service ProgressService, c *config.Config, logger *zap.Logger) *ProgressController {
	return &ProgressController{
		service: service,
		config:  c,
		logger:  logger,
	}
}

func (c *ProgressController) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.Create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to create achievement progress", zap.Error(err))
		}
		span.End()
	}()

	var req CreateProgressRequest

	err = ctx.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	result, err := c.service.Create(sctx, req)
	if err != nil {
		c.logger.Error("Failed to create achievement progress",
			zap.Int("criteria_id", req.CriterionId),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (c *ProgressController) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.Get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to get achievement progress", zap.Error(err))
		}
		span.End()
	}()

	var req GetProgressRequest

	err = ctx.ShouldBindUri(&req.ProgressId)
	if err != nil {
		ctx.Error(err)
		return
	}

	result, err := c.service.Get(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get achievement progress",
			zap.Int("progress_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *ProgressController) GetOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.GetOfUser")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to get achievement progress of user", zap.Error(err))
		}
		span.End()
	}()

	var req GetProgressOfUserRequest

	err = ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	result, err := c.service.GetOfUser(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get achievement progress of user", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *ProgressController) Update(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.Update")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to update achievement progress", zap.Error(err))
		}
		span.End()
	}()

	var req UpdateProgressRequest

	err = ctx.ShouldBindUri(&req.ProgressId)
	if err != nil {
		ctx.Error(err)
	}

	span.SetAttributes(attribute.Int("progress_id", req.Id))

	err = ctx.ShouldBindBodyWith(&req.ProgressUpdateData, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.Update(sctx, req)
	if err != nil {
		c.logger.Error("Failed to update acheievemnt progress",
			zap.Int("progress_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *ProgressController) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.Delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to delete achievement progress", zap.Error(err))
		}
		span.End()
	}()

	var req DeleteProgressRequest

	err = ctx.ShouldBindUri(&req.ProgressId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("progress_id", req.Id))

	_, err = c.service.Delete(sctx, req)
	if err != nil {
		c.logger.Error("Failed to delete achievement progress",
			zap.Int("progress_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
