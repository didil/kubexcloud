package main

import (
	"log"

	api "github.com/didil/kubexcloud/kxc-api"
)

func main() {
	err := api.StartServer()
	if err != nil {
		log.Fatal(err)
	}
}
