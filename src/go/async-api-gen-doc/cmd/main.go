package main

import (
	"context"
	"log"

	asyncapigendoc "github.com/dnitsch/async-api-generator/cmd/async-api-gen-doc"
)

func main() {
	cmd := asyncapigendoc.NewCmd(context.Background())
	cmd.WithCommands()
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
