package leaderboard

import (
	"net/http"
	"strconv"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
	sctx, span := tracer.Start(ctx.Request.Context(), "leaderboard.controller.getall")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	var dto GetAllDTO
	err = ctx.ShouldBindQuery(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	decodedCursor, err := utils.DecodeCursorUUID(dto.Cursor)
	if err != nil {
		ctx.Error(err)
		return
	}

	if dto.Limit == 0 {
		dto.Limit = 10
	}

	span.SetAttributes(attribute.String("cursor", decodedCursor.String()))
	span.SetAttributes(attribute.Int("limit", dto.Limit))

	leaderboards, err := c.service.GetAll(sctx, decodedCursor, dto.Limit)
	if err != nil {
		c.logger.Error("Failed to get all leaderboards",
			zap.String("cursor", decodedCursor.String()),
			zap.Int("limit", dto.Limit),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(leaderboards))
}

func (c *Controller) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "leaderboard.controller.get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	var dto GetDTO

	err = ctx.ShouldBindUri(&dto.LeaderboardId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("leaderboard_id", dto.Id))

	leaderboardId, _ := uuid.Parse(dto.Id)
	leaderboard, err := c.service.Get(sctx, leaderboardId)
	if err != nil {
		c.logger.Error("Failed to get leaderboard",
			zap.String("leaderboard_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, leaderboard)
}

func (c *Controller) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "leaderboard.controller.create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	var dto CreateDTO

	err = ctx.ShouldBindWith(&dto, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	leaderboard, err := c.service.Create(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to create leaderboard", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, leaderboard)
}

func (c *Controller) Modify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx, "leaderboard.controller.modify")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	var dto ModifyDTO

	err = ctx.ShouldBindUri(&dto.LeaderboardId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("leaderboard_id", dto.Id))

	err = ctx.ShouldBindWith(&dto.LeaderboardModifications, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	modifiedLeaderboard, err := c.service.Modify(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to update leaderboard",
			zap.String("leaderboard_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, modifiedLeaderboard)
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "leaderboard.controller.delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	var dto DeleteDTO

	err = ctx.ShouldBindUri(&dto.LeaderboardId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("leaderboard_id", dto.Id))

	_, err = c.service.Delete(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to delete leaderboard",
			zap.String("leaderboard_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) GetScore(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "leaderboard.controller.getscore")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	var dto GetScoreDTO

	err = ctx.ShouldBindUri(&dto.LeaderboardId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("leaderboard_id", dto.Id))

	limit := ctx.DefaultQuery("limit", "100")
	limitNum, _ := strconv.Atoi(limit)

	span.SetAttributes(attribute.Int("limit", limitNum))

	dto.Limit = limitNum

	scores, err := c.service.GetScore(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get score from leaderboard",
			zap.String("leaderboard_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, scores)
}

func (c *Controller) UpdateScore(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "leaderboard.controller.updatescore")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()
	var dto UpdateScoreDTO

	err = ctx.ShouldBindUri(&dto.LeaderboardId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("leaderboard_id", dto.Id))

	err = ctx.ShouldBindQuery(&dto.UpdateOptions)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("operator", dto.AggregateType))

	err = ctx.ShouldBindWith(&dto.ScoreDTO, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.UpdateScore(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to update score in leaderboard",
			zap.String("leaderboard_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) DeleteScore(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "leaderboard.controller.deletescore")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
	}()

	var dto DeleteScoreDTO

	err = ctx.ShouldBindUri(&dto.LeaderboardId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("leaderboard_id", dto.LeaderboardId.Id))

	if err := ctx.ShouldBindQuery(&dto.UserId); err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("user_id", dto.UserId.Id))

	err = c.service.DeleteScore(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to delete score in leaderboard",
			zap.String("leaderboard_id", dto.LeaderboardId.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
