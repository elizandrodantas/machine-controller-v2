package controller

import (
	"context"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
)

// Register handles the HTTP POST request for user registration.
// @Summary Register a user
// @Description Register a user
// @Tags User
// @Accept json
// @Produce json
// @Param data body domain.UserRegisterRequest true "UserRegisterRequest"
// @Success 201
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /user/register [post]
// @Security JWT
func (u UserController) Register(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), u.Timeout)
	defer cancel()

	var data domain.UserRegisterRequest
	if err := g.ShouldBindJSON(&data); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
		return
	}

	if err := u.UserUsecase.Create(ctx, &domain.UserCreate{
		Name:     data.Name,
		Username: data.Username,
		Password: data.Password,
	}); err != nil {
		if err == domain.ErrUserAlreayExist {
			g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(err))
			return
		}

		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	g.AbortWithStatus(http.StatusCreated)
}
