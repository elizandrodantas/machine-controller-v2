package controller

import (
	"context"
	"net/http"
	"strconv"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
)

// List handles the HTTP GET request for user list.
// @Summary List users
// @Description List users
// @Tags User
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(0)
// @Success 200 {object} domain.UserListResponse
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /user [get]
//
//	@Security JWT
func (u *UserController) List(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), u.Timeout)
	defer cancel()

	page, err := strconv.Atoi(g.DefaultQuery("page", "0"))
	if err != nil || page < 0 {
		page = 0
	}

	users, err := u.UserUsecase.List(ctx, page)
	if err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	count := len(users)
	total, err := u.UserUsecase.Count(ctx)
	if err != nil {
		total = count
	}

	g.JSON(http.StatusOK, domain.UserListResponse{
		Total: total,
		Count: count,
		Page:  page,
		Data:  users,
	})
}
