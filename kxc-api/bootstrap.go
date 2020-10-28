package api

import (
	"context"
	"fmt"
	"log"

	"github.com/didil/kubexcloud/kxc-api/lib"
	"github.com/didil/kubexcloud/kxc-api/requests"
	"github.com/didil/kubexcloud/kxc-api/services"
	"github.com/sethvargo/go-password/password"
)

// Bootstrap bootstraps the server
func Bootstrap() error {
	// try load env if .env file found
	err := lib.LoadEnv(".env")
	if err != nil {
		// skip file not found errors to allow .env file to be optional
		if err.Error() != fmt.Sprintf("open .env: no such file or directory") {
			return fmt.Errorf("load env err: %v", err)
		}
	}

	log.Printf("Initializing k8s service ...\n")
	k8sSvc, err := services.NewK8sService()
	if err != nil {
		return fmt.Errorf("init k8s service: %v", err)
	}

	userSvc := services.NewUserService(k8sSvc)

	userName := "admin"

	pwd, err := password.Generate(10, 2, 0, false, false)
	if err != nil {
		return fmt.Errorf("password generate: %v", err)
	}

	reqData := &requests.CreateUser{
		Name:     userName,
		Password: pwd,
		Role:     services.UserRoleAdmin,
	}

	err = userSvc.Create(context.Background(), reqData)
	if err != nil {
		return fmt.Errorf("user create: %v", err)
	}

	log.Printf("KXC Api Server bootstraped successfully\n")
	log.Printf("Admin User credentials:\n")
	log.Printf("username: %s\n", userName)
	log.Printf("password: %s\n", pwd)

	return nil
}
