package router

import (
	"time"

	"github.com/elizandrodantas/machine-controller-v2/api/middleware"
	"github.com/elizandrodantas/machine-controller-v2/internal/config"
	"github.com/elizandrodantas/machine-controller-v2/internal/scope"
	_userControler "github.com/elizandrodantas/machine-controller-v2/user/http/controller"
	_userRepository "github.com/elizandrodantas/machine-controller-v2/user/repository"
	_userUsecase "github.com/elizandrodantas/machine-controller-v2/user/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewUserRouter(g *gin.RouterGroup, env *config.Env, pool *pgxpool.Pool, timeout time.Duration) {
	userRepo := _userRepository.NewUserRepository(pool)
	userUS := _userUsecase.NewUserUsecase(userRepo)

	controller := _userControler.UserController{
		UserUsecase: userUS,
		Timeout:     timeout,
	}

	g.POST("/user/register",
		middleware.EnsureAutenticate(env),
		middleware.EnsurePermission(scope.ScopeAdmin),
		controller.Register)
	g.GET("/user/info", middleware.EnsureAutenticate(env), controller.Info)
	g.PUT("/user/password", middleware.EnsureAutenticate(env), controller.ChangePassword)
	g.GET("/user",
		middleware.EnsureAutenticate(env),
		middleware.EnsurePermission(scope.ScopeAdmin),
		controller.List)
	g.POST("/user/scope",
		middleware.EnsureAutenticate(env),
		middleware.EnsurePermission(scope.ScopeAdmin),
		controller.UserScope)
	g.DELETE("/user/scope",
		middleware.EnsureAutenticate(env),
		middleware.EnsurePermission(scope.ScopeAdmin),
		controller.UserScope)
	g.PUT("/user/disable/:id",
		middleware.EnsureAutenticate(env),
		middleware.EnsurePermission(scope.ScopeAdmin),
		controller.Status)
	g.PUT("/user/enable/:id",
		middleware.EnsureAutenticate(env),
		middleware.EnsurePermission(scope.ScopeAdmin),
		controller.Status)
}
