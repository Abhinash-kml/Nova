package users

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, c *Controller) {
	group := router.Group("/users")
	{
		// General routes
		group.GET("", c.GetAll)
		group.GET("/:id", c.Get)
		group.POST("/", c.Create)
		group.PATCH("/:id", c.Modify)
		group.PUT("/:id", c.Replace)
		group.DELETE("/:id", c.Delete)

		// Auth route
		group.GET("/login", c.Login)
		group.GET("/refresh", c.Refresh)

		// Bulk operations
		group.POST("/bulk", c.BulkAdd)
		group.PATCH("/bulk", c.BulkModify)
		group.DELETE("/bulk", c.BulkDelete)
	}
}
