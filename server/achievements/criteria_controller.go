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
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.criteria.controller.getall")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto GetAllCriteriaDTO

	err = ctx.ShouldBindQuery(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	if dto.Cursor < 0 {
		dto.Cursor = 0
	}
	if dto.Limit < 10 || dto.Limit > 20 {
		dto.Limit = 10
	}

	results, err := c.service.GetAll(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get all achievement criteria",
			zap.Int("cursor", dto.Cursor),
			zap.Int("limit", dto.Limit),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, results)
}

func (c *CriteriaController) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.criteria.controller.get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto GetCriteriaDTO

	err = ctx.ShouldBindUri(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("criteria_id", dto.CriteriaID))

	criteria, err := c.service.Get(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get achievement criteria",
			zap.Int("criteria_id", dto.CriteriaID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, criteria)
}

func (c *CriteriaController) GetByAchievement(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.criteria.controller.getbyachievement")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto GetCriteriaByAchievementDTO

	err = ctx.ShouldBindUri(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("achievement_id", dto.Id))

	criteria, err := c.service.GetByAchievement(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to get all criteria of an acheievement",
			zap.Int("achievement_id", dto.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, criteria)
}

func (c *CriteriaController) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievementcriteria.controller.create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto CreateCriteriaDTO

	err = ctx.ShouldBindWith(&dto, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	criteria, err := c.service.Add(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to create acheievement criteria",
			zap.Int("acheievement_id", dto.AchievementId),
			zap.Int("stat_id", dto.StatId),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, criteria)
}

func (c *CriteriaController) Modify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievementcriteria.controller.modify")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto UpdateCriteriaDTO

	err = ctx.ShouldBindUri(&dto.AchievementCriteriaId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("criteria_id", dto.CriteriaID))

	err = ctx.ShouldBindWith(&dto.UpdateCriteriaData, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	criteria, err := c.service.Update(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to update achievement criteria",
			zap.Int("criteria_id", dto.CriteriaID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, criteria)
}

func (c *CriteriaController) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievementcriteria.controller.delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		}
		span.End()
	}()

	var dto DeleteCriteriaDTO

	err = ctx.ShouldBindUri(&dto)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("criteria_id", dto.CriteriaID))

	err = c.service.Delete(sctx, dto)
	if err != nil {
		c.logger.Error("Failed to delete achievement criteria",
			zap.Int("criteria_id", dto.CriteriaID),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
