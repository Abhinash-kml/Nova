package users

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
	limit := ctx.DefaultQuery("limit", "20")
	limitNum, err := strconv.Atoi(limit)
	if err != nil || limitNum < 10 || limitNum > 20 {
		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/errors/validation-error",
			Title:       "Validation Error",
			Description: "The provided resource contains validation error",
			StatusCode:  400,
			Instance:    "GET /users",
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
			Instance:    "GET /users",
			Errors: []utils.ProblemDetailErrors{
				{
					Field:   "cursor",
					Message: "The provided cursor is invalid",
					Code:    "400",
				},
			},
		})
	}

	users := c.service.GetAll(ctx.Request.Context(), decodedCursor, limitNum)
	ctx.JSON(http.StatusOK, users)
}

func (c *Controller) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	userid, err := uuid.FromBytes([]byte(id))
	if err != nil {
		c.logger.Error("Failed to convert provided userid from string to uuid", zap.Error(err))

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/validation-error",
			Title:       "Validation Error",
			Description: "The provided field is invalid",
			Instance:    fmt.Sprintf("GET /users/%s", id),
			Errors: []utils.ProblemDetailErrors{
				{
					Field:   "id",
					Message: "invalid uuid",
					Code:    "100",
				},
			},
		})
	}

	user, found := c.service.GetById(ctx.Request.Context(), userid)
	if !found {
		ctx.JSON(http.StatusNotFound, utils.ProblemDetails{
			Type:        "nova.com/not-found",
			Title:       "Not found",
			Description: "The requested resource is not found",
			StatusCode:  404,
			Instance:    fmt.Sprintf("GET /users/%s", id),
		})
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) Create(ctx *gin.Context) {
	var dto UserCreateDTO

	err := ctx.BindJSON(&dto)
	if err != nil {
		c.logger.Error("Failed to bind UserCreateDTO", zap.Error(err))

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/errors/bad-request",
			Title:       "Bad request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "POST /users",
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
			Instance:    "POST /users",
		})
	}

	ctx.Status(http.StatusCreated)
}

func (c *Controller) Modify(ctx *gin.Context) {
	var dto UserUpdateDTO

	err := ctx.Bind(&dto)
	if err != nil {
		c.logger.Error("Failed to bind UserUpdateDTO", zap.Error(err))

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/errors/bad-request",
			Title:       "Bad request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "PATCH /users",
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
			Instance:    "PATCH /users",
		})
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Delete(ctx *gin.Context) {
	var dto UserDeleteDTO

	err := ctx.BindJSON(&dto)
	if err != nil {
		c.logger.Error("Failed to bind UserDeleteDTO", zap.Error(err))

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/errors/bad-request",
			Title:       "Bad request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "DELETE /users",
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
			Instance:    "DELETE /users",
		})
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Replace(ctx *gin.Context) {
	var dto UserReplaceDTO

	err := ctx.Bind(&dto)
	if err != nil {
		c.logger.Error("Failed to bind UserReplaceDTO", zap.Error(err))

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/errors/bad-request",
			Title:       "Bad request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "PUT /users",
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
			Instance:    "PUT /users",
		})
	}

	ctx.Status(http.StatusNoContent)
}
