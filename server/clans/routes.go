package clans

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine, c *Controller) {
	group := router.Group("/clans")
	{
		group.GET("/", c.GetAll)
		group.GET("/:id", c.Get)
		group.POST("/", c.Create)
		group.PATCH("/:id", c.Modify)
		group.DELETE("/:id", c.Delete)
	}

	memGroup := router.Group("/clans/members")
	{
		memGroup.GET("/:id", c.GetAllMembers)
	}

	actGroup := router.Group("/clans/action")
	{
		actGroup.POST("/send-invite", c.InviteMember)
		actGroup.POST("/accept-invite", c.AcceptInvite)
		actGroup.POST("/reject-invite", c.RejectInvite)
		actGroup.POST("/list-invites", c.ListAllOutogingInvites)
		actGroup.POST("/send-join-request", c.SendJoinRequest)
		actGroup.POST("/accept-join-request", c.AcceptJoinRequest)
		actGroup.POST("/reject-join-request", c.RejectJoinRequest)
		actGroup.POST("/list-join-requests", c.ListAllJoinRequests)
	}
}
