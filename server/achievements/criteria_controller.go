package achievements

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type CriteriaController struct {
	service CriteriaService
	logger  *zap.Logger
}

func NewCriteriaController(s CriteriaService, l *zap.Logger) *CriteriaController {
	return &CriteriaController{
		service: s,
		logger:  l,
	}
}

func (c *CriteriaController) GetAll(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.criteria.controller.getall")
	defer span.End()

	var dto GetAllCriteriaDTO

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
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
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, results)
}

func (c *CriteriaController) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.criteria.controller.get")
	defer span.End()

	var dto GetCriteriaDTO

	if err := ctx.ShouldBindUri(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	span.SetAttributes(attribute.Int("criteriaid", dto.CriteriaID))

	criteria, err := c.service.Get(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, criteria)
}

func (c *CriteriaController) GetByAchievement(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.criteria.controller.getbyachievement")
	defer span.End()

	var dto GetCriteriaByAchievementDTO

	if err := ctx.ShouldBindUri(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	span.SetAttributes(attribute.Int("achievementid", dto.Id))

	criteria, err := c.service.GetByAchievement(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, criteria)
}

func (c *CriteriaController) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievementcriteria.controller.create")
	defer span.End()

	var dto CreateCriteriaDTO

	if err := ctx.ShouldBindWith(&dto, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	criteria, err := c.service.Add(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusCreated, criteria)
}

func (c *CriteriaController) Modify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievementcriteria.controller.modify")
	defer span.End()

	var dto UpdateCriteriaDTO

	if err := ctx.ShouldBindUri(&dto.AchievementCriteriaId); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	span.SetAttributes(attribute.Int("criteria_id", dto.CriteriaID))

	if err := ctx.ShouldBindWith(&dto.UpdateCriteriaData, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	criteria, err := c.service.Update(sctx, dto)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, criteria)
}

func (c *CriteriaController) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievementcriteria.controller.delete")
	defer span.End()

	var dto DeleteCriteriaDTO

	if err := ctx.ShouldBindUri(&dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	span.SetAttributes(attribute.Int("criteriaid", dto.CriteriaID))

	if err := c.service.Delete(sctx, dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
