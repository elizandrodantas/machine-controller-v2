package controller

import (
	"time"

	"github.com/elizandrodantas/machine-controller-v2/domain"
)

type UserController struct {
	UserUsecase domain.UserUsecase
	Timeout     time.Duration
}
