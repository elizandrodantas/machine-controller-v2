package controller

import (
	"context"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
)

// Remove handles the HTTP DELETE request for service removal.
// @Summary Remove a service
// @Description Remove a service
// @Tags Service
// @Accept json
// @Produce json
// @Param data body domain.ServiceRemoveRequest true "ServiceRemoveRequest"
// @Success 204 "No Content"
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /service/{id} [delete]
//
//	@Security JWT
func (s *ServiceController) Remove(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), s.Timeout)
	defer cancel()

	var req domain.ServiceRemoveRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
		return
	}

	if _, err := s.ServiceUsecase.FindById(ctx, req.ID); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrServiceNotFound))
		return
	}

	if err := s.ServiceUsecase.Remove(ctx, req.ID); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	g.AbortWithStatus(http.StatusNoContent)
}
