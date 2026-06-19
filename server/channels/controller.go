package channels

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Controller struct {
	service Service
	logger  *zap.Logger
	tracer  trace.Tracer
}

func NewController(s Service, l *zap.Logger, t trace.Tracer) *Controller {
	return &Controller{
		service: s,
		logger:  l,
		tracer:  t,
	}
}

func (c *Controller) GetAll(ctx *gin.Context) {
	sctx, span := c.tracer.Start(ctx.Request.Context(), "channels.controller.getall")
	defer span.End()

	var dto GetAllDTO

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	decodedCursor, err := utils.DecodeCursor(dto.Cursor)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	channels, err := c.service.GetAll(sctx, decodedCursor, dto.Limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(channels))
}

func (c *Controller) Get(ctx *gin.Context) {
	sctx, span := c.tracer.Start(ctx.Request.Context(), "channels.controller.get")
	defer span.End()

	var dto GetDTO

	if err := ctx.ShouldBindUri(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	parsedId, _ := uuid.Parse(dto.Id)

	channel, err := c.service.GetById(sctx, parsedId)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, channel)
}

func (c *Controller) Create(ctx *gin.Context) {
	sctx, span := c.tracer.Start(ctx.Request.Context(), "channels.controller.create")
	defer span.End()

	var dto CreateDTO

	if err := ctx.ShouldBindWith(&dto, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.Add(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *Controller) Modify(ctx *gin.Context) {
	sctx, span := c.tracer.Start(ctx.Request.Context(), "channels.controller.modify")
	defer span.End()

	var dto UpdateDTO

	if err := ctx.ShouldBindUri(&dto.ChannelId); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindWith(&dto.ChannelModifications, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.Modify(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := c.tracer.Start(ctx.Request.Context(), "channels.controller.delete")
	defer span.End()

	var dto DeleteDTO

	if err := ctx.ShouldBindUri(&dto.ChannelId); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.Delete(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) BulkAdd(ctx *gin.Context) {
	sctx, span := c.tracer.Start(ctx.Request.Context(), "channels.controller.bulkadd")
	defer span.End()

	var dto BulkCreateDTO

	if err := ctx.ShouldBindWith(&dto, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.BulkAdd(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *Controller) BulkModify(ctx *gin.Context) {
	sctx, span := c.tracer.Start(ctx.Request.Context(), "channels.controller.bulkmodify")
	defer span.End()

	var dto BulkModifyDTO

	if err := ctx.ShouldBindWith(&dto, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.BulkModify(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) BulkDelete(ctx *gin.Context) {
	sctx, span := c.tracer.Start(ctx.Request.Context(), "channels.controller.bulkdelete")
	defer span.End()

	var dto BulkDeleteDTO

	if err := ctx.ShouldBindWith(&dto, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.BulkDelete(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
