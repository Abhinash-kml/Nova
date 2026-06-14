package auth

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func ExtractToken(c *gin.Context) (string, error) {
	token := c.GetHeader("authorization")
	parts := strings.Split(token, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("Invalid token format")
	}

	return parts[1], nil
}
