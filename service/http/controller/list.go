package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
)

// List handles the HTTP GET request for service list.
// @Summary List services
// @Description List services
// @Tags Service
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(0)
// @Success 200 {object} domain.ServiceListResponse
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /service [get]
//
//	@Security JWT
func (s *ServiceController) List(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), s.Timeout)
	defer cancel()

	page, err := strconv.Atoi(g.DefaultQuery("page", "0"))
	if err != nil {
		page = 0
	}

	services, err := s.ServiceUsecase.List(ctx, page)
	if err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	count := len(services)
	total, err := s.ServiceUsecase.Count(ctx)
	if err != nil {
		total = count
	}

	g.JSON(http.StatusOK, domain.ServiceListResponse{
		Count: count,
		Total: total,
		Page:  page,
		Data:  services,
	})
}
