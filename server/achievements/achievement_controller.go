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

type Controller struct {
	config  *config.Config
	logger  *zap.Logger
	service Service
}

func NewController(service Service, c *config.Config, l *zap.Logger) *Controller {
	return &Controller{
		config:  c,
		service: service,
		logger:  l,
	}
}

func (c *Controller) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.Create")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to create achievement", zap.Error(err))
		}
		span.End()
	}()

	var req CreateAchievementRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.Create(sctx, req)
	if err != nil {
		c.logger.Error("Failed to create achievement", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusCreated, response)
}

func (c *Controller) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.Get")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to get achievement", zap.Error(err))
		}
		span.End()
	}()

	var req GetAchievementRequest

	err = ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	response, err := c.service.Get(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get achievement",
			zap.Int("achievement_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (c *Controller) Update(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.Update")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to update achievement", zap.Error(err))
		}
		span.End()
	}()

	var req UpdateAchievementRequest

	err = ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	achievement, err := c.service.Update(sctx, req)
	if err != nil {
		c.logger.Error("Failed to update achievement",
			zap.Int("achievement_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, achievement)
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.Delete")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to delete achievement", zap.Error(err))
		}
		span.End()
	}()

	var req DeleteAchievementRequest

	err = ctx.ShouldBindUri(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	err = c.service.Delete(sctx, req)
	if err != nil {
		c.logger.Error("Failed to delete achievement",
			zap.Int("achievement_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) CompleteAchievement(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.CompleteAchievement")
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

	err = c.service.CompleteAchievement(sctx, req)
	if err != nil {
		c.logger.Error("Failed to create achievement", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *Controller) GetProgress(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.GetProgress")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to get achievement progress", zap.Error(err))
		}
		span.End()
	}()

	var req GetProgressRequest

	err = ctx.ShouldBindUri(&req.ProgressId)
	if err != nil {
		ctx.Error(err)
		return
	}

	result, err := c.service.GetProgress(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get achievement progress",
			zap.Int("progress_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *Controller) DeleteProgress(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.controller.DeleteProgress")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to delete achievement progress", zap.Error(err))
		}
		span.End()
	}()

	var req DeleteProgressRequest

	err = ctx.ShouldBindUri(&req.ProgressId)
	if err != nil {
		ctx.Error(err)
		return
	}

	span.SetAttributes(attribute.Int("progress_id", req.Id))

	_, err = c.service.DeleteProgress(sctx, req)
	if err != nil {
		c.logger.Error("Failed to delete achievement progress",
			zap.Int("progress_id", req.Id),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) GetProgressOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.controller.GetProgressOfUser")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to get achievement progress of user", zap.Error(err))
		}
		span.End()
	}()

	var req GetProgressOfUserRequest

	err = ctx.ShouldBindQuery(&req)
	if err != nil {
		ctx.Error(err)
		return
	}

	result, err := c.service.GetProgressOfUser(sctx, req)
	if err != nil {
		c.logger.Error("Failed to get achievement progress of user", zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusOK, result)
}

func (c *Controller) UpdateProgressOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievement.controller.UpdateProgressOfUser")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to update achievement progress", zap.Error(err))
		}
		span.End()
	}()

	var req UpdateProgressRequest

	err = ctx.ShouldBindBodyWithJSON(&req)
	if err != nil {
		ctx.Error(err)
		return
	}
	span.SetAttributes(attribute.String("user_id", req.UserId))

	err = c.service.UpdateProgress(sctx, req)
	if err != nil {
		c.logger.Error("Failed to update achievement progress",
			zap.String("user_id", req.UserId),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) DeleteAllProgressOfUser(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "achievements.controller.DeleteAllProgressOfUser")
	var err error

	defer func() {
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			c.logger.Debug("Failed to delete achievement progress of user", zap.Error(err))
		}
		span.End()
	}()

	var req DeleteProgressOfUserRequest

	err = ctx.ShouldBindUri(&req.UserId)
	if err != nil {
		ctx.Error(err)
	}

	span.SetAttributes(attribute.String("user_id", req.UserId))

	err = c.service.DeleteAllProgressOfUser(sctx, req)
	if err != nil {
		c.logger.Error("Failed to delete acheievemnt progress of user",
			zap.String("user_id", req.UserId),
			zap.Error(err))
		ctx.Error(err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
