package social

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/utils"
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
	defer span.End()

	var dto AddFriendDTO

	if err := ctx.ShouldBindBodyWithJSON(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := c.service.AddFriend(sctx, dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) RemoveFriend(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.removefriend")
	defer span.End()

	var dto RemoveFriendDTO

	if err := ctx.ShouldBindBodyWithJSON(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := c.service.RemoveFriend(sctx, dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) AcceptFriendRequest(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.acceptfriendrequest")
	defer span.End()

	var dto AcceptFriendRequestDTO

	if err := ctx.ShouldBindBodyWithJSON(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := c.service.AcceptFriendRequest(sctx, dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) RejectFriendRequest(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.rejectfriendrequest")
	defer span.End()

	var dto RejectFriendRequestDTO

	if err := ctx.ShouldBindBodyWithJSON(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := c.service.RejectFriendRequest(sctx, dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) GetAllIncomingFriendRequests(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.getallincomingfriendrequests")
	defer span.End()

	var dto GetIncomingRequestsOfUserDTO

	if err := ctx.ShouldBindQuery(&dto.UserID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	incomingRequestIDs, err := c.service.GetAllIncomingFriendRequests(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, incomingRequestIDs)
}

func (c *Controller) GetAllOutgoingFriendRequests(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.getalloutgoingfriendrequests")
	defer span.End()

	var dto GetOutgoingRequestsOfUserDTO

	if err := ctx.ShouldBindQuery(&dto.UserID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	outgoingRequestIDs, err := c.service.GetAllOutgoingFriendRequests(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, outgoingRequestIDs)
}

func (c *Controller) BlockUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.blockuser")
	defer span.End()

	var dto BlockUserDTO

	if err := ctx.ShouldBindBodyWithJSON(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := c.service.BlockUser(sctx, dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) UnblockUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.unblockuser")
	defer span.End()

	var dto UnblockUserDTO

	if err := ctx.ShouldBindBodyWithJSON(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := c.service.UnblockUser(sctx, dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) GetAllFriends(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.getallfriends")
	defer span.End()

	var dto GetAllFriendsDTO

	if err := ctx.ShouldBindQuery(&dto.UserID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	friendIDs, err := c.service.GetAllFriends(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, friendIDs)
}

func (c *Controller) GetAllBlocked(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.getallblocked")
	defer span.End()

	var dto GetAllBlockedDTO

	if err := ctx.ShouldBindQuery(&dto.UserID); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	blockedIDs, err := c.service.GetAllBlocked(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, blockedIDs)
}

func (c *Controller) GetMutualFriends(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "social.controller.getmutualfriends")
	defer span.End()

	var dto GetMutualFriendsDTO

	if err := ctx.ShouldBindBodyWithJSON(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	mutualFriendIDs, err := c.service.GetMutualFriends(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, mutualFriendIDs)
}
