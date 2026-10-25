package common

import (
	"github.com/abhinash-kml/nova/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func ErrorHandler(c *gin.Context) {
	// Run the next handler first
	c.Next()

	if len(c.Errors) > 0 {
		err := c.Errors.Last().Err

		zap.L().Debug("", zap.Error(err))

		utils.SendProblemDetails(c, err)
		return
	}
}
