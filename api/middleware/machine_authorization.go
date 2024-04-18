package middleware

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/elizandrodantas/machine-controller-v2/internal/config"
	"github.com/elizandrodantas/machine-controller-v2/internal/protocol"
	"github.com/elizandrodantas/machine-controller-v2/internal/security"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// MachineAuthorization is a middleware function that performs authorization checks for machine requests.
// deprecated [v2.0.1]: MachineAuthorization is deprecated and should not be used.
func MachineAuthorization(env *config.Env) gin.HandlerFunc {
	return func(c *gin.Context) {
		// IF DEVELOPER MODE IS ENABLED, THE MACHINE DATA IS NOT ENCRYPTED
		if env.ENVIRONMENT == "development" {
			var machineData domain.MachineData
			if err := c.ShouldBindJSON(&machineData); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
				return
			}

			if err := validator.New().Struct(machineData); err != nil {
				c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
				return
			}

			c.Set("machineData", machineData)
			c.Next()
			return
		}

		var data domain.MachineRequest
		if err := c.ShouldBindJSON(&data); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
			return
		}

		// SEND DECODED DATA TO PROTOCOL TOOL
		protocolTool, err := protocol.ImportProtocolBase64(data.Data)
		if err != nil || !protocolTool.IsVersionValid() {
			c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
			return
		}

		// DECRYPT DATA
		crypto := security.Crypto(protocolTool.GetKey())
		decrypted, err := crypto.Decrypt(protocolTool.GetMessage(), protocolTool.GetIV())
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
			return
		}

		// UNMARSHAL DECRYPTED DATA
		var machineData domain.MachineData
		if err := json.Unmarshal(decrypted, &machineData); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
			return
		}

		// VALIDATE MACHINE DATA
		if err := validator.New().Struct(machineData); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
			return
		}

		// CHECK EXPIRATION TIME
		if machineData.Expire < int(time.Now().Unix()) {
			c.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
			return
		}

		// SET MACHINE DATA TO CONTEXT
		// AND ALLOW REQUEST TO PROCEED
		c.Set("machineData", machineData)
		c.Next()
	}
}
