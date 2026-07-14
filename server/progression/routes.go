package progression

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, c *Controller) {
	group := router.Group("/progression")
	{
		// group.GET("/:id", c.GetByPlayerId)
		group.PATCH("/:id", c.Update)
	}
}
