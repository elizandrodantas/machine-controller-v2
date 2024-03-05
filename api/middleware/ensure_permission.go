package middleware

import (
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/elizandrodantas/machine-controller-v2/internal/scope"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// EnsurePermission is a middleware function that checks if the user has the required permissions.
func EnsurePermission(r ...scope.Scope) gin.HandlerFunc {
	return func(c *gin.Context) {
		if len(r) <= 0 {
			c.Next()
			return
		}

		userId := c.GetString("user_id")
		if _, err := uuid.Parse(userId); err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorHttpMessageFromError(domain.ErrUnauthorized))
			return
		}

		scopes := c.GetStringSlice("scope")
		if len(scopes) <= 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorHttpMessageFromError(domain.ErrUnauthorized))
			return
		}

		var scopeFound bool
		for _, s := range scopes {
			parsed := scope.Parse(s)
			if parsed != -1 {
				if parsed.Has(r...) {
					scopeFound = true
					break
				}
			}
		}

		if !scopeFound {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorHttpMessageFromError(domain.ErrUnauthorized))
			return
		}

		c.Next()
	}
}
