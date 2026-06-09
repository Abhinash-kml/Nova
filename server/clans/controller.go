package clans

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
	if err != nil {
		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/validation-error",
			Title:       "Validation Error",
			Description: "The provided resource contains validation errors",
			StatusCode:  400,
			Instance:    "GET /clans",
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
			Type:        "nova.com/validation-error",
			Title:       "Vaslidation Error",
			Description: "The provided resource contains validation errors",
			StatusCode:  400,
			Instance:    "GET /clans",
			Errors: []utils.ProblemDetailErrors{
				{
					Field:   "cursor",
					Message: "The provided cursor is invalid",
					Code:    "400",
				},
			},
		})
	}

	clans := c.service.GetAll(ctx.Request.Context(), decodedCursor, limitNum)
	ctx.JSON(http.StatusOK, clans)
}

func (c *Controller) Get(ctx *gin.Context) {
	id := ctx.Query("id")
	clanid, err := uuid.FromBytes([]byte(id))
	if err != nil {
		c.logger.Error("Failed to convert provided clanid from string to uuid", zap.Error(err))

		ctx.Header("Context-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/validation-error",
			Title:       "Validation Error",
			Description: "the provided field contains validation error",
			StatusCode:  400,
			Instance:    fmt.Sprintf("GET /clans/%s", id),
			Errors: []utils.ProblemDetailErrors{
				{
					Field:   "id",
					Message: "The provided id is invalid",
					Code:    "100",
				},
			},
		})
	}

	clan, found := c.service.GetById(ctx.Request.Context(), clanid)
	if !found {
		ctx.Header("Context-Type", "application/problem+json")
		ctx.JSON(http.StatusNotFound, utils.ProblemDetails{
			Type:        "nova.com/not-found",
			Title:       "Not Found",
			Description: "The requested resource is not found",
			StatusCode:  404,
			Instance:    fmt.Sprintf("GET /clans/%s", id),
		})
	}

	ctx.JSON(http.StatusOK, clan)
}

func (c *Controller) Create(ctx *gin.Context) {
	var dto CreateDTO

	err := ctx.Bind(&dto)
	if err != nil {
		c.logger.Error("Failed to bind ClanCreateDTO", zap.Error(err))

		ctx.Header("Context-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/bad-request",
			Title:       "Bad Request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "POST /clans",
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
			Instance:    "POST /clans",
		})
	}

	ctx.Status(http.StatusCreated)
}

func (c *Controller) Modify(ctx *gin.Context) {
	var dto UpdateDTO

	err := ctx.Bind(&dto)
	if err != nil {
		c.logger.Error("Failed to bind ClanUpateDTO", zap.Error(err))

		ctx.Header("Context-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/bad-request",
			Title:       "Bad Request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "PATCH /clans",
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
			Instance:    "PATCH /clans",
		})
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Delete(ctx *gin.Context) {
	var dto DeleteDTO

	err := ctx.BindJSON(&dto)
	if err != nil {
		c.logger.Error("Failed to bind ClanDeleteDTO", zap.Error(err))

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "nova.com/errors/bad-request",
			Title:       "Bad request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "DELETE /clans",
		})
	}

	deleted := c.service.Delete(ctx.Request.Context(), dto)
	if !deleted {
		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusInternalServerError, utils.ProblemDetails{
			Type:        "nova.com/errors/no-delete",
			Title:       "Not Deleted",
			Description: "The provided resource couldn't be deleted",
			StatusCode:  500,
			Instance:    "DELETE /clans",
		})
	}

	ctx.Status(http.StatusNoContent)
}

func (c *Controller) Replace(ctx *gin.Context) {
	ctx.Status(http.StatusNoContent)
}
