package router

import (
	"time"

	"github.com/elizandrodantas/machine-controller-v2/api/middleware"
	"github.com/elizandrodantas/machine-controller-v2/internal/config"
	"github.com/elizandrodantas/machine-controller-v2/internal/scope"
	_serviceController "github.com/elizandrodantas/machine-controller-v2/service/http/controller"
	_serviceRepository "github.com/elizandrodantas/machine-controller-v2/service/repository"
	_serviceUsecase "github.com/elizandrodantas/machine-controller-v2/service/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewServiceRouter(g *gin.RouterGroup, env *config.Env, pool *pgxpool.Pool, timeout time.Duration) {
	serviceRepo := _serviceRepository.NewServiceRepository(pool)
	serviceUS := _serviceUsecase.NewServiceUsecase(serviceRepo)

	controller := _serviceController.ServiceController{
		ServiceUsecase: serviceUS,
		Timeout:        timeout,
	}

	g.POST("/service",
		middleware.EnsureAutenticate(env),
		middleware.EnsurePermission(scope.ScopeAdmin),
		controller.Create)
	g.GET("/service", middleware.EnsureAutenticate(env), controller.List)
	g.DELETE("/service",
		middleware.EnsureAutenticate(env),
		middleware.EnsurePermission(scope.ScopeAdmin),
		controller.Remove)
}
