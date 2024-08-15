package main

import (
	"os"

	"github.com/rarimo/verificator-svc/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
