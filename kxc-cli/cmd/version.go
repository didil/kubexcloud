package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	BuildVersion string = ""
	BuildTime    string = ""
)

func buildVersionCmd() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "KubeXCloud CLI Version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("KubeXCloud CLI\nBuild Version: %v\nBuild Time: %v\n", BuildVersion, BuildTime)
		},
	}

	return versionCmd
}
