package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// List handles the HTTP request for notes list.
// @Summary List notes
// @Description List notes
// @Tags Notes
// @Accept json
// @Produce json
// @Param id path string true "Machine ID"
// @Param page query int false "Page"
// @Success 200 {object} domain.NotesListResponse
// @Failure 401 {object} domain.HttpError
// @Failure 400 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /notes/{id} [get]
//
//	@Security JWT
func (n *NotesController) List(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), n.Timeout)
	defer cancel()

	page, err := strconv.Atoi(g.DefaultQuery("page", "0"))
	if err != nil || page < 0 {
		page = 0
	}

	machineId, exist := g.Params.Get("id")
	if _, err := uuid.Parse(machineId); !exist || err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessage("invalid machine id"))
		return
	}

	list, err := n.NotesUsecase.ListByMachineId(ctx, machineId, page)
	if err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	count := len(list)
	total, err := n.NotesUsecase.CountByMachineId(ctx, machineId)
	if err != nil {
		total = count
	}

	g.JSON(http.StatusOK, &domain.NotesListResponse{
		Total: total,
		Count: count,
		Page:  page,
		Data:  list,
	})
}
