package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// RegisterRequest is a middleware function that generates a unique request ID and sets it as a header in the incoming request.
func RegisterRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		registerRequestId := uuid.NewString()

		c.Header("x-request-id", registerRequestId)

		c.Next()
	}
}
