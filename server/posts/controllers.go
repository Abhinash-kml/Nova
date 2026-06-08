package posts

import (
	"fmt"
	"net/http"
	"strconv"

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

func (c *Controller) GetAll(ctx *gin.Context) {
	limit := ctx.DefaultQuery("limit", "10")
	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 10 || limitNum > 20 {
		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/validation-error",
			Title:       "Validation Error",
			Description: "The provided resource contains validation errors",
			StatusCode:  400,
			Instance:    "GET /posts",
			Errors: []utils.ProblemDetailErrors{
				{
					Field:   "limit",
					Message: "The provided limit is invalid. Valid range: 10 - 20",
					Code:    "400",
				},
			},
		})
	}

	cursor := ctx.DefaultQuery("cursor", "nil")
	decodedCursor, err := utils.DecodeCursor(cursor)
	if err != nil {
		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/errors/validation-error",
			Title:       "Bad request",
			Description: "The provided resource contains validation errors",
			StatusCode:  400,
			Instance:    "GET /posts",
			Errors: []utils.ProblemDetailErrors{
				{
					Field:   "cursor",
					Message: "The provided cursor is invalid",
					Code:    "400",
				},
			},
		})
	}

	posts := c.service.GetAll(ctx.Request.Context(), decodedCursor, limitNum)
	ctx.JSON(http.StatusOK, posts)
}

func (c *Controller) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	postid, err := uuid.FromBytes([]byte(id))
	if err != nil {
		c.logger.Error("Failed to convert provided postid from string to uuid", zap.Error(err))

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/validation-error",
			Title:       "Validation Error",
			Description: "The provided field is invalid",
			Instance:    fmt.Sprintf("GET /posts/%s", id),
			Errors: []utils.ProblemDetailErrors{
				{
					Field:   "id",
					Message: "invalid uuid",
					Code:    "100",
				},
			},
		})
	}

	post, found := c.service.GetById(ctx.Request.Context(), postid)
	if !found {
		ctx.JSON(http.StatusNotFound, utils.ProblemDetails{
			Type:        "nova.com/not-found",
			Title:       "Not found",
			Description: "The requested resource is not found",
			StatusCode:  404,
			Instance:    fmt.Sprintf("GET /posts/%s", id),
		})
	}

	ctx.JSON(http.StatusOK, post)
}

func (c *Controller) Create(ctx *gin.Context) {
	var dto PostCreateDTO

	err := ctx.BindJSON(&dto)
	if err != nil {
		c.logger.Error("Failed to bind PostCreateDTO", zap.Error(err))

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/errors/bad-request",
			Title:       "Bad request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "POST /posts",
		})
	}

	added := c.service.Add(ctx.Request.Context(), dto)
	if !added {
		ctx.Header("Content-Type", "application/problem+json")
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
	var dto PostUpdateDTO

	err := ctx.Bind(&dto)
	if err != nil {
		c.logger.Error("Failed to bind PostUpdateDTO", zap.Error(err))

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/errors/bad-request",
			Title:       "Bad request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "PATCH /posts",
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
			Instance:    "PATCH /posts",
		})
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Delete(ctx *gin.Context) {
	var dto PostDeleteDTO

	err := ctx.BindJSON(&dto)
	if err != nil {
		c.logger.Error("Failed to bind PostDeleteDTO", zap.Error(err))

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/errors/bad-request",
			Title:       "Bad request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "DELETE /posts",
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
			Instance:    "DELETE /posts",
		})
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Replace(ctx *gin.Context) {
	var dto PostReplaceDTO

	err := ctx.Bind(&dto)
	if err != nil {
		c.logger.Error("Failed to bind PostReplaceDTO", zap.Error(err))

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
