package controller

import (
	"context"
	"net/http"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ChangePassword handles the HTTP request for user change password.
// @Summary Change user password
// @Description Change user password
// @Tags User
// @Accept json
// @Produce json
// @Param data body domain.UserChangePasswordRequest true "UserChangePasswordRequest"
// @Success 200 {object} domain.UserChangePasswordResponse
// @Failure 400 {object} domain.HttpError
// @Failure 401 {object} domain.HttpError
// @Failure 500 {object} domain.HttpError
// @Router /user/change-password [put]
//
//	@Security JWT
func (u *UserController) ChangePassword(g *gin.Context) {
	ctx, cancel := context.WithTimeout(g.Request.Context(), u.Timeout)
	defer cancel()

	var data domain.UserChangePasswordRequest
	if err := g.ShouldBindJSON(&data); err != nil {
		g.JSON(http.StatusBadRequest, domain.ErrorHttpMessageFromError(domain.ErrInvalidParameters))
		return
	}

	userId := g.GetString("user_id")
	if _, err := uuid.Parse(userId); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	if err := u.UserUsecase.UpdatePassword(ctx, userId, data.NewPassword); err != nil {
		g.JSON(http.StatusInternalServerError, domain.ErrorHttpMessageFromError(domain.ErrInternalError))
		return
	}

	g.JSON(http.StatusOK, domain.UserChangePasswordResponse{
		Message: "Password changed successfully."})
}
