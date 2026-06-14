package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Token() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := ExtractToken(c)
		if err != nil {
			c.Status(http.StatusUnauthorized)
		}

		claims, err := GetJwtService().ValidateAccessToken(c.Request.Context(), token)
		if err != nil {
			c.Status(http.StatusUnauthorized)
		}

		_ = claims // Dummy usage

		c.Next()
	}
}
