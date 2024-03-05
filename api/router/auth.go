package router

import (
	"strconv"
	"time"

	_userController "github.com/elizandrodantas/machine-controller-v2/auth/http/controller"
	_userUsecase "github.com/elizandrodantas/machine-controller-v2/auth/usecase"
	"github.com/elizandrodantas/machine-controller-v2/internal/config"
	"github.com/elizandrodantas/machine-controller-v2/internal/security"
	_userRepository "github.com/elizandrodantas/machine-controller-v2/user/repository"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewAuthRouter(g *gin.RouterGroup, env *config.Env, pool *pgxpool.Pool, timeout time.Duration) {
	userRepo := _userRepository.NewUserRepository(pool)
	authUS := _userUsecase.NewAuthUsecase(userRepo)

	jsonwebtoken := security.JsonWebTokenSecurity{
		Key: []byte(env.JWT_TOKEN),
	}

	jwtexpireh, err := strconv.Atoi(env.JWT_EXPIRE_HOUR)
	if err != nil {
		jwtexpireh = 1
	} else if jwtexpireh <= 0 {
		jwtexpireh = 1
	}

	controller := _userController.AuthController{
		AuthUsecase:   authUS,
		Timeout:       timeout,
		JsonWebToken:  jsonwebtoken,
		JWTExpireHour: jwtexpireh,
		TokenType:     env.TOKEN_TYPE,
	}

	g.POST("/auth", controller.Auth)
}
