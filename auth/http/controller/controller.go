package controller

import (
	"time"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/elizandrodantas/machine-controller-v2/internal/security"
)

type AuthController struct {
	AuthUsecase   domain.AuthUsecase
	Timeout       time.Duration
	JsonWebToken  security.JsonWebTokenSecurity
	JWTExpireHour int
	TokenType     string
}
