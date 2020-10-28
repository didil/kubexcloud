package main

import (
	"log"
	"os"

	api "github.com/didil/kubexcloud/kxc-api"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "bootstrap" {
		err := api.Bootstrap()
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	err := api.StartServer()
	if err != nil {
		log.Fatal(err)
	}
}
