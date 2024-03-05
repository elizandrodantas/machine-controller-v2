package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"
	"os"

	"github.com/elizandrodantas/machine-controller-v2/domain"
	"github.com/elizandrodantas/machine-controller-v2/internal/scope"
	"github.com/elizandrodantas/machine-controller-v2/user/repository"
	"github.com/elizandrodantas/machine-controller-v2/user/usecase"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spf13/cobra"
)

func createUserUsecase(pool *pgxpool.Pool) domain.UserUsecase {
	userRepo := repository.NewUserRepository(pool)
	return usecase.NewUserUsecase(userRepo)
}

func random(length int, special bool) string {
	var charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if special {
		charset += "!@#$%^&*()_+"
	}
	var result string
	for i := 0; i < length; i++ {
		index, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return ""
		}
		result += string(charset[index.Int64()])
	}
	return result
}

var userCommand = &cobra.Command{
	Use:   "user",
	Short: "Manage users.",
}

var userCreateCommand = &cobra.Command{
	Use:   "create",
	Short: "Create a new user.",
	RunE: func(cmd *cobra.Command, args []string) error {
		pool, err := createPoolConnection()
		if err != nil {
			return fmt.Errorf("error connecting to database: %s", err)
		}
		defer pool.Close()

		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")
		name, _ := cmd.Flags().GetString("name")
		scope, _ := cmd.Flags().GetString("scope")
		force, _ := cmd.Flags().GetBool("force")

		if username == "" {
			username = fmt.Sprintf("user-%s", random(8, false))
		}

		if password == "" {
			password = random(16, true)
		}

		if name == "" {
			name = username
		}

		if scope == "" {
			scope = "user"
		}

		if len(password) < 6 {
			return fmt.Errorf("password must be at least 6 characters")
		}

		useCase := createUserUsecase(pool)

		if force {

			useCase.DeleteByUsername(context.Background(), username)
		}

		if err := useCase.Create(context.Background(), &domain.UserCreate{
			Name:     name,
			Username: username,
			Password: password,
			Scope:    []string{scope},
		}); err != nil {
			return fmt.Errorf("error creating user: %s", err)
		}

		fmt.Fprintf(os.Stdout, "User created successfully.\n")
		fmt.Fprintf(os.Stdout, "Name: %s\n", name)
		fmt.Fprintf(os.Stdout, "Username: %s\n", username)
		fmt.Fprintf(os.Stdout, "Password: %s\n", password)
		fmt.Fprintf(os.Stdout, "Scope: %s\n", scope)

		return nil
	},
}

var userListCommand = &cobra.Command{
	Use:   "list",
	Short: "List all users.",
	RunE: func(cmd *cobra.Command, args []string) error {
		pool, err := createPoolConnection()
		if err != nil {
			return fmt.Errorf("error connecting to database: %s", err)
		}
		defer pool.Close()

		useCase := createUserUsecase(pool)

		var listAll []domain.UserJson
		var page int

		for {
			user, err := useCase.List(context.Background(), page)
			if err != nil {
				return fmt.Errorf("error listing users: %s", err)
			}

			if len(user) == 0 {
				break
			}

			listAll = append(listAll, user...)

			page++
		}

		if len(listAll) == 0 {
			fmt.Printf("No users found.\n")
			return nil
		}

		for _, u := range listAll {
			fmt.Fprintf(os.Stdout, "ID: %s\n", u.ID)
			fmt.Fprintf(os.Stdout, "Name: %s\n", u.Name)
			fmt.Fprintf(os.Stdout, "Username: %s\n", u.Username)
			fmt.Fprintf(os.Stdout, "Scope: %s\n", u.Scope)
			fmt.Fprintf(os.Stdout, "Status: %t\n", u.Status)
			fmt.Fprintf(os.Stdout, "Created At: %s\n", u.CreatedAt)
			fmt.Fprintf(os.Stdout, "Updated At: %s\n", u.UpdatedAt)
			fmt.Fprintf(os.Stdout, "\n")
		}

		return nil
	},
}

var userScopeCommand = &cobra.Command{
	Use:   "scope",
	Short: "Manage user scopes.",
}

