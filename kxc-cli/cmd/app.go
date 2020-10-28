package main

import (
	"fmt"
	"log"
	"os"

	"github.com/didil/kubexcloud/kxc-cli/client"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func buildAppsCmd() *cobra.Command {
	var appsCmd = &cobra.Command{
		Use:   "apps",
		Short: "KubeXCloud Apps",
	}

	var projectName string
	appsCmd.PersistentFlags().StringVarP(&projectName, "project", "p", "", "project")

	appsListCmd := buildAppsListCmd()
	appsCmd.AddCommand(appsListCmd)

	appRestartCmd := buildAppRestartCmd()
	appsCmd.AddCommand(appRestartCmd)

	return appsCmd
}

func buildAppsListCmd() *cobra.Command {
	var appsListCmd = &cobra.Command{
		Use:   "list",
		Short: "KubeXCloud Apps List",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName, err := cmd.Flags().GetString("project")
			if err != nil {
				return err
			}
			if projectName == "" {
				return fmt.Errorf("project name required")
			}

			err = listAppsRun(projectName)
			if err != nil {
				log.Fatalf("run: %v", err)
			}

			return nil
		},
	}

	return appsListCmd
}

func listAppsRun(projectName string) error {
	cl := client.NewClient()

	fmt.Printf("Fetching Apps for project %v ...\n", projectName)

	appsList, err := cl.ListApps(projectName)
	if err != nil {
		return fmt.Errorf("list apps: %v", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name"})

	for _, proj := range appsList.Apps {
		table.Append([]string{proj.Name})
	}
	table.Render()

	return nil
}

func buildAppRestartCmd() *cobra.Command {
	var appsListCmd = &cobra.Command{
		Use:   "restart <app>",
		Short: "KubeXCloud Apps List",
		RunE: func(cmd *cobra.Command, args []string) error {
			projectName, err := cmd.Flags().GetString("project")
			if err != nil {
				return err
			}
			if projectName == "" {
				return fmt.Errorf("project name required")
			}

			if len(args) == 0 {
				return fmt.Errorf("app name required")
			}

			appName := args[0]

			err = restartAppRun(projectName, appName)
			if err != nil {
				log.Fatalf("run: %v", err)
			}

			return nil
		},
	}

	return appsListCmd
}

func restartAppRun(projectName, appName string) error {
	cl := client.NewClient()

	fmt.Printf("Restarting App %s [Project %v] ...\n", appName, projectName)

	err := cl.RestartApp(projectName, appName)
	if err != nil {
		return fmt.Errorf("restart app: %v", err)
	}

	fmt.Printf("App restart triggered successfully")

	return nil
}
