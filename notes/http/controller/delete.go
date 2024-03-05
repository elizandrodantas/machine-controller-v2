package controller

import (
	"context"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Delete handles the HTTP DELETE request for notes removal.
// @Summary Delete a note
// @Description Delete a note
// @Tags Notes
// @Accept json
// @Produce json
// @Param data body domain.NotesDeleteRequest true "NotesDeleteRequest"
// @Success 204 "No Content"
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 403 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /notes/delete [delete]
//
//	@Security JWT
func (n *NotesController) Delete(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), n.Timeout)
	defer cancel()

	var req domain.NotesDeleteRequest
	if err := g.ShouldBindJSON(&req); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
		return
	}

	userId := g.GetString("user_id")
	if _, err := uuid.Parse(userId); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	if _, err := uuid.Parse(req.ID); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessage("invalid note id"))
		return
	}

	if err := n.NotesUsecase.Delete(ctx, userId, req.ID); err != nil {
		if err == domain.ErrOnlyCreatorCanDelete {
			g.JSON(http.StatusForbidden, domain.ErrorHttpMessageFromError(err))
			return
		}

		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	g.AbortWithStatus(http.StatusNoContent)
}