var userScopeAddCommand = &cobra.Command{
	Use:   "add",
	Short: "Add a new scope to a user.",
	RunE: func(cmd *cobra.Command, args []string) error {
		pool, err := createPoolConnection()
		if err != nil {
			return fmt.Errorf("error connecting to database: %s", err)
		}
		defer pool.Close()

		id, _ := cmd.Flags().GetString("id")
		username, _ := cmd.Flags().GetString("username")
		scop, _ := cmd.Flags().GetString("scope")

		if id == "" && username == "" {
			return fmt.Errorf("id or username is required")
		}

		if scope.Parse(scop) == -1 {
			return fmt.Errorf("scope is invalid")
		}

		useCase := createUserUsecase(pool)

		if id != "" {
			if err := useCase.AddScopeWithUserId(context.Background(), id, scop); err != nil {
				return fmt.Errorf("error adding scope to user: %s", err)
			}
		} else {
			user, err := useCase.FindByUsername(context.Background(), username)
			if err != nil {
				return fmt.Errorf("error finding username: %s", err)
			}

			if err := useCase.AddScopeWithUserId(context.Background(), user.ID, scop); err != nil {
				return fmt.Errorf("error adding scope to user: %s", err)
			}
		}

		fmt.Fprintf(os.Stdout, "Scope added successfully.\n")
		return nil
	},
}

var userScopeRemoveCommand = &cobra.Command{
	Use:   "remove",
	Short: "Remove a scope from a user.",
	RunE: func(cmd *cobra.Command, args []string) error {
		pool, err := createPoolConnection()
		if err != nil {
			return fmt.Errorf("error connecting to database: %s", err)
		}
		defer pool.Close()

		id, _ := cmd.Flags().GetString("id")
		username, _ := cmd.Flags().GetString("username")
		scop, _ := cmd.Flags().GetString("scope")

		if id == "" && username == "" {
			return fmt.Errorf("id or username is required")
		}

		if scope.Parse(scop) == -1 {
			return fmt.Errorf("scope is invalid")
		}

		useCase := createUserUsecase(pool)

		if id != "" {
			if err := useCase.RemoveScopeWithUserId(context.Background(), id, scop); err != nil {
				return fmt.Errorf("error adding scope to user: %s", err)
			}
		} else {
			user, err := useCase.FindByUsername(context.Background(), username)
			if err != nil {
				return fmt.Errorf("error finding username: %s", err)
			}

			if err := useCase.RemoveScopeWithUserId(context.Background(), user.ID, scop); err != nil {
				return fmt.Errorf("error adding scope to user: %s", err)
			}
		}

		fmt.Fprintf(os.Stdout, "Scope added successfully.\n")
		return nil
	},
}

func userCLI() *cobra.Command {
	userCreateCommand.Flags().StringP("name", "n", "", "User name.")
	userCreateCommand.Flags().StringP("username", "u", "", "Username.")
	userCreateCommand.Flags().StringP("password", "p", "", "Password.")
	userCreateCommand.Flags().StringP("scope", "s", "", "User scope.")
	userCreateCommand.Flags().BoolP("force", "f", false, "Force user creation.")
	userListCommand.Flags().BoolP("help", "h", false, "Show help.")
	userListCommand.Flags().BoolP("version", "v", false, "Show version.")

	userScopeAddCommand.Flags().String("id", "", "User ID.")
	userScopeAddCommand.Flags().StringP("username", "u", "", "Username.")
	userScopeAddCommand.Flags().StringP("scope", "s", "", "User scope. ( admin | user )")
	userScopeAddCommand.Flags().BoolP("help", "h", false, "Show help.")

	userScopeRemoveCommand.Flags().String("id", "", "User ID.")
	userScopeRemoveCommand.Flags().StringP("username", "u", "", "Username.")
	userScopeRemoveCommand.Flags().StringP("scope", "s", "", "User scope. ( admin | user )")
	userScopeRemoveCommand.Flags().BoolP("help", "h", false, "Show help.")

	userScopeCommand.AddCommand(userScopeAddCommand)
	userScopeCommand.AddCommand(userScopeRemoveCommand)
	userCommand.AddCommand(userCreateCommand)
	userCommand.AddCommand(userListCommand)
	userCommand.AddCommand(userScopeCommand)

	return userCommand
}
