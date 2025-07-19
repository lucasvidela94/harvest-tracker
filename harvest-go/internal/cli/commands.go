package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lucasvidela94/harvest-cli/internal/core"
	"github.com/lucasvidela94/harvest-cli/pkg/harvest"
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

	// Comando add
	rootCmd.AddCommand(addCmd)

	// Comando status
	rootCmd.AddCommand(statusCmd)
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

// addCmd es el comando para agregar tareas
var addCmd = &cobra.Command{
	Use:   "add [description] [hours] [category]",
	Short: "Add a new task",
	Long: `Add a new task to your daily tracking.

Examples:
  harvest add "Fix bug" 2.0
  harvest add "Development" 3.5 tech
  harvest add "Meeting" 1.0 meeting`,
	Args: cobra.MinimumNArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		description := args[0]
		hoursStr := args[1]

		// Validar horas
		var hours float64
		var err error

		// Intentar parsear como número
		if hours, err = parseHours(hoursStr); err != nil {
			printError(fmt.Errorf("invalid hours: %s", hoursStr))
			return
		}

		// Categoría por defecto
		category := "general"
		if len(args) > 2 {
			category = args[2]
		}

		// Agregar tarea
		taskManager := core.NewTaskManager()
		if err := taskManager.AddTask(description, hours, category); err != nil {
			printError(err)
			return
		}

		printSuccess(fmt.Sprintf("Added task: %s (%.1fh %s)", description, hours, category))

		// Mostrar estado actual
		showStatus(taskManager)
	},
}

// parseHours convierte una cadena de horas a float64
func parseHours(hoursStr string) (float64, error) {
	// Intentar parsear como float
	var hours float64
	var err error

	// Reemplazar comas por puntos para compatibilidad
	hoursStr = strings.ReplaceAll(hoursStr, ",", ".")

	if hours, err = strconv.ParseFloat(hoursStr, 64); err != nil {
		return 0, fmt.Errorf("could not parse hours: %s", hoursStr)
	}

	if hours <= 0 {
		return 0, fmt.Errorf("hours must be greater than 0")
	}

	return hours, nil
}

// showStatus muestra el estado actual
func showStatus(taskManager *core.TaskManager) {
	todayTasks, err := taskManager.GetTodayTasks()
	if err != nil {
		printError(err)
		return
	}

	totalHours := taskManager.GetTotalHours(todayTasks)
	targetHours := taskManager.GetDailyHoursTarget()

	fmt.Printf("\n📊 Today's Status: %.1fh / %.1fh (%.1fh remaining)\n",
		totalHours, targetHours, targetHours-totalHours)

	if len(todayTasks) > 0 {
		fmt.Println("\n📝 Today's tasks:")
		for _, task := range todayTasks {
			icon := harvest.GetIcon(task.Category)
			fmt.Printf("  %s %s - %.1fh\n", icon, task.Description, task.Hours)
		}
	}
}

// statusCmd es el comando para mostrar el estado
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show today's status",
	Long: `Show today's task status and progress.

Displays:
- Current date
- Hours worked vs target
- Remaining hours
- List of today's tasks with icons`,
	Run: func(cmd *cobra.Command, args []string) {
		taskManager := core.NewTaskManager()
		showDetailedStatus(taskManager)
	},
}

// showDetailedStatus muestra el estado detallado
func showDetailedStatus(taskManager *core.TaskManager) {
	todayTasks, err := taskManager.GetTodayTasks()
	if err != nil {
		printError(err)
		return
	}

	totalHours := taskManager.GetTotalHours(todayTasks)
	targetHours := taskManager.GetDailyHoursTarget()
	remainingHours := targetHours - totalHours

	// Mostrar fecha
	fmt.Printf("📅 Today (%s): %.2fh / %.1fh\n",
		time.Now().Format("2006-01-02"), totalHours, targetHours)

	// Mostrar horas restantes
	if remainingHours > 0 {
		fmt.Printf("📈 Remaining: %.2fh\n", remainingHours)
	} else if remainingHours < 0 {
		fmt.Printf("📈 Overtime: %.2fh\n", -remainingHours)
	} else {
		fmt.Printf("📈 Perfect! Target reached\n")
	}

	// Mostrar tareas
	if len(todayTasks) > 0 {
		for _, task := range todayTasks {
			icon := harvest.GetIcon(task.Category)
			fmt.Printf("  %s %s (%.1fh)\n", icon, task.Description, task.Hours)
		}
	} else {
		fmt.Println("  No tasks for today")
	}

	// Mostrar barra de progreso
	percentage := (totalHours / targetHours) * 100
	if percentage > 100 {
		percentage = 100
	}

	barLength := 20
	filledLength := int((percentage / 100) * float64(barLength))

	fmt.Printf("📊 [")
	for i := 0; i < barLength; i++ {
		if i < filledLength {
			fmt.Print("█")
		} else {
			fmt.Print("░")
		}
	}
	fmt.Printf("] %.1f%%\n", percentage)
}
