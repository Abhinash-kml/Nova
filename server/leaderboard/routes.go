package leaderboard

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, c *Controller) {
	group := router.Group("/leaderboards")
	{
		// General routes
		group.GET("", c.GetAll)        // Get list of  all leaderboards
		group.GET("/:id", c.Get)       // Get details of a particular leaderboard
		group.POST("", c.Create)       // Create a leaderboard
		group.PUT("/:id", c.Modify)    // Modify an existing leaderboard
		group.DELETE("/:id", c.Delete) // Delete an existing leaderboard

		// Score routes
		group.GET("/:id/score", c.GetScore)       // Get scores in leaderboard
		group.POST("/:id/score", c.UpdateScore)   // Add score in leaderboard
		group.PUT("/:id/score", c.UpdateScore)    // Update/Modify score in leaderboard
		group.DELETE("/:id/score", c.DeleteScore) // Delete score in leaderboard
	}
}
