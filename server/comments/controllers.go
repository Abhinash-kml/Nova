package comments

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
	sctx, span := c.tracer.Start(ctx.Request.Context(), "comments.controller.getall")
	defer span.End()

	var dto GetAllDTO

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	decodedCursor, err := utils.DecodeCursor(dto.Cursor)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	comments, err := c.service.GetAll(sctx, decodedCursor, dto.Limit)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(comments))
}

func (c *Controller) Get(ctx *gin.Context) {
	sctx, span := c.tracer.Start(ctx.Request.Context(), "comments.controller.get")
	defer span.End()

	var dto GetDTO

	if err := ctx.ShouldBindUri(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	parsedUuid, _ := uuid.Parse(dto.Id)

	comment, err := c.service.GetById(sctx, parsedUuid)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, comment)
}

func (c *Controller) Create(ctx *gin.Context) {
	sctx, span := c.tracer.Start(ctx.Request.Context(), "comments.controller.create")
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
	sctx, span := c.tracer.Start(ctx.Request.Context(), "comments.controller.modify")
	defer span.End()

	var dto UpdateDTO

	if err := ctx.ShouldBindUri(&dto.CommentId); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindWith(&dto.Body, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.Update(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := c.tracer.Start(ctx.Request.Context(), "comments.controller.delete")
	defer span.End()

	var dto DeleteDTO

	if err := ctx.ShouldBindUri(&dto.CommentId); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindQuery(&dto.DeleteOptions); err != nil {
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

func (c *Controller) Replace(ctx *gin.Context) {
	sctx, span := c.tracer.Start(ctx.Request.Context(), "comments.controller.replace")
	defer span.End()

	var dto ReplaceDTO

	if err := ctx.ShouldBindUri(&dto.CommentId); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindWith(&dto.Body, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.Replace(sctx, dto)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) BulkAdd(ctx *gin.Context) {
	sctx, span := c.tracer.Start(ctx.Request.Context(), "comments.controller.bulkadd")
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
	sctx, span := c.tracer.Start(ctx.Request.Context(), "comments.controller.bulkmodify")
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
	sctx, span := c.tracer.Start(ctx.Request.Context(), "comments.controller.bulkdelete")
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
