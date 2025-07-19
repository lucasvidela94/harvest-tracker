package cli

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/lucasvidela94/harvest-cli/internal/core"
	"github.com/lucasvidela94/harvest-cli/internal/upgrade"
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

	// Comandos específicos
	rootCmd.AddCommand(techCmd)
	rootCmd.AddCommand(meetingCmd)
	rootCmd.AddCommand(qaCmd)
	rootCmd.AddCommand(dailyCmd)

	// Comando report
	rootCmd.AddCommand(reportCmd)

	// Comando upgrade
	rootCmd.AddCommand(upgradeCmd)

	// Agregar comando rollback
	rootCmd.AddCommand(rollbackCmd)
}

// rollbackCmd es el comando para gestionar rollbacks
var rollbackCmd = &cobra.Command{
	Use:   "rollback",
	Short: "Manage rollback operations",
	Long: `Manage rollback operations for Harvest CLI.

This command allows you to:
- Check rollback availability
- View rollback information
- Perform manual rollback
- View rollback logs`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		showRollbackInfo()
	},
}

// showRollbackInfo muestra información sobre el rollback
func showRollbackInfo() {
	rollbackManager := upgrade.NewRollbackManager()

	fmt.Println("🔄 Rollback Information")
	fmt.Println(strings.Repeat("─", 50))

	// Verificar disponibilidad
	if rollbackManager.IsRollbackAvailable() {
		printSuccess("Rollback is available")

		// Obtener información detallada
		info, err := rollbackManager.GetRollbackInfo()
		if err != nil {
			printError(fmt.Errorf("could not get rollback info: %v", err))
			return
		}

		// Mostrar información del backup del binario
		if size, ok := info["binary_backup_size"]; ok {
			fmt.Printf("Binary backup size: %d bytes\n", size)
		}

		if backupTime, ok := info["binary_backup_time"]; ok {
			if t, ok := backupTime.(time.Time); ok {
				fmt.Printf("Binary backup time: %s\n", t.Format("2006-01-02 15:04:05"))
			}
		}

		// Mostrar información del backup de datos
		if path, ok := info["data_backup_path"]; ok {
			fmt.Printf("Data backup path: %s\n", path)
		}

		if backupTime, ok := info["data_backup_time"]; ok {
			if t, ok := backupTime.(time.Time); ok {
				fmt.Printf("Data backup time: %s\n", t.Format("2006-01-02 15:04:05"))
			}
		}

		// Mostrar log de rollback si existe
		if log, err := rollbackManager.GetRollbackLog(); err == nil && log != "" {
			fmt.Println("\n📋 Recent rollback activity:")
			fmt.Println(log)
		}

	} else {
		printInfo("No rollback available")
		fmt.Println("This means either:")
		fmt.Println("- This is a fresh installation")
		fmt.Println("- No previous version was backed up")
		fmt.Println("- Backup files were cleaned up")
	}
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

// techCmd es el comando para agregar tareas técnicas
var techCmd = &cobra.Command{
	Use:   "tech [description] [hours]",
	Short: "Add a technical task",
	Long: `Add a technical task (development, coding, etc.).

Examples:
  harvest tech "Fix bug" 2.0
  harvest tech "Development" 3.5`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		description := args[0]
		hoursStr := args[1]

		// Validar horas
		hours, err := parseHours(hoursStr)
		if err != nil {
			printError(fmt.Errorf("invalid hours: %s", hoursStr))
			return
		}

		// Agregar tarea técnica
		taskManager := core.NewTaskManager()
		if err := taskManager.AddTask(description, hours, "tech"); err != nil {
			printError(err)
			return
		}

		printSuccess(fmt.Sprintf("Added tech task: %s (%.1fh)", description, hours))
		showStatus(taskManager)
	},
}

