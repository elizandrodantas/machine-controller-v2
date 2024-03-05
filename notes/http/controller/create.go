package controller

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

// Create handles the HTTP POST request for notes creation.
// @Summary Create a note
// @Description Create a note
// @Tags Notes
// @Accept json
// @Produce json
// @Param data body domain.NotesCreateRequest true "NotesCreateRequest"
// @Success 204
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /notes/create [post]
//
//	@Security JWT
func (n *NotesController) Create(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), n.Timeout)
	defer cancel()

	userId := g.GetString("user_id")
	if _, err := uuid.Parse(userId); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	var req domain.NotesCreateRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
		return
	}

	if err := n.NotesUsecase.Create(ctx, req.Text, req.MachineId, userId); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr); pgerrcode.IsIntegrityConstraintViolation(pgErr.Code) {
			g.JSON(http.StatusBadRequest, domain.ErrorHttpMessage(fmt.Sprintf("machine with id '%s' not found", req.MachineId)))
			return
		}

		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	g.AbortWithStatus(http.StatusNoContent)
}
