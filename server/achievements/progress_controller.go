package achievements

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type ProgressController struct {
	service ProgressService
	logger  *zap.Logger
}

func NewProgressController(service ProgressService, logger *zap.Logger) *ProgressController {
	return &ProgressController{
		service: service,
		logger:  logger,
	}
}

func (c *ProgressController) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.create")
	defer span.End()

	var dto CreateProgressDTO

	if err := ctx.ShouldBindBodyWith(&dto, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	result, err := c.service.Create(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, result)
}

func (c *ProgressController) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.get")
	defer span.End()

	var dto GetProgressDTO

	if err := ctx.ShouldBindUri(&dto.ProgressId); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	result, err := c.service.Get(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *ProgressController) GetOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.getofuser")
	defer span.End()

	var dto GetProgressOfUserDTO

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	result, err := c.service.GetOfUser(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *ProgressController) Update(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.update")
	defer span.End()

	var dto UpdateProgressDTO

	if err := ctx.ShouldBindUri(&dto.ProgressId); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
	}

	span.SetAttributes(attribute.Int("progress_id", dto.Id))

	if err := ctx.ShouldBindBodyWith(&dto.ProgressUpdateData, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := c.service.Update(sctx, dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *ProgressController) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.progress.controller.delete")
	defer span.End()

	var dto DeleteProgressDTO

	if err := ctx.ShouldBindUri(&dto.ProgressId); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	span.SetAttributes(attribute.Int("progress_id", dto.Id))

	if _, err := c.service.Delete(sctx, dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
