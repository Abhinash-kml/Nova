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
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.addfriend")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto AddFriendDTO

	err = ctx.ShouldBindBodyWithJSON(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.AddFriend(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to add friend",
			zap.String("user_id", dto.UserID),
			zap.String("target_id", dto.TargetID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) RemoveFriend(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.removefriend")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto RemoveFriendDTO

	err = ctx.ShouldBindBodyWithJSON(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.RemoveFriend(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to remove friend",
			zap.String("user_id", dto.UserID),
			zap.String("target_id", dto.TargetID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) AcceptFriendRequest(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.acceptfriendrequest")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto AcceptFriendRequestDTO

	err = ctx.ShouldBindBodyWithJSON(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.AcceptFriendRequest(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to accept friend request",
			zap.String("user_id", dto.UserID),
			zap.String("target_id", dto.TargetID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) RejectFriendRequest(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.rejectfriendrequest")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto RejectFriendRequestDTO

	err = ctx.ShouldBindBodyWithJSON(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.RejectFriendRequest(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to reject friend request",
			zap.String("user_id", dto.UserID),
			zap.String("target_id", dto.TargetID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) GetAllIncomingFriendRequests(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.getallincomingfriendrequests")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto GetIncomingRequestsOfUserDTO

	err = ctx.ShouldBindQuery(&dto.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	incomingRequestIDs, err := c.service.GetAllIncomingFriendRequests(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get all friend requests",
			zap.String("user_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, incomingRequestIDs)
}

func (c *Controller) GetAllOutgoingFriendRequests(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.getalloutgoingfriendrequests")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto GetOutgoingRequestsOfUserDTO

	err = ctx.ShouldBindQuery(&dto.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	outgoingRequestIDs, err := c.service.GetAllOutgoingFriendRequests(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get all outgoing friend requests",
			zap.String("user_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, outgoingRequestIDs)
}

func (c *Controller) BlockUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.blockuser")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto BlockUserDTO

	err = ctx.ShouldBindBodyWithJSON(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.BlockUser(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to block user",
			zap.String("user_id", dto.UserID),
			zap.String("target_id", dto.TargetID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) UnblockUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.unblockuser")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto UnblockUserDTO

	err = ctx.ShouldBindBodyWithJSON(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.UnblockUser(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to unblock user",
			zap.String("user_id", dto.UserID),
			zap.String("target_id", dto.TargetID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) GetAllFriends(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.getallfriends")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto GetAllFriendsDTO

	err = ctx.ShouldBindQuery(&dto.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	friendIDs, err := c.service.GetAllFriends(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get all friends",
			zap.String("user_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, friendIDs)
}

func (c *Controller) GetAllBlocked(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.getallblocked")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto GetAllBlockedDTO

	err = ctx.ShouldBindQuery(&dto.UserID)
	if err != nil {
		ctx.Error(err)
		return
	}

	blockedIDs, err := c.service.GetAllBlocked(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get all blocked users",
			zap.String("user_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, blockedIDs)
}

func (c *Controller) GetMutualFriends(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.getmutualfriends")
	var err error
	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto GetMutualFriendsDTO

	err = ctx.ShouldBindBodyWithJSON(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	mutualFriendIDs, err := c.service.GetMutualFriends(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get mutual friends",
			zap.String("user_id", dto.UserOneID),
			zap.String("target_id", dto.UserTwoID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, mutualFriendIDs)
}
