package entitlements

import "github.com/gin-gonic/gin"

func SetupRoutes(r *gin.Engine, c *Controller) {
	private := r.Group("/private/entitlements")
	{
		private.POST("", c.Create)
		private.DELETE("/:id", c.Delete)
		private.GET("/:id", c.CheckIfActive)
		private.GET("", c.GetAllOfUser)
	}
}
