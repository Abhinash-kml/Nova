package clans

import "github.com/gin-gonic/gin"

type Controller struct {
	service Service
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
