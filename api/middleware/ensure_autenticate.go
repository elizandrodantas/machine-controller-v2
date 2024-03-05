package middleware

import (
	"net/http"
	"strings"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/elizandrodantas/machine-controller-v2/internal/config"
	"github.com/elizandrodantas/machine-controller-v2/internal/security"
	"github.com/gin-gonic/gin"
)

// EnsureAutenticate is a middleware function that ensures the authentication of the request.
func EnsureAutenticate(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")
		if authorization == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorHttpMessageFromError(domain.ErrUnauthorized))
			return
		}

		parts := strings.Split(authorization, " ")
		if len(parts) != 2 || parts[0] != env.TOKEN_TYPE {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorHttpMessageFromError(domain.ErrUnauthorized))
			return
		}

		jsonwebtoken := security.JsonWebTokenSecurity{
			Key: []byte(env.JWT_TOKEN),
		}

		user, err := jsonwebtoken.Verify(parts[1])
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, domain.ErrorHttpMessageFromError(domain.ErrUnauthorized))
			return
		}

		c.Set("user_id", user.Subject)
		c.Set("scope", user.Scope)
		c.Next()
	}
}
