package controller

import (
	"context"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Actives handles the HTTP request for rule actives.
// @Summary Get actives rules
// @Description Get actives rules
// @Tags Rule
// @Accept json
// @Produce json
// @Param machine_id path string true "Machine ID"
// @Success 200 {object} domain.RuleActivesResponse
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /rule/actives/{machine_id} [get]
//
//	@Security JWT
func (r *RuleController) Actives(g *gin.Context) {
	_, cancel := context.WithTimeout(g.Request.Context(), r.Timeout)
	defer cancel()

	machineId := g.Params.ByName("machine_id")
	if _, err := uuid.Parse(machineId); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrParamIdNotUuid))
		return
	}

	list, err := r.RuleUsecase.Actives(g.Request.Context(), machineId)
	if err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	response := domain.RuleActivesResponse{
		MachineId: machineId,
		Total:     len(list),
		Data:      list,
	}

	g.JSON(http.StatusOK, response)
}
