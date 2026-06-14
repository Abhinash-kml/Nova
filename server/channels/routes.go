package channels

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, c *Controller) {
	group := router.Group("/comments")
	{
		group.GET("/", c.GetAll)
		group.GET("/:id", c.Get)
		group.POST("/", c.Create)
		group.PATCH("/:id", c.Modify)
		group.DELETE("/:id", c.Delete)

		group.POST("/bulk", c.BulkAdd)
		group.PATCH("/bulk", c.BulkModify)
		group.DELETE("/bulk", c.BulkDelete)
	}
}
