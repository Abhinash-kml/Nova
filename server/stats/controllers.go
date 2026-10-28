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
	config  *config.StatsConfig
	logger  *zap.Logger
}

func NewController(s Service, c *config.StatsConfig, l *zap.Logger) *Controller {
	return &Controller{
		service: s,
		config:  c,
		logger:  l,
	}
}

func (c *Controller) GetAll(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.GetAll")
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

	stats, err := c.service.GetAll(sctx, req.Cursor, req.Limit)
	if err != nil {
		c.logger.Error("Failed to get all stats",
			zap.Int("cursor", req.Cursor),
			zap.Int("limit", req.Limit),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(stats))
}

func (c *Controller) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.Get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var data GetRequest

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
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.Create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req CreateRequest

	err = ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	stat, err := c.service.Add(sctx, req)
	if err != nil {
		c.logger.Error("Failed to create stat", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, stat)
}

func (c *Controller) Modify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.Modify")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req UpdateRequest

	err = ctx.ShouldBindUri(&req.StatsId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("stat_id", req.Id))

	err = ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	modifiedStat, err := c.service.Update(sctx, req)
	if err != nil {
		c.logger.Error("Failed to update stat",
			zap.Int("stat_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, modifiedStat)
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.Delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req DeleteRequest

	err = ctx.ShouldBindUri(&req.StatsId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("stat_id", req.Id))

	_, err = c.service.Delete(sctx, req)
	if err != nil {
		c.logger.Error("Failed to delete stat",
			zap.Int("stat_it", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Replace(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.Replace")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req ReplaceRequest

	err = ctx.ShouldBindUri(&req.StatsId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("stat_id", req.Id))

	err = ctx.ShouldBindWith(&req.ReplacementData, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	replacedStat, err := c.service.Replace(sctx, req)
	if err != nil {
		c.logger.Error("Failed to replace stat",
			zap.Int("stat_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, replacedStat)
}

func (c *Controller) UpdatePlayerStats(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "stats.controller.UpdatePlayerStats")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var req UpdatePlayerStatRequest

	err = ctx.ShouldBindUri(&req.UserId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(
		attribute.String("user_id", req.Id),
		attribute.String("stat_id", req.StatId),
	)

	err = ctx.ShouldBindWith(&req.IncomingPlayerStat, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.UpdatePlayerStats(sctx, req)
	if err != nil {
		c.logger.Error("Failed to update player stat",
			zap.String("user_id", req.Id),
			zap.String("stat_it", req.StatId))
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

	var req DeletePlayerStatRequest

	err = ctx.ShouldBindUri(&req.UserId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("user_id", req.Id))

	err = c.service.DeletePlayerStats(sctx, req)
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

	var req GetPlayerStatRequest

	err = ctx.ShouldBindUri(&req.UserId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.String("user_id", req.Id))

	playerStats, err := c.service.GetPlayerStats(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get player stats", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, playerStats)
}
