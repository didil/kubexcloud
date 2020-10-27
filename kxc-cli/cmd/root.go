package main

import (
	"fmt"

	"github.com/didil/kubexcloud/kxc-cli/config"
	"github.com/spf13/cobra"
)

func Execute() error {
	err := config.InitConfig()
	if err != nil {
		return fmt.Errorf("initconfig: %v", err)
	}

	var rootCmd = &cobra.Command{
		Use:   "kxc",
		Short: "KubeXCloud CLI",
	}

	versionCmd := buildVersionCmd()
	rootCmd.AddCommand(versionCmd)

	authCmd := buildAuthCmd()
	rootCmd.AddCommand(authCmd)

	projectsCmd := buildProjectsCmd()
	rootCmd.AddCommand(projectsCmd)

	err = rootCmd.Execute()
	if err != nil {
		return fmt.Errorf("execute: %v", err)
	}

	return nil
}
