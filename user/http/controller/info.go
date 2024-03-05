package controller

import (
	"context"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Info handles the HTTP request for user info.
// @Summary Get user info
// @Description Get user info
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} domain.UserJson
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /user/{id} [get]
//
//	@Security JWT
func (u *UserController) Info(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), u.Timeout)
	defer cancel()

	userId := g.GetString("user_id")
	if _, err := uuid.Parse(userId); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrParamIdNotUuid))
		return
	}

	info, err := u.UserUsecase.FindById(ctx, userId)
	if err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	g.JSON(http.StatusOK, domain.UserJson{
		ID:       info.ID,
		Name:     info.Name,
		Username: info.Username,
		Status:   info.Status,
		Scope:    info.Scope,
	})
}
