package clans

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Controller struct {
	service Service
	config  *config.Config
	logger  *zap.Logger
}

func NewController(s Service, c *config.Config, l *zap.Logger) *Controller {
	return &Controller{
		config:  c,
		service: s,
		logger:  l,
	}
}

func (c *Controller) GetAll(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.getall")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req GetAllRequest

	err = ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	decodedCursor, err := utils.DecodeCursorUUID(req.Cursor)
	if err != nil {
		ctx.Error(err)
		return
	}

	clans, err := c.service.GetAll(sctx, decodedCursor, req.Limit)
	if err != nil {
		c.logger.Error("Failed to get all clans",
			zap.String("cursor", req.Cursor),
			zap.Int("limit", req.Limit),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(clans))
}

func (c *Controller) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req GetRequest

	err = ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("clan_id", req.Id))

	parsedID, _ := uuid.Parse(req.Id)

	clan, err := c.service.GetById(sctx, parsedID)
	if err != nil {
		c.logger.Error("Failed to get clan",
			zap.String("clan_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, clan)
}

func (c *Controller) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req CreateRequest

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	clan, err := c.service.Add(sctx, req)
	if err != nil {
		c.logger.Error("Failed to create clan", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, clan)
}

func (c *Controller) Modify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.modify")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req UpdateRequest

	err = ctx.ShouldBindUri(&req.ClanId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("clan_id", req.Id))

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	modifiedClan, err := c.service.Update(sctx, req)
	if err != nil {
		c.logger.Error("Failed to update clan",
			zap.String("clan_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, modifiedClan)
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req DeleteRequest

	err = ctx.ShouldBindUri(&req.ClanId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("clan_id", req.Id))

	err = ctx.ShouldBindQuery(&req.DeleteOptions)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.Delete(sctx, req)
	if err != nil {
		c.logger.Error("Failed to delete clan",
			zap.String("clan_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, response)
}

func (c *Controller) AddMember(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.addmember")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Error("Failed to add member", zap.Error(err))
		}
		span.End()
	}()

	var req ActionRequest
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.AddMember(sctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, response)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) RemoveMember(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.removemember")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Error("Failed to remove member", zap.Error(err))
		}
		span.End()
	}()

	var req ActionRequest
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.RemoveMember(sctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, response)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) ListAllOutogingInvites(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.removemember")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Error("Failed to remove member", zap.Error(err))
		}
		span.End()
	}()

	var req ActionRequest
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.ListAllOutogingInvites(sctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) InviteMember(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.invitemember")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Error("Failed to invite member", zap.Error(err))
		}
		span.End()
	}()

	var req ActionRequest
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.InviteMember(sctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, response)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) AcceptInvite(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.acceptinvite.")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Error("Failed to accept invite", zap.Error(err))
		}
		span.End()
	}()

	var req ActionRequest
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.AcceptInvite(sctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, response)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) RejectInvite(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.rejectinvite")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Error("Failed to reject invite", zap.Error(err))
		}
		span.End()
	}()

	var req ActionRequest
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.RejectInvite(sctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, response)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) ListAllJoinRequests(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.listalljoinrequests")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Error("Failed to list all join requests", zap.Error(err))
		}
		span.End()
	}()

	var req ActionRequest
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.ListAllJoinRequests(sctx, req)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) SendJoinRequest(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.sendjoinrequest")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Error("Failed to send join request", zap.Error(err))
		}
		span.End()
	}()

	var req ActionRequest
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.SendJoinRequest(sctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, response)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) AcceptJoinRequest(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.acceptjoinrequest")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Error("Failed to accept join request", zap.Error(err))
		}
		span.End()
	}()

	var req ActionRequest
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.AcceptJoinRequest(sctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, response)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) RejectJoinRequest(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "clans.controller.rejectjoinrequest")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Error("Failed to reject join request", zap.Error(err))
		}
		span.End()
	}()

	var req ActionRequest
	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.RejectJoinRequest(sctx, req)
	if err != nil {
		ctx.JSON(http.StatusOK, response)
		return
	}

	ctx.JSON(http.StatusOK, response)
}
