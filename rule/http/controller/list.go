package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
)

// List handles the HTTP GET request for rule list.
// @Summary List rules
// @Description List rules
// @Tags Rule
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(0)
// @Success 200 {object} domain.RuleListResponse
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /rule [get]
//
//	@Security JWT
func (r *RuleController) List(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), r.Timeout)
	defer cancel()

	page, err := strconv.Atoi(g.DefaultQuery("page", "0"))
	if err != nil || page < 0 {
		page = 0
	}

	list, err := r.RuleUsecase.List(ctx, page)
	if err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	count := len(list)
	total, err := r.RuleUsecase.Count(ctx)
	if err != nil {
		total = count
	}

	response := domain.RuleListResponse{
		Total: total,
		Count: count,
		Page:  page,
		Data:  list,
	}

	g.JSON(http.StatusOK, response)
}
