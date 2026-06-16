package posts

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

	posts, err := c.service.GetAll(ctx.Request.Context(), decodedCursor, dto.Limit)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, utils.Paginate(posts))
}

func (c *Controller) Get(ctx *gin.Context) {
	var dto GetDTO

	if err := ctx.ShouldBindUri(&dto); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	parsedUuid, _ := uuid.Parse(dto.Id)

	post, err := c.service.GetById(ctx.Request.Context(), parsedUuid)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.JSON(http.StatusOK, post)
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

	if err := ctx.ShouldBindUri(&dto.PostId); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindWith(&dto.FieldUpdates, binding.JSON); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.Update(ctx.Request.Context(), dto)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Delete(ctx *gin.Context) {
	var dto DeleteDTO

	if err := ctx.ShouldBindUri(&dto.PostId); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindQuery(&dto.DeleteOptions); err != nil {
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

func (c *Controller) Replace(ctx *gin.Context) {
	var dto ReplaceDTO

	if err := ctx.ShouldBindUri(&dto.PostId); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	if err := ctx.ShouldBindWith(&dto.ReplacementData, binding.JSON); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.Replace(ctx.Request.Context(), dto)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) BulkAdd(ctx *gin.Context) {
	var dto BulkCreateDTO

	if err := ctx.ShouldBindWith(&dto, binding.JSON); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.BulkAdd(ctx.Request.Context(), dto)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusCreated)
}

func (c *Controller) BulkModify(ctx *gin.Context) {
	var dto BulkModifyDTO

	if err := ctx.ShouldBindWith(&dto, binding.JSON); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.BulkModify(ctx.Request.Context(), dto)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) BulkDelete(ctx *gin.Context) {
	var dto BulkDeleteDTO

	if err := ctx.ShouldBindWith(&dto, binding.JSON); err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	err := c.service.BulkDelete(ctx.Request.Context(), dto)
	if err != nil {
		utils.SendProblemDetails(ctx, err)
		return
	}

	ctx.Status(http.StatusNoContent)
}
