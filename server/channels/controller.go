package channels

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
}

func (c *Controller) GetAll(ctx *gin.Context) {
	var dto GetAllDTO

	if err := ctx.ShouldBindQuery(&dto); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	decodedCursor, err := utils.DecodeCursor(dto.Cursor)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	channels, err := c.service.GetAll(ctx.Request.Context(), decodedCursor, dto.Limit)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(channels))
}

func (c *Controller) Get(ctx *gin.Context) {
	var dto GetDTO

	if err := ctx.ShouldBindUri(&dto); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	parsedId, _ := uuid.Parse(dto.Id)

	channel, err := c.service.GetById(ctx.Request.Context(), parsedId)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, channel)
}

func (c *Controller) Create(ctx *gin.Context) {
	var dto CreateDTO

	if err := ctx.ShouldBindWith(&dto, binding.JSON); err != nil {
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
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindWith(&dto, binding.JSON); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.Modify(ctx.Request.Context(), dto)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Delete(ctx *gin.Context) {
	var dto DeleteDTO

	if err := ctx.ShouldBindUri(&dto); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.Delete(ctx.Request.Context(), dto)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
