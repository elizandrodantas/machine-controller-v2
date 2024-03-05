package controller

import (
	"context"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

// Machine handles the HTTP POST request for machine registration.
// @Summary Machine registration
// @Description Machine registration
// @Tags Machine
// @Accept json
// @Produce json
// @Param data body domain.MachineData true "MachineData"
// @Success 200 {object} domain.MachineResponse
// @Failure 401 {object} domain.HttpError
// @Failure 422 {object} domain.MachineResponse
// @Failure 500 {object} domain.HttpError
// @Router /machine [post]
func (m *MachineController) Machine(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), m.Timeout)
	defer cancel()

	machineData, exist := g.Get("machineData")
	if !exist {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	data, ok := machineData.(domain.MachineData)
	if !ok {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	machineInfo, err := m.MachineUsecase.FindByGuid(g, data.Guid)
	if err != nil {
		if err == pgx.ErrNoRows {
			err := m.MachineUsecase.Create(ctx, &data)
			if err != nil {
				g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
				return
			}

			g.JSON(http.StatusUnprocessableEntity, domain.MachineResponse{
				Status:          domain.MACHINERESPONSE_CREATE,
				Message:         "Machine not found, register created",
				Identify:        data.Guid,
				Name:            data.Name,
				ServiceIdentify: data.ServiceId,
			})
			return
		}

		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	if _, err := m.RuleUsecase.ActiveByMachineIdAndServiceId(ctx, machineInfo.ID, data.ServiceId); err != nil {
		g.JSON(http.StatusUnprocessableEntity, domain.MachineResponse{
			Status:          domain.MACHINERESPONSE_UNAVAILABLE,
			Message:         "Machine not available",
			Identify:        data.Guid,
			Name:            data.Name,
			ServiceIdentify: data.ServiceId,
		})
		return
	}

	g.JSON(http.StatusOK, domain.MachineResponse{
		Status:          domain.MACHINERESPONSE_SUCCESS,
		Message:         "Welcome back",
		Identify:        data.Guid,
		ServiceIdentify: data.ServiceId,
		Name:            data.Name,
	})
}
