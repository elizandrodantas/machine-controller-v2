package main

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/elizandrodantas/machine-controller-v2/api/middleware"
	"github.com/elizandrodantas/machine-controller-v2/api/router"
	"github.com/elizandrodantas/machine-controller-v2/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"

	_ "github.com/elizandrodantas/machine-controller-v2/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var apiCommand = &cobra.Command{
	Use:   "api",
	Short: "Start the API server.",
	RunE: func(cmd *cobra.Command, args []string) error {
		pool, err := createPoolConnection()
		if err != nil {
			return fmt.Errorf("error connecting to database: %s", err)
		}
		defer pool.Close()

		env, _ := config.GetEnv()

		port, _ := cmd.Flags().GetString("port")
		if port != "" {
			env.WEB_PORT = port
		}

		if !strings.HasPrefix(env.WEB_PORT, ":") {
			env.WEB_PORT = ":" + env.WEB_PORT
		}

		APIResolve(&env, pool)
		return nil
	},
}

func apiCLI() *cobra.Command {
	apiCommand.Flags().StringP("port", "p", "3000", "Port to run the API server.")

	return apiCommand
}

func APIResolve(env *config.Env, pool *pgxpool.Pool) {
	timeout_int, err := strconv.Atoi(env.TIMEOUT)
	if err != nil {
		timeout_int = 60
	} else if timeout_int <= 0 {
		timeout_int = 60
	}

	engine := gin.Default()

	engine.Use(middleware.Cors())
	engine.Use(middleware.RegisterRequest())

	timeout := time.Duration(timeout_int) * time.Second

	if env.ENVIRONMENT == "development" {
		engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	}

	baseUrl := engine.Group("/api/v2")

	router.NewUserRouter(baseUrl, env, pool, timeout)
	router.NewAuthRouter(baseUrl, env, pool, timeout)
	router.NewMachineRouter(baseUrl, env, pool, timeout)
	router.NewRuleRouter(baseUrl, env, pool, timeout)
	router.NewServiceRouter(baseUrl, env, pool, timeout)
	router.NewNotesRouter(baseUrl, env, pool, timeout)

	engine.Run(env.WEB_PORT)
}