// meetingCmd es el comando para agregar reuniones
var meetingCmd = &cobra.Command{
	Use:   "meeting [description] [hours]",
	Short: "Add a meeting task",
	Long: `Add a meeting task (team sync, planning, etc.).

Examples:
  harvest meeting "Team sync" 1.0
  harvest meeting "Planning" 2.0`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		description := args[0]
		hoursStr := args[1]

		// Validar horas
		hours, err := parseHours(hoursStr)
		if err != nil {
			printError(fmt.Errorf("invalid hours: %s", hoursStr))
			return
		}

		// Agregar tarea de reunión
		taskManager := core.NewTaskManager()
		if err := taskManager.AddTask(description, hours, "meeting"); err != nil {
			printError(err)
			return
		}

		printSuccess(fmt.Sprintf("Added meeting: %s (%.1fh)", description, hours))
		showStatus(taskManager)
	},
}

// qaCmd es el comando para agregar tareas de QA
var qaCmd = &cobra.Command{
	Use:   "qa [description] [hours]",
	Short: "Add a QA/testing task",
	Long: `Add a QA/testing task (testing, quality assurance, etc.).

Examples:
  harvest qa "Testing" 1.5
  harvest qa "Bug fixes" 2.0`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		description := args[0]
		hoursStr := args[1]

		// Validar horas
		hours, err := parseHours(hoursStr)
		if err != nil {
			printError(fmt.Errorf("invalid hours: %s", hoursStr))
			return
		}

		// Agregar tarea de QA
		taskManager := core.NewTaskManager()
		if err := taskManager.AddTask(description, hours, "qa"); err != nil {
			printError(err)
			return
		}

		printSuccess(fmt.Sprintf("Added QA task: %s (%.1fh)", description, hours))
		showStatus(taskManager)
	},
}

// dailyCmd es el comando para agregar daily standup
var dailyCmd = &cobra.Command{
	Use:   "daily",
	Short: "Add daily standup",
	Long: `Add daily standup task (automatically 0.25h).

This command automatically adds a daily standup task with the default
duration configured in your settings.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		// Obtener duración del daily desde configuración
		taskManager := core.NewTaskManager()
		dailyHours := taskManager.GetDailyStandupHours()

		// Agregar tarea de daily
		if err := taskManager.AddTask("Daily Standup", dailyHours, "daily"); err != nil {
			printError(err)
			return
		}

		printSuccess(fmt.Sprintf("Added daily standup (%.2fh)", dailyHours))
		showStatus(taskManager)
	},
}

// reportCmd es el comando para generar reportes
var reportCmd = &cobra.Command{
	Use:   "report",
	Short: "Generate report for Harvest",
	Long: `Generate a formatted report for Harvest.

The report shows all today's tasks in the format:
"Description - X.Xh"

This format is ready to copy and paste into Harvest.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		taskManager := core.NewTaskManager()
		generateReport(taskManager)
	},
}

// generateReport genera el reporte para Harvest
func generateReport(taskManager *core.TaskManager) {
	todayTasks, err := taskManager.GetTodayTasks()
	if err != nil {
		printError(err)
		return
	}

	if len(todayTasks) == 0 {
		printInfo("No tasks for today")
		return
	}

	totalHours := taskManager.GetTotalHours(todayTasks)

	fmt.Printf("📋 Harvest Report for %s\n", time.Now().Format("2006-01-02"))
	fmt.Printf("Total hours: %.2fh\n\n", totalHours)

	fmt.Println("Copy the following lines to Harvest:")
	fmt.Println(strings.Repeat("─", 50))

	for _, task := range todayTasks {
		fmt.Printf("%s - %.1fh\n", task.Description, task.Hours)
	}

	fmt.Println(strings.Repeat("─", 50))

	// Intentar copiar al portapapeles (opcional)
	if err := copyToClipboard(todayTasks); err != nil {
		printInfo("Note: Could not copy to clipboard automatically")
	} else {
		printSuccess("Report copied to clipboard!")
	}
}

// copyToClipboard intenta copiar el reporte al portapapeles
func copyToClipboard(tasks []harvest.Task) error {
	// Construir el texto del reporte
	var reportText strings.Builder
	for _, task := range tasks {
		reportText.WriteString(fmt.Sprintf("%s - %.1fh\n", task.Description, task.Hours))
	}

	// Intentar usar xclip en Linux
	cmd := exec.Command("xclip", "-selection", "clipboard")
	cmd.Stdin = strings.NewReader(reportText.String())

	if err := cmd.Run(); err != nil {
		// Si xclip falla, intentar con xsel
		cmd = exec.Command("xsel", "--input", "--clipboard")
		cmd.Stdin = strings.NewReader(reportText.String())
		return cmd.Run()
	}

	return nil
}

