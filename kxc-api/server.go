package api

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/didil/kubexcloud/kxc-api/handlers"
	"github.com/didil/kubexcloud/kxc-api/services"
)

// StartServer starts the server
func StartServer() error {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	log.Printf("Initializing k8s service ...\n")
	k8sSvc, err := services.NewK8sService()
	if err != nil {
		return err
	}

	projectSvc := services.NewProjectService(k8sSvc)

	app := &handlers.App{
		ProjectSvc: projectSvc,
	}

	log.Printf("Initializing router ...\n")

	mux := BuildRouter(app)

	fmt.Printf("Listening on port %s\n", port)

	return http.ListenAndServe(fmt.Sprintf(":%s", port), mux)
}
