package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("homedir: %v", err)
	}

	configPath := filepath.Join(homeDir, ".kxc")

	err = os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		return fmt.Errorf("create config dir: %v", err)
	}

	viper.AddConfigPath(configPath)

	err = viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found, ignore

		} else {
			return fmt.Errorf("read config: %v", err)
		}
	}

	return nil
}

const apiURLKey = "apiURL"
const authTokenKey = "authToken"

func SetApiUrl(apiURL string) {
	viper.Set(apiURLKey, apiURL)
}

func SetAuthToken(authToken string) {
	viper.Set(authTokenKey, authToken)
}

func GetApiUrl() string {
	return viper.GetString(apiURLKey)
}

func GetAuthToken() string {
	return viper.GetString(authTokenKey)
}

func WriteConfig() error {
	err := viper.WriteConfig()
	if _, ok := err.(viper.ConfigFileNotFoundError); ok {
		// file doesn't exist
		err = viper.SafeWriteConfig()
	}
	if err != nil {
		return err
	}

	return nil
}
