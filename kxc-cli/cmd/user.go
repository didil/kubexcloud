package main

import (
	"fmt"
	"log"
	"os"

	"github.com/didil/kubexcloud/kxc-cli/client"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func buildUsersCmd() *cobra.Command {
	var usersCmd = &cobra.Command{
		Use:   "users",
		Short: "KubeXCloud Users",
	}

	usersCreateCmd := buildUsersCreateCmd()
	usersCmd.AddCommand(usersCreateCmd)

	usersListCmd := buildUsersListCmd()
	usersCmd.AddCommand(usersListCmd)

	return usersCmd
}

func buildUsersCreateCmd() *cobra.Command {
	var userName, password, role string

	var usersCreateCmd = &cobra.Command{
		Use:   "create",
		Short: "KubeXCloud Users Create (admin only)",
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

func buildUsersListCmd() *cobra.Command {
	var usersListCmd = &cobra.Command{
		Use:   "list",
		Short: "KubeXCloud Users List (admin only)",
		Run: func(cmd *cobra.Command, args []string) {
			err := listUsersRun()
			if err != nil {
				log.Fatalf("run: %v", err)
			}
		},
	}

	return usersListCmd
}

func listUsersRun() error {
	cl := client.NewClient()

	fmt.Printf("Fetching Users ...\n")

	usersList, err := cl.ListUsers()
	if err != nil {
		return fmt.Errorf("list users: %v", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name", "Role"})

	for _, user := range usersList.Users {
		table.Append([]string{user.Name, user.Role})
	}
	table.Render()

	return nil
}
