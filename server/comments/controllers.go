package comments

import (
	"fmt"
	"net/http"

	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
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

func (c *Controller) GetAll(ctx *gin.Context) {}

func (c *Controller) Get(ctx *gin.Context) {
	id := ctx.Query("id")
	commentid, err := uuid.FromBytes([]byte(id))
	if err != nil {
		c.logger.Error("Failed to convert provided commentid from string to uuid", zap.Error(err))

		ctx.Header("Context-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/validation-error",
			Title:       "Validation Error",
			Description: "the provided field contains validation error",
			StatusCode:  400,
			Instance:    fmt.Sprintf("GET /comments/%s", id),
			Errors: []utils.ProblemDetailErrors{
				{
					Field:   "id",
					Message: "The provided id is invalid",
					Code:    "100",
				},
			},
		})
	}

	comment, found := c.service.GetById(ctx.Request.Context(), commentid)
	if !found {
		ctx.Header("Context-Type", "application/problem+json")
		ctx.JSON(http.StatusNotFound, utils.ProblemDetails{
			Type:        "nova.com/not-found",
			Title:       "Not Found",
			Description: "The requested resource is not found",
			StatusCode:  404,
			Instance:    "GET /comments",
		})
	}

	ctx.JSON(http.StatusOK, comment)
}

func (c *Controller) Create(ctx *gin.Context) {
	var dto CommentCreateDTO

	err := ctx.Bind(&dto)
	if err != nil {
		c.logger.Error("Failed to bind CommentCreateDTO", zap.Error(err))

		ctx.Header("Context-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/bad-request",
			Title:       "Bad Request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "POST /comments",
		})
	}

	added := c.service.Add(ctx.Request.Context(), dto)
	if !added {
		ctx.Header("Context-Type", "application/problem+json")
		ctx.JSON(http.StatusInternalServerError, utils.ProblemDetails{
			Type:        "nova.com/errors/record-no-create",
			Title:       "Record not created",
			Description: "The provided record could not be created",
			StatusCode:  500,
			Instance:    "POST /posts",
		})
	}

	ctx.Status(http.StatusCreated)
}

func (c *Controller) Modify(ctx *gin.Context) {
	var dto CommentUpdateDTO

	err := ctx.Bind(&dto)
	if err != nil {
		c.logger.Error("Failed to bind CommentUpateDTO", zap.Error(err))

		ctx.Header("Context-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/bad-request",
			Title:       "Bad Request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "PATCH /comments",
		})
	}

	updated := c.service.Update(ctx.Request.Context(), dto)
	if !updated {
		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/errors/not-updated",
			Title:       "Not updated",
			Description: "The provided resource couldn't be updated",
			StatusCode:  500,
			Instance:    "PATCH /comments",
		})
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Delete(ctx *gin.Context) {
	var dto CommentDeleteDTO

	err := ctx.BindJSON(&dto)
	if err != nil {
		c.logger.Error("Failed to bind CommentDeleteDTO", zap.Error(err))

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/errors/bad-request",
			Title:       "Bad request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "DELETE /comments",
		})
	}

	deleted := c.service.Delete(ctx.Request.Context(), dto.Id)
	if !deleted {
		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusInternalServerError, utils.ProblemDetails{
			Type:        "nova.com/errors/no-delete",
			Title:       "Not Deleted",
			Description: "The provided resource couldn't be deleted",
			StatusCode:  500,
			Instance:    "DELETE /comments",
		})
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Replace(ctx *gin.Context) {
	var dto CommentReplaceDTO

	err := ctx.Bind(&dto)
	if err != nil {
		c.logger.Error("Failed to bind CommentReplaceDTO", zap.Error(err))

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/errors/bad-request",
			Title:       "Bad request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "PUT /posts",
		})
	}

	replaced := c.service.Replace(ctx.Request.Context(), dto)
	if !replaced {
		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/errors/not-replaced",
			Title:       "Not replaced",
			Description: "The provided resource couldn't be replaced",
			StatusCode:  500,
			Instance:    "PUT /posts",
		})
	}

	ctx.Status(http.StatusNoContent)
}
