package inventory

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/config"
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
	config  *config.Config
	logger  *zap.Logger
}

func NewItemsController(s ItemsService, c *config.Config, l *zap.Logger) *ItemsController {
	return &ItemsController{
		config:  c,
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

	var dto GetAllDTO

	err = ctx.ShouldBindQuery(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	items, err := c.service.GetAll(sctx, dto.Cursor, dto.Limit)
	if err != nil {
		c.logger.Error("Failed to get all items",
			zap.Int("cursor", dto.Cursor),
			zap.Int("limit", dto.Limit),
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

	var dto GetItemDTO

	err = ctx.ShouldBindUri(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("item_id", dto.Id))

	itemID, _ := uuid.Parse(dto.Id)

	item, err := c.service.GetById(sctx, itemID)
	if err != nil {
		c.logger.Error("Failed to get item",
			zap.String("item_id", dto.Id),
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

	var dto CreateItemDTO

	err = ctx.ShouldBindWith(&dto, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	item, err := c.service.Add(sctx, dto)
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

	var dto UpdateItemDTO

	err = ctx.ShouldBindUri(&dto.ItemId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("item_id", dto.Id))

	err = ctx.ShouldBindWith(&dto.ItemUpdateData, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	item, err := c.service.Update(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to update item",
			zap.String("item_id", dto.Id),
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

	var dto ReplaceItemDTO

	err = ctx.ShouldBindUri(&dto.ItemId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("item_id", dto.Id))

	err = ctx.ShouldBindWith(&dto.ReplacementData, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	item, err := c.service.Replace(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to replace item",
			zap.String("item_id", dto.Id),
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

	var dto DeleteItemDTO

	err = ctx.ShouldBindUri(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("item_id", dto.Id))

	_, err = c.service.Delete(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to delete item",
			zap.String("item_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
