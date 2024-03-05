package controller

import (
	"context"
	"net/http"
	"strings"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// Status handles the HTTP request for user status.
// @Summary Change user status
// @Description Change user status
// @Tags User
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /user/{id}/enable [put]
//
//	@Security JWT
func (u *UserController) Status(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), u.Timeout)
	defer cancel()

	userId := g.GetString("user_id")
	if _, err := uuid.Parse(userId); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	userIdParam, exist := g.Params.Get("id")
	if _, err := uuid.Parse(userIdParam); !exist || err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrParamIdNotUuid))
		return
	}

	if userId == userIdParam {
		g.JSON(http.StatusForbidden, domain.ErrorHttpMessage("user cannot change own status"))
		return
	}

	if i := strings.Index(g.Request.URL.Path, "enable"); i > 0 {
		err := u.UserUsecase.Enable(ctx, userIdParam)
		if err != nil {
			if err == domain.ErrUserNotExist {
				g.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(err))
				return
			}

			g.AbortWithStatusJSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
			return
		}

		g.AbortWithStatus(http.StatusNoContent)
		return
	} else if i := strings.Index(g.Request.URL.Path, "disable"); i > 0 {
		err := u.UserUsecase.Disable(ctx, userIdParam)
		if err != nil {
			if err == domain.ErrUserNotExist {
				g.AbortWithStatusJSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(err))
				return
			}

			g.AbortWithStatusJSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
			return
		}

		g.AbortWithStatus(http.StatusNoContent)
		return
	}

	g.AbortWithStatus(http.StatusNotFound)
}
