package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/elizandrodantas/machine-controller-v2/internal/security"
	"github.com/gin-gonic/gin"
)

// Auth is a method of the AuthController struct that handles the authentication logic.
// @Summary Authenticate a user
// @Description Authenticate a user
// @Tags Auth
// @Accept json
// @Produce json
// @Param data body domain.AuthRequest true "AuthRequest"
// @Success 200 {object} domain.AuthResponse
// @Failure 400 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /auth [post]
func (a *AuthController) Auth(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), a.Timeout)
	defer cancel()

	var data domain.AuthRequest
	if err := g.ShouldBindJSON(&data); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
		return
	}

	user, err := a.AuthUsecase.Autenticate(ctx, data.Username, data.Password)
	if err != nil {
		if err == domain.ErrUsernameOrPasswordInvalid || err == domain.ErrUserIsNotActive {
			g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(err))
			return
		}

		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrUsernameOrPasswordInvalid))
		return
	}

	jwtclaims := security.JsonWebTokenClaims{
		Sub:   user.ID,
		Scope: user.Scope,
	}

	token, err := a.JsonWebToken.Create(&jwtclaims, a.JWTExpireHour)
	if err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	response := domain.AuthResponse{
		Token:  token,
		Type:   a.TokenType,
		Expire: time.Now().UnixMilli() + (int64(a.JWTExpireHour) * time.Hour.Milliseconds()),
	}

	g.JSON(http.StatusOK, response)
}
