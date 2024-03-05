package controller

import (
	"context"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
)

// Create handles the HTTP POST request for service creation.
// @Summary Create a service
// @Description Create a service
// @Tags Service
// @Accept json
// @Produce json
// @Param data body domain.ServiceCreateRequest true "ServiceCreateRequest"
// @Success 201
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 409 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /service/create [post]
//
//	@Security JWT
func (s *ServiceController) Create(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), s.Timeout)
	defer cancel()

	var request domain.ServiceCreateRequest
	if err := g.ShouldBindJSON(&request); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
		return
	}

	if service, err := s.ServiceUsecase.FindByName(ctx, request.Name); err == nil && service.Name == request.Name {
		g.JSON(http.StatusConflict, domain.ErrorHttpMessageFromError(domain.ErrServiceAlreayExist))
		return
	}

	if err := s.ServiceUsecase.Create(ctx, request.Name, request.Description); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	g.AbortWithStatus(http.StatusCreated)
}
