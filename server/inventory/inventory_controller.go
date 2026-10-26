package inventory

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type InventoryController struct {
	service InventoryService
	config  *config.InventoryConfig
	logger  *zap.Logger
}

func NewInventoryController(s InventoryService, c *config.InventoryConfig, l *zap.Logger) *InventoryController {
	return &InventoryController{
		config:  c,
		service: s,
		logger:  l,
	}
}

func (c *InventoryController) GetInventoryOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.getinventoryofuser")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req GetInventoryOfUserRequest
	err = ctx.ShouldBindQuery(&req.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	inventory, err := c.service.GetInventoryOfUser(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get inventory of user",
			zap.String("user_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, inventory)
}

func (c *InventoryController) DeleteInventoryOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.deleteinventoryofuser")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req DeleteInventoryOfUserRequest

	err = ctx.ShouldBindQuery(&req.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	_, err = c.service.DeleteInventoryOfUser(sctx, req)
	if err != nil {
		c.logger.Error("Failed to delete inventory of user",
			zap.String("user_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *InventoryController) GetInventoryItemOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.getinventoryofuser")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req GetInventoryItemOfUserRequest

	err = ctx.ShouldBindUri(&req.ItemID)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = ctx.ShouldBindQuery(&req.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	item, err := c.service.GetInventoryItemOfUser(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get inventory item of user",
			zap.String("user_id", req.Id),
			zap.String("item_id", req.ItemId),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *InventoryController) AddInventoryItemOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.getinventoryofuser")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req AddInventoryItemOfUserRequest

	err = ctx.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	item, err := c.service.AddInventoryItemOfUser(sctx, req)
	if err != nil {
		c.logger.Error("Failed to add inventory item of user",
			zap.String("user_id", req.UserID.String()),
			zap.String("item_id", req.ItemID.String()),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *InventoryController) UpdateInventoryItemOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.getinventoryofuser")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req UpdateInventoryItemOfUserRequest

	err = ctx.ShouldBindUri(&req.ItemID)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = ctx.ShouldBindQuery(&req.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = ctx.ShouldBindBodyWith(&req.UpdateInventoryData, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	updatedItem, err := c.service.UpdateInventoryItemOfUser(sctx, req)
	if err != nil {
		c.logger.Error("Failed to update inventory item of user",
			zap.String("user_id", req.Id),
			zap.String("item_id", req.ItemId),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, updatedItem)
}

func (c *InventoryController) DeleteInventoryItemOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.getinventoryofuser")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req DeleteInventoryItemOfUserRequest

	err = ctx.ShouldBindUri(&req.ItemID)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = ctx.ShouldBindQuery(&req.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.DeleteInventoryItemOfUser(sctx, req)
	if err != nil {
		c.logger.Error("Failed to delete inventory item of user",
			zap.String("user_id", req.Id),
			zap.String("item_id", req.ItemId),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
