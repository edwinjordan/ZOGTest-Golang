package main

import (
	"context"
	"fmt"
	"os"

	"github.com/edwinjordan/ZOGTest-Golang.git/cmd/commands"
	"github.com/edwinjordan/ZOGTest-Golang.git/config"
	"github.com/edwinjordan/ZOGTest-Golang.git/internal/logging"
)

func init() {
	config.LoadEnv()
}

func main() {
	if len(os.Args) < 2 {
		logging.LogErrorMessage(context.Background(), "Expected a command")
		os.Exit(1)
	}

	command := os.Args[1]
	args := os.Args[2:]

	err := commands.Execute(command, args)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
