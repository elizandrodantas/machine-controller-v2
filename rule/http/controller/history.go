package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// History handles the HTTP request for rule history.
// @Summary Get rules history
// @Description Get rules history
// @Tags Rule
// @Accept json
// @Produce json
// @Param machine_id path string true "Machine ID"
// @Param service query string false "Filter by service ID"
// @Param page query int false "Page number" default(0)
// @Success 200 {object} domain.RuleHistoryResponse
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /rule/history/{machine_id} [get]
//
//	@Security JWT
func (r *RuleController) History(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), r.Timeout)
	defer cancel()

	machineId := g.Params.ByName("machine_id")
	if _, err := uuid.Parse(machineId); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrParamIdNotUuid))
		return
	}

	serviceId, ok := g.GetQuery("service")
	if ok {
		if _, err := uuid.Parse(serviceId); err != nil {
			g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrParamIdNotUuid))
			return
		}
	}

	page, err := strconv.Atoi(g.DefaultQuery("page", "0"))
	if err != nil || page < 0 {
		page = 0
	}

	list, err := r.RuleUsecase.History(ctx, machineId, serviceId, page)
	if err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	count := len(list)
	total, err := r.RuleUsecase.CountByMachineId(ctx, machineId, serviceId)
	if err != nil {
		total = count
	}

	response := domain.RuleHistoryResponse{
		MachineId: machineId,
		Total:     total,
		Count:     count,
		Page:      page,
		Data:      list,
	}

	g.JSON(http.StatusOK, response)
}
