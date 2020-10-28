package auth

import (
	"fmt"

	"github.com/didil/kubexcloud/kxc-api/services"
)

// Login is a support function to fake login
func Login(username string) (string, error) {
	if username == "" {
		return "", fmt.Errorf("username empty")
	}

	token, err := services.SignJWT(username)
	if err != nil {
		return "", err
	}

	return token, nil
}
