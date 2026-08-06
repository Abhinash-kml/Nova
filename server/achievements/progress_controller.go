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
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto CreateProgressDTO

	err = ctx.ShouldBindBodyWith(&dto, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	result, err := c.service.Create(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to create achievement progress",
			zap.Int("criteria_id", dto.CriterionId),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (c *ProgressController) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto GetProgressDTO

	err = ctx.ShouldBindUri(&dto.ProgressId)
	if err != nil {
		ctx.Error(err)
		return
	}

	result, err := c.service.Get(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get achievement progress",
			zap.Int("progress_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *ProgressController) GetOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.getofuser")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto GetProgressOfUserDTO

	err = ctx.ShouldBindQuery(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	result, err := c.service.GetOfUser(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get achievement progress of user", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *ProgressController) Update(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.update")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto UpdateProgressDTO

	err = ctx.ShouldBindUri(&dto.ProgressId)
	if err != nil {
		ctx.Error(err)
	}

	span.SetAttributes(attribute.Int("progress_id", dto.Id))

	err = ctx.ShouldBindBodyWith(&dto.ProgressUpdateData, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.Update(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to update acheievemnt progress",
			zap.Int("progress_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *ProgressController) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto DeleteProgressDTO

	err = ctx.ShouldBindUri(&dto.ProgressId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("progress_id", dto.Id))

	_, err = c.service.Delete(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to delete achievement progress",
			zap.Int("progress_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
