package entitlements

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Controller struct {
	service Service
	logger  *zap.Logger
}

func NewController(service Service, l *zap.Logger) *Controller {
	return &Controller{
		service: service,
		logger:  l,
	}
}

func (c *Controller) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "entitlements.controller.Create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req CreateRequest

	err = ctx.ShouldBind(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ent, err := c.service.Create(sctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, ent)
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "entitlements.controller.Delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req DeleteRequest

	err = ctx.ShouldBindUri(&req.EntitlementID)
	if err != nil {
		ctx.Error(err)
		return
	}

	deletedID, err := c.service.Delete(sctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, deletedID)
}

func (c *Controller) CheckIfActive(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "entitlements.controller.CheckIfActive")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req CheckIfActiveRequest

	err = ctx.ShouldBindUri(&req.EntitlementID)
	if err != nil {
		ctx.Error(err)
		return
	}

	isActive, err := c.service.CheckIfActive(sctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, isActive)
}

func (c *Controller) GetAllOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "entitlements.controller.GetAllOfUser")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req GetAllOfUserRequest

	err = ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	entitlements, err := c.service.GetAllOfUser(sctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, entitlements)
}
