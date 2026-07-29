package social

import "github.com/gin-gonic/gin"

func SetupRoutes(engine *gin.Engine, c *Controller) {
	friends := engine.Group("/public/social/friends")
	{
		friends.GET("", c.GetAllFriends)
		friends.POST("/add", c.AddFriend)
		friends.POST("/remove", c.RemoveFriend)
		friends.POST("/accept", c.AcceptFriendRequest)
		friends.POST("/reject", c.RejectFriendRequest)
		friends.POST("/block", c.BlockUser)
		friends.POST("/unblock", c.UnblockUser)
		friends.GET("/mutual", c.GetMutualFriends)
		friends.GET("/blocked", c.GetAllBlocked)
		friends.GET("/requests.incoming", c.GetAllIncomingFriendRequests)
		friends.GET("/requests.outgoing", c.GetAllOutgoingFriendRequests)
	}

	// messages := engine.Group("/public/social/messages")
	// {

	// }

}

// User relationships
// GET  social/friends?userid=
// POST social/friends/accept
// POST social/friends/reject
// POST social/friends/add
// POST social/friends/remove
// POST social/friends/block
// GET  social/friends/blocked?userid=
// GET  social/friends?userid=
// GET  social/friends/mutual
// GET  social/friends/incoming.requests?userid=
// GET  social/friends/outgoing.requests?userid=

// User messages
// GET social/conversations?userid=
// GET social/messages/conversation?begin=xx&end=xx
// POST social/messages/conversation?convid=xx
// DELETE social/messages/conversation?convid=xx&messageid=xx
