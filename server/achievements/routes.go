package achievements

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, c *Controller, cc *CriteriaController, pc *ProgressController, cac *CompletedAchievementsController) {
	achievement := router.Group("/private/achievements")
	{
		// Achievement routes
		//group.GET("", c.GetAll)            // Get list of  all achievements
		achievement.GET("/:id", c.Get)       // Get details of a particular achievement
		achievement.POST("", c.Create)       // Create an achievement
		achievement.PATCH("/:id", c.Update)  // Modify an existing achievement
		achievement.DELETE("/:id", c.Delete) // Delete an existing achievement
	}

	// Criteria private routes
	criteria := router.Group("/private/achievements-criteria")
	{
		criteria.GET("", cc.GetAll)        // Get all achievementsd criterias
		criteria.GET("/:id", cc.Get)       // Get details of a particular achievement criteria
		criteria.POST("", cc.Create)       // Create an achievement criteria
		criteria.PATCH("/:id", cc.Modify)  // Modify an existing achievement criteria
		criteria.DELETE("/:id", cc.Delete) // Delete an existing achievement criteria
	}

	// Criteria public routes
	pCriteria := router.Group("/public/achievements-criteria")
	{
		pCriteria.GET("/:achievement_id", cc.GetByAchievement) // Get all criterias of achievement
	}

	// Progress public routes
	publicProgress := router.Group("/public/progress")
	{
		publicProgress.GET("", pc.GetOfUser)
	}

	// Progress private routes
	privateProgress := router.Group("/private/progress")
	{
		privateProgress.GET("/:id", pc.Get)
		privateProgress.POST("", pc.Create)
		privateProgress.PATCH("/:id", pc.Update)
		privateProgress.DELETE("/:id", pc.Delete)
	}

	// Completed achievement private routes
	privateCompleted := router.Group("/private/completed")
	{
		privateCompleted.POST("", cac.Create)
	}
}
