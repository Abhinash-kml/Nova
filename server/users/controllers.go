package users

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

func (c *Controller) GetAll(ctx *gin.Context) {
	limit := ctx.DefaultQuery("limit", "20")
	cursor := ctx.DefaultQuery("cursor", "nil")

	_, _ = limit, cursor

}

func (c *Controller) Get(ctx *gin.Context) {
	id := ctx.Param("id")
	userid, err := uuid.FromBytes([]byte(id))
	if err != nil {
		// Log internally
		c.logger.Error("Failed to convert provided userid from string to uuid", zap.Error(err))

		ctx.Header("Content-Type", "application/problem+json")
		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "invalid field",
			Title:       "Invalid field",
			Description: "The provided field is invalid",
			Instance:    fmt.Sprintf("/users/%s", id),
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
			Type:        "not found",
			Title:       "Not found",
			Description: "The requested resource is not found",
			StatusCode:  404,
			Instance:    fmt.Sprintf("/users/%s", id),
		})
	}

	ctx.JSON(http.StatusOK, user)
}

func (c *Controller) Create(ctx *gin.Context) {
	var user UserCreateDTO

	err := ctx.BindJSON(&user)
	if err != nil {

	}

	err = c.service.Create(ctx.Request.Context(), user)
	if err != nil {

		ctx.JSON(http.StatusBadRequest, utils.ProblemDetails{
			Type:        "bad request",
			Title:       "Bad request",
			Description: "The provided resource is of bad format",
			StatusCode:  400,
			Instance:    "/users",
		})
	}

	ctx.JSON(http.StatusCreated, "OK")
}

func (c *Controller) Modify(ctx *gin.Context) {

}

func (c *Controller) Delete(ctx *gin.Context) {

}

func (c *Controller) Replace(ctx *gin.Context) {

}
