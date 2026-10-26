package social

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type Controller struct {
	service Service
	config  *config.Config
	logger  *zap.Logger
}

func NewController(service Service, c *config.Config, l *zap.Logger) *Controller {
	return &Controller{
		service: service,
		config:  c,
		logger:  l,
	}
}

func (c *Controller) AddFriend(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.AddFriend")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req AddFriend_Request

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.AddFriend(sctx, req)
	if err != nil {
		c.logger.Error("Failed to add friend",
			zap.String("user_id", req.UserID),
			zap.String("target_id", req.TargetID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) RemoveFriend(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.RemoveFriend")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req RemoveFriend_Request

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.RemoveFriend(sctx, req)
	if err != nil {
		c.logger.Error("Failed to remove friend",
			zap.String("user_id", req.UserID),
			zap.String("target_id", req.TargetID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) AcceptFriendRequest(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.AcceptFriendRequest")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req AcceptFriendRequest_Request

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.AcceptFriendRequest(sctx, req)
	if err != nil {
		c.logger.Error("Failed to accept friend request",
			zap.String("user_id", req.UserID),
			zap.String("target_id", req.TargetID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) RejectFriendRequest(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.RejectFriendRequest")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req RejectFriendRequest_Request

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.RejectFriendRequest(sctx, req)
	if err != nil {
		c.logger.Error("Failed to reject friend request",
			zap.String("user_id", req.UserID),
			zap.String("target_id", req.TargetID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) GetAllIncomingFriendRequests(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.GetAllIncomingFriendRequests")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req GetIncomingRequestsOfUser_Request

	err = ctx.ShouldBindQuery(&req.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	incomingRequestIDs, err := c.service.GetAllIncomingFriendRequests(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get all friend requests",
			zap.String("user_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, incomingRequestIDs)
}

func (c *Controller) GetAllOutgoingFriendRequests(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.GetAllOutgoingFriendRequests")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req GetOutgoingRequestsOfUser_Request

	err = ctx.ShouldBindQuery(&req.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	outgoingRequestIDs, err := c.service.GetAllOutgoingFriendRequests(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get all outgoing friend requests",
			zap.String("user_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, outgoingRequestIDs)
}

func (c *Controller) BlockUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.BlockUser")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req BlockUser_Request

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.BlockUser(sctx, req)
	if err != nil {
		c.logger.Error("Failed to block user",
			zap.String("user_id", req.UserID),
			zap.String("target_id", req.TargetID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) UnblockUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.UnblockUser")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req UnblockUser_Request

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.UnblockUser(sctx, req)
	if err != nil {
		c.logger.Error("Failed to unblock user",
			zap.String("user_id", req.UserID),
			zap.String("target_id", req.TargetID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) GetAllFriends(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.GetAllFriends")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req GetAllFriends_Request

	err = ctx.ShouldBindQuery(&req.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	friendIDs, err := c.service.GetAllFriends(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get all friends",
			zap.String("user_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, friendIDs)
}

func (c *Controller) GetAllBlocked(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.GetAllBlocked")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req GetAllBlocked_Request

	err = ctx.ShouldBindQuery(&req.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	blockedIDs, err := c.service.GetAllBlocked(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get all blocked users",
			zap.String("user_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, blockedIDs)
}

func (c *Controller) GetMutualFriends(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.GetMutualFriends")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req GetMutualFriends_Request

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	mutualFriendIDs, err := c.service.GetMutualFriends(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get mutual friends",
			zap.String("user_id", req.UserOneID),
			zap.String("target_id", req.UserTwoID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mutualFriendIDs)
}
