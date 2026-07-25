package achievements

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type CompletedAchievementsController struct {
	service CompletedAchievementService
	logger  *zap.Logger
}

func NewCompletedAchievementController(service CompletedAchievementService, l *zap.Logger) *CompletedAchievementsController {
	return &CompletedAchievementsController{
		service: service,
		logger:  l,
	}
}

func (c *CompletedAchievementsController) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.completed.controller.create")
	defer span.End()

	var dto CompleteAchievementDTO

	if err := ctx.ShouldBindBodyWith(&dto, binding.JSON); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
	}

	if err := c.service.Create(sctx, dto); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		utils.SendProblemDetails(ctx, err)
	}

	ctx.Status(http.StatusCreated)
}
