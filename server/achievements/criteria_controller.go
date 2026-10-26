package achievements

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type CriteriaController struct {
	service CriteriaService
	config  *config.Config
	logger  *zap.Logger
}

func NewCriteriaController(s CriteriaService, c *config.Config, l *zap.Logger) *CriteriaController {
	return &CriteriaController{
		service: s,
		config:  c,
		logger:  l,
	}
}

func (c *CriteriaController) GetAll(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.criteria.controller.GetAll")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to get all achievement criteria", zap.Error(err))
		}
		span.End()
	}()

	var req GetAllCriteriaRequest

	err = ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	if req.Cursor < 0 {
		req.Cursor = 0
	}
	if req.Limit < 10 || req.Limit > 20 {
		req.Limit = 10
	}

	results, err := c.service.GetAll(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get all achievement criteria",
			zap.Int("cursor", req.Cursor),
			zap.Int("limit", req.Limit),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, results)
}

func (c *CriteriaController) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.criteria.controller.Get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to get achievement criteria", zap.Error(err))
		}
		span.End()
	}()

	var req GetCriteriaRequest

	err = ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("criteria_id", req.CriteriaID))

	criteria, err := c.service.Get(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get achievement criteria",
			zap.Int("criteria_id", req.CriteriaID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, criteria)
}

func (c *CriteriaController) GetByAchievement(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.criteria.controller.GetByAchievement")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to get achievement criteria by achievement", zap.Error(err))
		}
		span.End()
	}()

	var req GetCriteriaByAchievementRequest

	err = ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("achievement_id", req.Id))

	criteria, err := c.service.GetByAchievement(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get all criteria of an acheievement",
			zap.Int("achievement_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, criteria)
}

func (c *CriteriaController) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievementcriteria.controller.Create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to create achievement criteria", zap.Error(err))
		}
		span.End()
	}()

	var req CreateCriteriaRequest

	err = ctx.ShouldBindWith(&req, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	criteria, err := c.service.Add(sctx, req)
	if err != nil {
		c.logger.Error("Failed to create acheievement criteria",
			zap.Int("acheievement_id", req.AchievementId),
			zap.Int("stat_id", req.StatId),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, criteria)
}

func (c *CriteriaController) Modify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievementcriteria.controller.Modify")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to update achievement criteria", zap.Error(err))
		}
		span.End()
	}()

	var req UpdateCriteriaRequest

	err = ctx.ShouldBindUri(&req.AchievementCriteriaId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("criteria_id", req.CriteriaID))

	err = ctx.ShouldBindWith(&req.UpdateCriteriaData, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	criteria, err := c.service.Update(sctx, req)
	if err != nil {
		c.logger.Error("Failed to update achievement criteria",
			zap.Int("criteria_id", req.CriteriaID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, criteria)
}

func (c *CriteriaController) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievementcriteria.controller.Delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to delete achievement criteria", zap.Error(err))
		}
		span.End()
	}()

	var req DeleteCriteriaRequest

	err = ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("criteria_id", req.CriteriaID))

	err = c.service.Delete(sctx, req)
	if err != nil {
		c.logger.Error("Failed to delete achievement criteria",
			zap.Int("criteria_id", req.CriteriaID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
