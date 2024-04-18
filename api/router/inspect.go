package router

import (
	"time"

	"github.com/elizandrodantas/machine-controller-v2/api/middleware"
	_inspectController "github.com/elizandrodantas/machine-controller-v2/inspect/http/controller"
	"github.com/elizandrodantas/machine-controller-v2/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewInspectRouter(g *gin.RouterGroup, env *config.Env, pool *pgxpool.Pool, timeout time.Duration) {
	controller := _inspectController.InspectController{
		Timeout: timeout,
	}

	g.POST("/inspect", middleware.SecurityEndToEnd(env), controller.Create)
}
