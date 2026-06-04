package clans

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, c *Controller) {
	group := router.Group("/comments")
	{
		group.GET("/", c.GetAll)
		group.GET("/:id", c.Get)
		group.POST("/:id", c.Create)
		group.PATCH("/:id", c.Modify)
		group.PUT("/:id", c.Replace)
		group.DELETE("/:id", c.Delete)
	}
}
