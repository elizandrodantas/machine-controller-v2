package router

import (
	"time"

	"github.com/elizandrodantas/machine-controller-v2/api/middleware"
	"github.com/elizandrodantas/machine-controller-v2/internal/config"
	_ruleController "github.com/elizandrodantas/machine-controller-v2/rule/http/controller"
	_ruleRepository "github.com/elizandrodantas/machine-controller-v2/rule/repository"
	_ruleUsecase "github.com/elizandrodantas/machine-controller-v2/rule/usecase"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewRuleRouter(g *gin.RouterGroup, env *config.Env, pool *pgxpool.Pool, timeout time.Duration) {
	ruleRepo := _ruleRepository.NewRuleRepository(pool)
	ruleUS := _ruleUsecase.NewRuleUsecase(ruleRepo)

	controller := _ruleController.RuleController{
		RuleUsecase: ruleUS,
		Timeout:     timeout,
	}

	g.GET("/machine/rule", middleware.EnsureAutenticate(env), controller.List)
	g.POST("/machine/rule", middleware.EnsureAutenticate(env), controller.Create)
	g.GET("/machine/rule/history/:machine_id", middleware.EnsureAutenticate(env), controller.History)
	g.GET("/machine/rule/:machine_id", middleware.EnsureAutenticate(env), controller.Actives)
	g.DELETE("/machine/rule", middleware.EnsureAutenticate(env), controller.Remove)
}
