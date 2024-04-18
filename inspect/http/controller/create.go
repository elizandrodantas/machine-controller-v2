package controller

import (
	"fmt"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
)

func (c *InspectController) Create(g *gin.Context) {
	data, exist := g.Get("data")

	if !exist || data == "" {
		g.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
		return
	}

	fmt.Println("Inspect: ", data)
	g.AbortWithStatus(http.StatusOK)
}