// upgradeCmd es el comando para actualizar Harvest CLI
var upgradeCmd = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade to latest version",
	Long: `Upgrade Harvest CLI to the latest version.

This command will:
1. Check for available updates
2. Backup your current data
3. Download the latest version
4. Install the new version
5. Restore your data

If you're currently using the Python version, this will migrate you to Go.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		performUpgrade()
	},
}

// performUpgrade ejecuta el proceso de upgrade
func performUpgrade() {
	fmt.Println("🔍 Checking for updates...")

	// Crear gestor de versiones
	vm := upgrade.NewVersionManager()

	// Obtener información de upgrade
	info, err := vm.GetUpgradeInfo()
	if err != nil {
		printError(fmt.Errorf("could not check for updates: %v", err))
		return
	}

	// Mostrar información actual
	fmt.Printf("📊 Current version: %s\n", info.CurrentVersion)
	fmt.Printf("📦 Latest version: %s\n", info.LatestVersion)

	if !info.HasUpdate {
		printSuccess("You are currently on the latest version!")
		return
	}

	// Mostrar información de migración
	if !info.IsGoInstallation {
		fmt.Println("\n🔄 Migration detected: Python → Go")
		fmt.Println("This will migrate your installation from Python to Go while preserving all your data.")
	} else {
		fmt.Println("\n⬆️  Update available")
		fmt.Println("This will update your Go installation to the latest version.")
	}

	// Confirmar upgrade
	fmt.Print("\nDo you want to proceed with the upgrade? (y/N): ")
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(strings.TrimSpace(response)) != "y" && strings.ToLower(strings.TrimSpace(response)) != "yes" {
		printInfo("Upgrade cancelled")
		return
	}

	// Crear backup antes de proceder
	fmt.Println("\n💾 Creating backup of your data...")
	backupManager := upgrade.NewBackupManager()

	if err := backupManager.CreateBackup(); err != nil {
		printError(fmt.Errorf("could not create backup: %v", err))
		return
	}

	// Verificar integridad del backup
	if err := backupManager.VerifyBackup(); err != nil {
		printError(fmt.Errorf("backup verification failed: %v", err))
		return
	}

	backupPath := backupManager.GetBackupPath()
	printSuccess(fmt.Sprintf("Backup created successfully at: %s", backupPath))

	// Descargar nueva versión
	fmt.Println("\n📥 Downloading latest version...")
	downloadManager := upgrade.NewDownloadManager()

	// Obtener URL de descarga
	downloadURL, err := downloadManager.GetDownloadURL(info.LatestVersion)
	if err != nil {
		printError(fmt.Errorf("could not get download URL: %v", err))
		return
	}

	fmt.Printf("Download URL: %s\n", downloadURL)

	// Instalar nueva versión
	fmt.Println("\n🔧 Installing new version...")
	installManager := upgrade.NewInstallManager()

	// Por ahora, solo mostrar información de instalación
	fmt.Println("\n🚧 Installation system is under development")
	fmt.Println("This feature will be available in the next release.")
	fmt.Println("For now, you can manually download and install from:")
	fmt.Printf("https://github.com/%s/%s/releases\n", upgrade.RepoOwner, upgrade.RepoName)
	fmt.Printf("\nYour data has been safely backed up to: %s\n", backupPath)
	fmt.Printf("Installation path: %s\n", installManager.GetInstallPath())

	// Verificar disponibilidad de rollback
	rollbackManager := upgrade.NewRollbackManager()
	if rollbackManager.IsRollbackAvailable() {
		fmt.Println("\n🛡️  Rollback protection available")
		fmt.Println("If anything goes wrong, the system can automatically restore the previous version.")
	} else {
		fmt.Println("\n⚠️  No rollback protection available")
		fmt.Println("This is a fresh installation or no previous version was backed up.")
	}
}
