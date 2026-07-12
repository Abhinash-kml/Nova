package stats

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Controller struct {
	service Service
	logger  *zap.Logger
}

func NewController(s Service, l *zap.Logger) *Controller {
	return &Controller{
		service: s,
		logger:  l,
	}
}

func (c *Controller) GetAll(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.getall")
	defer span.End()

	var dto GetAllDTO

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	stats, err := c.service.GetAll(sctx, dto.Cursor, dto.Limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(stats))
}

func (c *Controller) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.get")
	defer span.End()

	var data GetDTO

	if err := ctx.ShouldBindUri(&data); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	span.SetAttributes(attribute.Int("statid", data.Id))

	stat, err := c.service.GetById(sctx, data.Id)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, stat)
}

func (c *Controller) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.create")
	defer span.End()

	var dto CreateDTO

	if err := ctx.ShouldBindWith(&dto, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	stat, err := c.service.Add(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, stat)
}

func (c *Controller) Modify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.modify")
	defer span.End()

	var dto UpdateDTO

	if err := ctx.ShouldBindUri(&dto.StatsId); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	span.SetAttributes(attribute.Int("statid", dto.Id))

	if err := ctx.ShouldBindWith(&dto, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	modifiedStat, err := c.service.Update(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, modifiedStat)
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.delete")
	defer span.End()

	var dto DeleteDTO

	if err := ctx.ShouldBindUri(&dto.StatsId); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	span.SetAttributes(attribute.Int("statid", dto.Id))

	_, err := c.service.Delete(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Replace(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.replace")
	defer span.End()

	var dto ReplaceDTO

	if err := ctx.ShouldBindUri(&dto.StatsId); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	span.SetAttributes(attribute.Int("statid", dto.Id))

	if err := ctx.ShouldBindWith(&dto.ReplacementData, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	replacedStat, err := c.service.Replace(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, replacedStat)
}
