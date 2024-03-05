package controller

import (
	"context"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
)

// Remove handles the HTTP DELETE request for rule removal.
// @Summary Remove a rule
// @Description Remove a rule
// @Tags Rule
// @Accept json
// @Produce json
// @Param data body domain.RuleRemoveRequest true "RuleRemoveRequest"
// @Success 204 "No Content"
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /rule/remove [delete]
//
//	@Security JWT
func (r *RuleController) Remove(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), r.Timeout)
	defer cancel()

	var data domain.RuleRemoveRequest
	if err := g.ShouldBindJSON(&data); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
		return
	}

	if _, err := r.RuleUsecase.ActiveByMachineIdAndServiceId(ctx, data.MachineId, data.ServiceId); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessage("Rule not found."))
		return
	}

	if err := r.RuleUsecase.Remove(ctx, data.MachineId, data.ServiceId); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	g.AbortWithStatus(http.StatusOK)
}
