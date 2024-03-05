package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Cors is a middleware function that handles Cross-Origin Resource Sharing (CORS) headers.
func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Type, Authorization, X-Request-Id")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusAccepted)
			return
		}

		c.Next()
	}
}
