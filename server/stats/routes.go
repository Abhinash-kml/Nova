package stats

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, c *Controller) {
	// Private internal routes
	privateGroup := router.Group("/private/stats")
	{
		// Meta routes
		privateGroup.GET("/", c.GetAll)
		privateGroup.GET("/:id", c.Get)
		privateGroup.POST("/", c.Create)
		privateGroup.PATCH("/:id", c.Modify)
		privateGroup.PUT("/:id", c.Replace)
		privateGroup.DELETE("/:id", c.Delete)

		// Player specific routes
		privateGroup.POST("/player/:id", c.UpdatePlayerStats)
		privateGroup.DELETE("/player/:id", c.DeletePlayerStats)
	}

	// Public routes
	publicGroup := router.Group("/public/stats")
	{
		publicGroup.GET("/player/:id", c.GetPlayerStats)
	}
}
