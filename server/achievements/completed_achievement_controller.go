package achievements

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/config"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.opentelemetry.io/otel/codes"
	"go.uber.org/zap"
)

type CompletedAchievementsController struct {
	service CompletedAchievementService
	config  *config.Config
	logger  *zap.Logger
}

func NewCompletedAchievementController(service CompletedAchievementService, c *config.Config, l *zap.Logger) *CompletedAchievementsController {
	return &CompletedAchievementsController{
		service: service,
		config:  c,
		logger:  l,
	}
}

func (c *CompletedAchievementsController) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.completed.controller.Create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to complete achievement", zap.Error(err))
		}
		span.End()
	}()

	var req CompleteAchievementRequest

	err = ctx.ShouldBindBodyWith(&req, binding.JSON)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.Create(sctx, req)
	if err != nil {
		c.logger.Error("Failed to create achievement", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusCreated)
}
