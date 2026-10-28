package achievements

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, c *Controller) {
	achievement := router.Group("/private/achievements")
	{
		// Achievement routes
		//group.GET("", c.GetAll)            // Get list of  all achievements
		achievement.GET("/:id", c.Get)       // Get details of a particular achievement
		achievement.POST("", c.Create)       // Create an achievement
		achievement.PATCH("/:id", c.Update)  // Modify an existing achievement
		achievement.DELETE("/:id", c.Delete) // Delete an existing achievement
	}

	// Progress public routes
	publicProgress := router.Group("/public/achievements/progress")
	{
		publicProgress.GET("", c.GetProgressOfUser)
	}

	// Progress private routes
	privateProgress := router.Group("/private/achievements/progress")
	{
		privateProgress.GET("/:id", c.GetProgress)
		privateProgress.DELETE("/:id", c.DeleteProgress)
	}

	// Player progress operation private routes
	privatePlayerProgress := router.Group("/private/achievements/progress/user")
	{
		privatePlayerProgress.PATCH("", c.UpdateProgressOfUser)
		privatePlayerProgress.DELETE("/:id", c.DeleteAllProgressOfUser)
	}

	// Completed achievement private routes
	privateComplete := router.Group("/private/achievements/complete")
	{
		privateComplete.POST("", c.CompleteAchievement)
	}
}
