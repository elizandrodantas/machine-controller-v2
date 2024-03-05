package controller

import (
	"context"
	"errors"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// Create handles the HTTP POST request for rule creation.
// @Summary Create a rule
// @Description Create a rule
// @Tags Rule
// @Accept json
// @Produce json
// @Param data body domain.RuleCreateRequest true "RuleCreateRequest"
// @Success 201 {object} domain.RuleCreateResponse
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /rule/create [post]
//
//	@Security JWT
func (r *RuleController) Create(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), r.Timeout)
	defer cancel()

	var data domain.RuleCreateRequest
	if err := g.ShouldBindJSON(&data); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
		return
	}

	if err := r.RuleUsecase.Create(ctx, data.MachineId, data.ServiceId, data.Expire); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			g.JSON(http.StatusBadRequest, domain.ErrorHttpMessage("machine or service not found"))
			return
		}

		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	g.JSON(http.StatusCreated, domain.RuleCreateResponse(data))
}
