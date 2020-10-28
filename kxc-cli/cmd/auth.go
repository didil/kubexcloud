package main

import (
	"fmt"
	"log"

	"github.com/didil/kubexcloud/kxc-cli/client"
	"github.com/didil/kubexcloud/kxc-cli/config"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

func buildAuthCmd() *cobra.Command {
	var apiURL, userName, password string

	var authCmd = &cobra.Command{
		Use:   "auth",
		Short: "KubeXCloud Auth",
		Run: func(cmd *cobra.Command, args []string) {
			err := authRun(apiURL, userName, password)
			if err != nil {
				log.Fatalf("run: %v", err)
			}
		},
	}

	authCmd.Flags().StringVarP(&apiURL, "apiurl", "a", "", "apiurl")
	authCmd.Flags().StringVarP(&userName, "username", "u", "", "username")
	authCmd.Flags().StringVarP(&password, "password", "p", "", "password")

	return authCmd
}

func authRun(apiURL, userName, password string) error {
	// prompt for api endpoint if missing
	if apiURL == "" {
		prompt := promptui.Prompt{
			Label: "Api Url",
		}

		apiURLResult, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("api url prompt failed: %v", err)
		}

		apiURL = apiURLResult
	}

	// prompt for username if missing
	if userName == "" {
		prompt := promptui.Prompt{
			Label: "Username",
		}

		userNameResult, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("username prompt failed: %v", err)
		}

		userName = userNameResult
	}

	// prompt for password if missing
	if password == "" {
		prompt := promptui.Prompt{
			Label: "Password",
			Mask:  '*',
		}

		passwordResult, err := prompt.Run()
		if err != nil {
			return fmt.Errorf("password prompt failed: %v", err)
		}

		password = passwordResult
	}

	cl := client.NewClient()

	fmt.Printf("Authenticating ...\n")

	token, err := cl.LoginUser(apiURL, userName, password)
	if err != nil {
		return fmt.Errorf("auth: %v", err)
	}

	config.SetApiUrl(apiURL)
	config.SetAuthToken(token)

	err = config.WriteConfig()
	if err != nil {
		return fmt.Errorf("writeconfig: %v", err)
	}

	fmt.Printf("Authenticated Successfully\n")

	return nil
}
