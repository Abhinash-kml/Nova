package posts

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, c *Controller) {
	group := router.Group("/posts")
	{
		group.GET("/", c.GetAll)
		group.GET("/:id", c.Get)
		group.POST("/", c.Create)
		group.PATCH("/", c.Modify)
		group.PUT("/", c.Replace)
		group.DELETE("/", c.Delete)
	}
}
