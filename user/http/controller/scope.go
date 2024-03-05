package controller

import (
	"context"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/elizandrodantas/machine-controller-v2/internal/scope"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// UserScope handles the HTTP request for managing user scopes.
// @Summary Manage user scopes
// @Description Manage user scopes
// @Tags User
// @Accept json
// @Produce json
// @Param data body domain.UserCremoveScopeRequest true "UserCremoveScopeRequest"
// @Success 204 "No Content"
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 409 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /user/{userId}/scope [delete]
//
//	@Security JWT
func (r *UserController) UserScope(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), r.Timeout)
	defer cancel()

	userId := g.GetString("user_id")
	if _, err := uuid.Parse(userId); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	var data domain.UserCremoveScopeRequest
	if err := g.ShouldBindJSON(&data); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
		return
	}

	parseScope := scope.Parse(data.Scope)
	if !parseScope.IsValid() {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrScopeInvalid))
		return
	}

	if g.Request.Method == http.MethodDelete {
		if userId == data.UserId && parseScope == scope.ScopeAdmin {
			g.JSON(http.StatusForbidden, domain.ErrorHttpMessage("user cannot remove admin scope from itself"))
			return
		}

		err := r.UserUsecase.RemoveScopeWithUserId(ctx, data.UserId, parseScope.String())
		if err == domain.ErrUserNotExist {
			g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(err))
			return
		}

		if err != nil {
			g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
			return
		}
	} else {
		err := r.UserUsecase.AddScopeWithUserId(ctx, data.UserId, parseScope.String())
		if err != nil {
			if err == domain.ErrUserNotExist || err == domain.ErrScopeAlreadyExist {
				g.JSON(http.StatusNotFound, domain.ErrorHttpMessageFromError(err))
				return
			}

			g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
			return
		}
	}

	g.AbortWithStatus(http.StatusNoContent)
}
