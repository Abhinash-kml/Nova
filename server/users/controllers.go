package users

import (
	"net/http"

	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type Controller struct {
	service Service
	logger  *zap.Logger
	// tracer otel.Tracer
}

func NewController(s Service, l *zap.Logger) *Controller {
	return &Controller{
		service: s,
		logger:  l,
	}
}

func (c *Controller) GetAll(ctx *gin.Context) {
	var data GetAllDTO

	if err := ctx.ShouldBindQuery(&data); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	decodedCursor, err := utils.DecodeCursor(data.Cursor)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	users, err := c.service.GetAll(ctx.Request.Context(), decodedCursor, data.Limit)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(users))
}

func (c *Controller) Get(ctx *gin.Context) {
	var data GetDTO

	if err := ctx.ShouldBindUri(&data); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	parsedId, _ := uuid.Parse(data.Id)

	user, err := c.service.GetById(ctx.Request.Context(), parsedId)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) Create(ctx *gin.Context) {
	var dto CreateDTO

	if err := ctx.BindJSON(&dto); err != nil {
		c.logger.Error("Failed to bind UserCreateDTO", zap.Error(err))
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.Add(ctx.Request.Context(), dto)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *Controller) Modify(ctx *gin.Context) {
	var dto UpdateDTO

	if err := ctx.ShouldBindUri(&dto); err != nil {
		c.logger.Error("Failed to bind UserUpdateDTO", zap.Error(err))
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindBodyWithJSON(&dto); err != nil {
		c.logger.Error("Failed to bind UserUpdateDTO", zap.Error(err))
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.Update(ctx.Request.Context(), dto)
	if err != nil {
		c.logger.Error("Failed to update user", zap.Error(err))
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Delete(ctx *gin.Context) {
	var dto DeleteDTO

	if err := ctx.ShouldBindUri(&dto); err != nil {
		c.logger.Error("Failed to bind UserUpdateDTO", zap.Error(err))
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		c.logger.Error("Failed to bind UserUpdateDTO", zap.Error(err))
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.Delete(ctx.Request.Context(), dto)
	if err != nil {
		c.logger.Error("Failed to delete user", zap.Error(err))
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Replace(ctx *gin.Context) {
	var dto ReplaceDTO

	if err := ctx.ShouldBindUri(&dto); err != nil {
		c.logger.Error("Failed to bind UserUpdateDTO", zap.Error(err))
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindWith(&dto, binding.JSON); err != nil {
		c.logger.Error("Failed to bind UserReplaceDTO", zap.Error(err))
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.Replace(ctx.Request.Context(), dto)
	if err != nil {
		c.logger.Error("Failed to replace user", zap.Error(err))
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
