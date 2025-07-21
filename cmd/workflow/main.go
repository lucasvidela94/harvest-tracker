package main

import (
	"fmt"
	"os"

	"github.com/lucasvidela94/workflow-cli/internal/cli"
	"github.com/lucasvidela94/workflow-cli/internal/core"
)

// Version es la versión actual del CLI (usando la versión centralizada)
var Version = core.Version

func main() {
	if err := cli.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
