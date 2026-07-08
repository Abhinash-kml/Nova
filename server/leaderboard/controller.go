package leaderboard

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	service Service
	logger  *zap.Logger
}

func NewController(s Service, l *zap.Logger) *Controller {
	return &Controller{
		service: s,
		logger:  l,
	}
}

func (c *Controller) GetAll(ctx *gin.Context) {

}

func (c *Controller) Get(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "leaderboard.controller.get")
	defer span.End()

}

func (c *Controller) Create(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "leaderboard.controller.create")
	defer span.End()
}

func (c *Controller) Modify(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "leaderboard.controller.modify")
	defer span.End()
}

func (c *Controller) Delete(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "leaderboard.controller.delete")
	defer span.End()
}

func (c *Controller) GetScore(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "leaderboard.controller.getscore")
	defer span.End()
}

func (c *Controller) UpdateScore(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "leaderboard.controller.updatescore")
	defer span.End()
}

func (c *Controller) DeleteScore(ctx *gin.Context) {
	sctx, span := tracer.Start(ctx.Request.Context(), "leaderboard.controller.deletescore")
	defer span.End()
}
