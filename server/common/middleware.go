package common

import (
	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
)

func ErrorHandler(c *gin.Context) {
	// Run the next handler first
	c.Next()

	if len(c.Errors) > 0 {
		err := c.Errors.Last().Err
		utils.SendProblemDetails(c, err)
		return
	}
}
