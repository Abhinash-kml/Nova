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
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req GetAllItemRequest

	err = ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	items, err := c.service.GetAll(sctx, req.Cursor, req.Limit)
	if err != nil {
		c.logger.Error("Failed to get all items",
			zap.Int("cursor", req.Cursor),
			zap.Int("limit", req.Limit),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(items))
}

func (c *ItemsController) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.items.controller.get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req GetItemRequest

	err = ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("item_id", req.Id))

	itemID, _ := uuid.Parse(req.Id)

	item, err := c.service.GetById(sctx, itemID)
	if err != nil {
		c.logger.Error("Failed to get item",
			zap.String("item_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *ItemsController) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.items.controller.create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req CreateItemRequest

	err = ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	item, err := c.service.Add(sctx, req)
	if err != nil {
		c.logger.Error("Failed to create item", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, item)
}

func (c *ItemsController) Modify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.items.controller.modify")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req UpdateItemRequest

	err = ctx.ShouldBindUri(&req.ItemId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("item_id", req.Id))

	err = ctx.ShouldBindWith(&req.ItemUpdateData, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	item, err := c.service.Update(sctx, req)
	if err != nil {
		c.logger.Error("Failed to update item",
			zap.String("item_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *ItemsController) Replace(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.items.controller.replace")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req ReplaceItemRequest

	err = ctx.ShouldBindUri(&req.ItemId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("item_id", req.Id))

	err = ctx.ShouldBindWith(&req.ReplacementData, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	item, err := c.service.Replace(sctx, req)
	if err != nil {
		c.logger.Error("Failed to replace item",
			zap.String("item_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *ItemsController) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.items.controller.delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req DeleteItemRequest

	err = ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("item_id", req.Id))

	_, err = c.service.Delete(sctx, req)
	if err != nil {
		c.logger.Error("Failed to delete item",
			zap.String("item_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
