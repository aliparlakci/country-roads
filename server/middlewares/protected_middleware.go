package middlewares

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Protected(handler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		if _, exists := c.Get("user"); !exists {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		handler(c)
	}
}
