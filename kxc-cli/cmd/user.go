package main

import (
	"fmt"
	"log"

	"github.com/didil/kubexcloud/kxc-cli/client"
	"github.com/spf13/cobra"
)

func buildUsersCmd() *cobra.Command {
	var usersCmd = &cobra.Command{
		Use:   "users",
		Short: "KubeXCloud Users",
	}

	usersCreateCmd := buildUsersCreateCmd()
	usersCmd.AddCommand(usersCreateCmd)

	return usersCmd
}

func buildUsersCreateCmd() *cobra.Command {
	var userName, password, role string

	var usersCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "KubeXCloud Users Create",
		RunE: func(cmd *cobra.Command, args []string) error {
			if userName == "" {
				return fmt.Errorf("username is empty")
			}
			if password == "" {
				return fmt.Errorf("password is empty")
			}
			if role == "" {
				return fmt.Errorf("role is empty")
			}
			err := createUsersRun(userName, password, role)
			if err != nil {
				log.Fatalf("run: %v", err)
			}

			return nil
		},
	}

	usersCreateCmd.Flags().StringVarP(&userName, "username", "u", "", "username")
	usersCreateCmd.Flags().StringVarP(&password, "password", "p", "", "password")
	usersCreateCmd.Flags().StringVarP(&role, "role", "r", "regular", "role")

	return usersCreateCmd
}

func createUsersRun(userName, password, role string) error {
	cl := client.NewClient()

	fmt.Printf("Creating User %s [role: %s]...\n", userName, role)

	err := cl.CreateUser(userName, password, role)
	if err != nil {
		return fmt.Errorf("create user: %v", err)
	}

	fmt.Printf("User created successfully\n")

	return nil
}
