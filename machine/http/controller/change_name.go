package controller

import (
	"context"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
)

// UpdateName handles the HTTP PUT request for machine name update.
// @Summary Update machine name
// @Description Update machine name
// @Tags Machine
// @Accept json
// @Produce json
// @Param data body domain.MachineUpdateNameRequest true "MachineUpdateNameRequest"
// @Success 200
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /machine/update-name [put]
//
//	@Security JWT
func (m *MachineController) UpdateName(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), m.Timeout)
	defer cancel()

	var req domain.MachineUpdateNameRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
		return
	}

	if err := m.MachineUsecase.UpdateName(ctx, req.Id, req.Name); err != nil {
		if err == domain.ErrMachineNotFound {
			g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(err))
			return
		}

		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	g.AbortWithStatus(http.StatusOK)
}
