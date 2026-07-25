package inventory

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type ItemsController struct {
	service ItemsService
	logger  *zap.Logger
}

func NewItemsController(s ItemsService, l *zap.Logger) *ItemsController {
	return &ItemsController{
		service: s,
		logger:  l,
	}
}

func (c *ItemsController) GetAll(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.items.controller.getall")
	defer span.End()

	var dto GetAllDTO

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	items, err := c.service.GetAll(sctx, dto.Cursor, dto.Limit)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(items))
}

func (c *ItemsController) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.items.controller.get")
	defer span.End()

	var dto GetItemDTO

	if err := ctx.ShouldBindUri(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	span.SetAttributes(attribute.String("item_id", dto.Id))

	itemID, _ := uuid.Parse(dto.Id)

	item, err := c.service.GetById(sctx, itemID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *ItemsController) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.items.controller.create")
	defer span.End()

	var dto CreateItemDTO

	if err := ctx.ShouldBindWith(&dto, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	item, err := c.service.Add(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, item)
}

func (c *ItemsController) Modify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.items.controller.modify")
	defer span.End()

	var dto UpdateItemDTO

	if err := ctx.ShouldBindUri(&dto.ItemId); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	span.SetAttributes(attribute.String("item_id", dto.Id))

	if err := ctx.ShouldBindWith(&dto.ItemUpdateData, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	item, err := c.service.Update(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *ItemsController) Replace(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.items.controller.replace")
	defer span.End()

	var dto ReplaceItemDTO

	if err := ctx.ShouldBindUri(&dto.ItemId); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	span.SetAttributes(attribute.String("itemid", dto.Id))

	if err := ctx.ShouldBindWith(&dto.ReplacementData, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	item, err := c.service.Replace(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *ItemsController) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.items.controller.delete")
	defer span.End()

	var dto DeleteItemDTO

	if err := ctx.ShouldBindUri(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	span.SetAttributes(attribute.String("itemid", dto.Id))

	_, err := c.service.Delete(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
