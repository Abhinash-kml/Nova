package stats

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
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
		service: s,
		config:  c,
		logger:  l,
	}
}

func (c *Controller) GetAll(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.getall")
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

	cursor, err := utils.DecodeCursor(dto.Cursor)
	if err != nil {
		ctx.Error(err)
		return
	}

	stats, err := c.service.GetAll(sctx, cursor, dto.Limit)
	if err != nil {
		c.logger.Error("Failed to get all stats",
			zap.String("cursor", dto.Cursor),
			zap.Int("limit", dto.Limit),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(stats))
}

func (c *Controller) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var data GetDTO

	err = ctx.ShouldBindUri(&data)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("stat_id", data.Id))

	stat, err := c.service.GetById(sctx, data.Id)
	if err != nil {
		c.logger.Error("Failed to get stat",
			zap.Int("stat_id", data.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, stat)
}

func (c *Controller) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto CreateDTO

	err = ctx.ShouldBindWith(&dto, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	stat, err := c.service.Add(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to create stat", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, stat)
}

func (c *Controller) Modify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.modify")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto UpdateDTO

	err = ctx.ShouldBindUri(&dto.StatsId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("stat_id", dto.Id))

	err = ctx.ShouldBindWith(&dto, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	modifiedStat, err := c.service.Update(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to update stat",
			zap.Int("stat_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, modifiedStat)
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto DeleteDTO

	err = ctx.ShouldBindUri(&dto.StatsId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("stat_id", dto.Id))

	_, err = c.service.Delete(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to delete stat",
			zap.Int("stat_it", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Replace(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.replace")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto ReplaceDTO

	err = ctx.ShouldBindUri(&dto.StatsId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("stat_id", dto.Id))

	err = ctx.ShouldBindWith(&dto.ReplacementData, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	replacedStat, err := c.service.Replace(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to replace stat",
			zap.Int("stat_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, replacedStat)
}

func (c *Controller) UpdatePlayerStats(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.updateplayerstats")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto UpdatePlayerStatDTO

	err = ctx.ShouldBindUri(&dto.UserId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(
		attribute.String("user_id", dto.Id),
		attribute.String("stat_id", dto.StatId),
	)

	err = ctx.ShouldBindWith(&dto.IncomingPlayerStat, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.UpdatePlayerStats(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to update player stat",
			zap.String("user_id", dto.Id),
			zap.String("stat_it", dto.StatId))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusOK)
}

func (c *Controller) DeletePlayerStats(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.deleteplayerstats")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto DeletePlayerStatDTO

	err = ctx.ShouldBindUri(&dto.UserId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("user_id", dto.Id))

	err = c.service.DeletePlayerStats(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to delete player stats", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) GetPlayerStats(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.updateplayerstats")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto GetPlayerStatDTO

	err = ctx.ShouldBindUri(&dto.UserId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("user_id", dto.Id))

	playerStats, err := c.service.GetPlayerStats(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get player stats", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, playerStats)
}
