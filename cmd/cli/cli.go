package main

import (
	"fmt"
	"os"

	"github.com/elizandrodantas/machine-controller-v2/database"
	"github.com/elizandrodantas/machine-controller-v2/internal/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

var mainCommand = &cobra.Command{
	Use:   "machine-cli",
	Short: "Machine CLI",
	Long:  `Machine CLI is a command line interface for managing the Machine API.11`,
}

func createPoolConnection() (*pgxpool.Pool, error) {
	env, err := config.GetEnv()
	if err != nil {
		return nil, fmt.Errorf("error getting environment: %s", err)
	}

	pool, err := database.NewPostgresConnect(&env)
	if err != nil {
		return nil, fmt.Errorf("error connecting to database: %s", err)
	}

	return pool, nil
}

func main() {
	for _, k := range []*cobra.Command{apiCLI(), userCLI()} {
		mainCommand.AddCommand(k)
	}

	if err := mainCommand.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	// EXIT SUCESS
	os.Exit(0)
}
