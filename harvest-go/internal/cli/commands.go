package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "harvest",
	Short: "Harvest CLI - Task tracking for Harvest",
	Long: `🌾 Harvest CLI - Simple command line interface for task tracking

A simple and efficient tool for tracking your daily tasks and generating reports for Harvest.`,
}

// Execute ejecuta el comando raíz
func Execute() error {
	return rootCmd.Execute()
}

// init inicializa los comandos
func init() {
	// Comando de ayuda personalizado
	rootCmd.SetHelpTemplate(`🌾 Harvest CLI

Usage:
  harvest [command] [args...]

Available Commands:
  add       Add a new task
  tech      Add a technical task
  meeting   Add a meeting task
  qa        Add a QA/testing task
  daily     Add daily standup
  status    Show today's status
  report    Generate report for Harvest
  upgrade   Upgrade to latest version

Examples:
  harvest add "Fix bug" 2.0
  harvest tech "Development" 3.5
  harvest meeting "Team sync" 1.0
  harvest status
  harvest report

Use "harvest [command] --help" for more information about a command.
`)

	// Comando de versión
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🌾 Harvest CLI v1.1.0")
		fmt.Println("Built with Go")
	},
}

// printError imprime errores de forma consistente
func printError(err error) {
	fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
}

// printSuccess imprime mensajes de éxito
func printSuccess(message string) {
	fmt.Printf("✅ %s\n", message)
}

// printInfo imprime información
func printInfo(message string) {
	fmt.Printf("ℹ️  %s\n", message)
}
