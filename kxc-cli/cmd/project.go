package main

import (
	"fmt"
	"log"
	"os"

	"github.com/didil/kubexcloud/kxc-cli/client"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func buildProjectsCmd() *cobra.Command {
	var projectsCmd = &cobra.Command{
		Use:   "projects",
		Short: "KubeXCloud Projects",
	}

	projectsListCmd := buildProjectsListCmd()
	projectsCmd.AddCommand(projectsListCmd)

	return projectsCmd
}

func buildProjectsListCmd() *cobra.Command {
	var projectsListCmd = &cobra.Command{
		Use:   "list",
		Short: "KubeXCloud Projects List",
		Run: func(cmd *cobra.Command, args []string) {
			err := listProjectsRun()
			if err != nil {
				log.Fatalf("run: %v", err)
			}
		},
	}

	return projectsListCmd
}

func listProjectsRun() error {
	cl := client.NewClient()

	fmt.Printf("Fetching Projects ...\n")

	projectsList, err := cl.ListProjects()
	if err != nil {
		return fmt.Errorf("list projects: %v", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Name"})

	for _, proj := range projectsList.Projects {
		table.Append([]string{proj.Name})
	}
	table.Render()

	return nil
}
