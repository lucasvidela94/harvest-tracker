package main

import (
	"fmt"
	"os"

	"github.com/lucasvidela94/harvest-cli/internal/cli"
)

// Version es la versión actual del CLI
var Version = "2.0.1"

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
