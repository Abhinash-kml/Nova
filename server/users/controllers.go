package users

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type Controller struct {
	service Service
	logger  *zap.Logger
	// tracer otel.Tracer
}

func NewController(s Service) *Controller {
	return &Controller{
		service: s,
	}
}

func (c *Controller) GetAll(ctx *gin.Context) {}

func (c *Controller) Get(ctx *gin.Context) {

}

func (c *Controller) Create(ctx *gin.Context) {

}

func (c *Controller) Modify(ctx *gin.Context) {

}

func (c *Controller) Delete(ctx *gin.Context) {

}

func (c *Controller) Replace(ctx *gin.Context) {

}
