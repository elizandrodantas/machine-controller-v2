package controller

import (
	"context"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Detail handles the HTTP request for machine detail.
// @Summary Get machine detail
// @Description Get machine detail
// @Tags Machine
// @Accept json
// @Produce json
// @Param id path string true "Machine ID"
// @Success 200 {object} domain.Machine
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /machine/{id} [get]
//
//	@Security JWT
func (c *MachineController) Detail(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), c.Timeout)
	defer cancel()

	machineId := g.Params.ByName("id")
	if _, err := uuid.Parse(machineId); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrParamIdNotUuid))
		return
	}

	machine, err := c.MachineUsecase.Detail(ctx, machineId)
	if err != nil {
		if err == domain.ErrMachineNotFound {
			g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(err))
			return
		}

		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	g.JSON(http.StatusOK, machine)
}
