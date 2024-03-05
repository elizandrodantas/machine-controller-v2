package controller

import (
	"context"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
)

// List handles the HTTP GET request for machine list.
// @Summary List machines
// @Description List machines
// @Tags Machine
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(0)
// @Param machine_id query string false "Filter by machine ID"
// @Param os query string false "Filter by OS"
// @Param q query string false "Filter with query"
// @Success 200 {object} domain.MachineListResponse
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /machine [get]
//
//	@Security JWT
func (m *MachineController) List(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), m.Timeout)
	defer cancel()

	var req domain.MachineListQuerys
	if err := g.ShouldBindQuery(&req); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInvalidQuery))
		return
	}

	list, err := m.MachineUsecase.List(ctx, req)
	if err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	count := len(list)
	total, err := m.MachineUsecase.Count(ctx)
	if err != nil {
		total = count
	}

	g.JSON(http.StatusOK, domain.MachineListResponse{
		Total: total,
		Count: count,
		Page:  req.Page,
		Data:  list,
	})
}
