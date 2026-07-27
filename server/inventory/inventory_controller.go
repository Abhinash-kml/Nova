package inventory

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type InventoryController struct {
	service InventoryService
	config  *config.Config
	logger  *zap.Logger
}

func NewInventoryController(s InventoryService, c *config.Config, l *zap.Logger) *InventoryController {
	return &InventoryController{
		config:  c,
		service: s,
		logger:  l,
	}
}

func (c *InventoryController) GetInventoryOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.getinventoryofuser")
	defer span.End()

	var dto GetInventoryOfUserDTO
	if err := ctx.ShouldBindQuery(&dto.UserID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	inventory, err := c.service.GetInventoryOfUser(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, inventory)
}

func (c *InventoryController) DeleteInventoryOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.deleteinventoryofuser")
	defer span.End()

	var dto DeleteInventoryOfUserDTO
	if err := ctx.ShouldBindQuery(&dto.UserID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	_, err := c.service.DeleteInventoryOfUser(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *InventoryController) GetInventoryItemOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.getinventoryofuser")
	defer span.End()

	var dto GetInventoryItemOfUserDTO

	if err := ctx.ShouldBindUri(&dto.ItemID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindQuery(&dto.UserID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	item, err := c.service.GetInventoryItemOfUser(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *InventoryController) AddInventoryItemOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.getinventoryofuser")
	defer span.End()

	var dto AddInventoryItemOfUserDTO

	if err := ctx.ShouldBindBodyWith(&dto, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	item, err := c.service.AddInventoryItemOfUser(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, item)
}

func (c *InventoryController) UpdateInventoryItemOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.getinventoryofuser")
	defer span.End()

	var dto UpdateInventoryItemOfUserDTO

	if err := ctx.ShouldBindUri(&dto.ItemID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindQuery(&dto.UserID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindBodyWith(&dto.UpdateInventoryData, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	updatedItem, err := c.service.UpdateInventoryItemOfUser(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, updatedItem)
}

func (c *InventoryController) DeleteInventoryItemOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "inventory.getinventoryofuser")
	defer span.End()

	var dto DeleteInventoryItemOfUserDTO

	if err := ctx.ShouldBindUri(&dto.ItemID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindQuery(&dto.UserID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.DeleteInventoryItemOfUser(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
