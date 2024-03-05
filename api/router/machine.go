package router

import (
	"time"

	"github.com/elizandrodantas/machine-controller-v2/api/middleware"
	"github.com/elizandrodantas/machine-controller-v2/internal/config"
	_machineController "github.com/elizandrodantas/machine-controller-v2/machine/http/controller"
	_machineRepository "github.com/elizandrodantas/machine-controller-v2/machine/repository"
	_machineUsecase "github.com/elizandrodantas/machine-controller-v2/machine/usecase"
	_ruleRepository "github.com/elizandrodantas/machine-controller-v2/rule/repository"
	_ruleUsecase "github.com/elizandrodantas/machine-controller-v2/rule/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewMachineRouter(g *gin.RouterGroup, env *config.Env, pool *pgxpool.Pool, timeout time.Duration) {
	machineRepo := _machineRepository.NewMachineRepository(pool)
	machineUC := _machineUsecase.NewMachineUsecase(machineRepo)

	ruleRepo := _ruleRepository.NewRuleRepository(pool)
	ruleUC := _ruleUsecase.NewRuleUsecase(ruleRepo)

	controller := _machineController.MachineController{
		MachineUsecase: machineUC,
		RuleUsecase:    ruleUC,
		Timeout:        timeout,
	}

	g.POST("/machine", middleware.MachineAuthorization(env), controller.Machine)
	g.PUT("/machine", middleware.EnsureAutenticate(env), controller.UpdateName)
	g.GET("/machine", middleware.EnsureAutenticate(env), controller.List)
	g.GET("/machine/:id", middleware.EnsureAutenticate(env), controller.Detail)
}
