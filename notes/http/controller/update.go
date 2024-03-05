package controller

import (
	"context"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Update handles the HTTP request for notes update.
// @Summary Update a note
// @Description Update a note
// @Tags Notes
// @Accept json
// @Produce json
// @Param data body domain.NotesUpdateRequest true "NotesUpdateRequest"
// @Success 200
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 403 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /notes/update [put]
//
//	@Security JWT
func (n *NotesController) Update(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), n.Timeout)
	defer cancel()

	userId := g.GetString("user_id")
	if _, err := uuid.Parse(userId); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	var req domain.NotesUpdateRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
		return
	}

	if _, err := uuid.Parse(req.ID); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessage("invalid note id"))
		return
	}

	if err := n.NotesUsecase.ChangeNote(ctx, req.ID, userId, req.Text); err != nil {
		if err == domain.ErrNoteNotFound || err == domain.ErrOnlyCreatorCanChange {
			g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
			return
		}

		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	g.AbortWithStatus(http.StatusOK)
}
