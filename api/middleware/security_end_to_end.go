package middleware

import (
	"net/http"
	"strings"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/elizandrodantas/machine-controller-v2/internal/config"
	"github.com/elizandrodantas/machine-controller-v2/internal/protocol"
	"github.com/elizandrodantas/machine-controller-v2/internal/security"
	"github.com/gin-gonic/gin"
)

// SecurityEndToEnd is a middleware function that performs request security interpretation.
// recivied data is encrypted and must be decrypted to be used.
func SecurityEndToEnd(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		// IF DEVELOPER MODE IS ENABLED, THE MACHINE DATA IS NOT ENCRYPTED
		if strings.Contains(env.ENVIRONMENT, "dev") {
			var data domain.MiddleSETERequest
			if err := c.ShouldBindJSON(&data); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
				return
			}

			c.Set("data", data.Data)
			c.Next()
			return
		}

		var data domain.MiddleSETERequest
		if err := c.ShouldBindJSON(&data); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
			return
		}

		protocolTool, err := protocol.ImportProtocolBase64(data.Data)
		if err != nil || !protocolTool.IsVersionValid() {
			c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
			return
		}

		crypto := security.Crypto(protocolTool.GetKey())
		decrypted, err := crypto.Decrypt(protocolTool.GetMessage(), protocolTool.GetIV())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
			return
		}

		c.Set("data", string(decrypted))
		c.Next()
	}
}
