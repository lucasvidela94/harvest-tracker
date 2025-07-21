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
	rootCmd.SetHelpTemplate(`🌾 Harvest CLI - Enterprise Task Management Tool

Usage:
  harvest [command] [args...]

Available Commands:
  add         Add a new task with description and hours
  tech        Add a technical development task
  meeting     Add a meeting or collaboration task
  qa          Add a QA/testing task
  daily       Add daily standup meeting
  status      Show today's task status and progress
  report      Generate detailed report for Harvest
  list        List tasks with filters (date, category, status)
  search      Search tasks by text, category, or status
  edit        Edit existing task (description, hours, category)
  delete      Delete a task with confirmation
  complete    Mark task as completed
  duplicate   Duplicate a task with date options
  export      Export tasks to CSV/JSON format
  migrate     Migrate from JSON to SQLite database
  upgrade     Upgrade to latest version
  check-update Check for available updates
  version     Show version information

Enterprise Features:
  • SQLite Database: Fast and scalable task storage
  • Auto-Update: Automatic version management
  • Multiplatform: Linux, macOS, Windows support
  • Professional Distribution: One-liner installation

Examples:
  harvest add "Fix critical bug" 2.0
  harvest tech "API development" 4.0
  harvest meeting "Sprint planning" 1.5
  harvest status
  harvest report
  harvest list --date 2025-07-21
  harvest search "bug"
  harvest export --format csv

Installation:
  curl -fsSL https://raw.githubusercontent.com/lucasvidela94/harvest-tracker/main/install-latest.sh | bash

Use "harvest [command] --help" for more information about a command.
`)

	// Comando de versión
	rootCmd.AddCommand(versionCmd)

	// Comando add
	addCmd.Flags().String("date", "", "Specific date for the task (format: YYYY-MM-DD)")
	addCmd.Flags().Bool("yesterday", false, "Add task for yesterday")
	addCmd.Flags().Bool("tomorrow", false, "Add task for tomorrow")
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

	// Comando list
	listCmd.Flags().String("date", "", "Date to list tasks for (format: YYYY-MM-DD)")
	rootCmd.AddCommand(listCmd)

	// Flags para edit
	editCmd.Flags().String("description", "", "New description for the task")
	editCmd.Flags().String("hours", "", "New hours for the task")
	editCmd.Flags().String("category", "", "New category for the task")

	// Flags para delete
	deleteCmd.Flags().Bool("force", false, "Force deletion without confirmation")

	// Agregar comandos
	rootCmd.AddCommand(editCmd)
	rootCmd.AddCommand(deleteCmd)

	// Flags para complete
	completeCmd.Flags().Bool("force", false, "Force completion without confirmation")

	// Agregar comando
	rootCmd.AddCommand(completeCmd)

	// Flags para search
	searchCmd.Flags().String("category", "", "Filter by category")
	searchCmd.Flags().String("status", "", "Filter by status (pending, in_progress, completed, paused)")
	searchCmd.Flags().String("date", "", "Filter by date (format: YYYY-MM-DD)")

	// Flags para report
	reportCmd.Flags().String("date", "", "Generate report for specific date (format: YYYY-MM-DD)")
	reportCmd.Flags().Bool("week", false, "Generate weekly report")
	reportCmd.Flags().Bool("month", false, "Generate monthly report")
	reportCmd.Flags().String("category", "", "Filter by category")
	reportCmd.Flags().String("status", "", "Filter by status (pending, in_progress, completed, paused)")
	reportCmd.Flags().Bool("harvest", false, "Generate legacy Harvest format report")

	// Agregar comando
	rootCmd.AddCommand(searchCmd)
	rootCmd.AddCommand(migrateCmd)
	rootCmd.AddCommand(duplicateCmd)
	rootCmd.AddCommand(exportCmd)
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
		fmt.Println("🌾 Harvest CLI v1.0.1")
		fmt.Println("Built with Go")
		fmt.Println("Migrated from Python to Go")
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
  harvest add "Meeting" 1.0 meeting
  harvest add --date 2025-07-20 "Tarea del lunes" 3.0
  harvest add --yesterday "Tarea olvidada" 2.0
  harvest add --tomorrow "Planificación" 1.5`,
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

		// Determinar fecha basada en flags
		date := ""
		dateFlag, _ := cmd.Flags().GetString("date")
		yesterdayFlag, _ := cmd.Flags().GetBool("yesterday")
		tomorrowFlag, _ := cmd.Flags().GetBool("tomorrow")

		// Validar que solo se use un flag de fecha
		dateFlagsCount := 0
		if dateFlag != "" {
			dateFlagsCount++
		}
		if yesterdayFlag {
			dateFlagsCount++
		}
		if tomorrowFlag {
			dateFlagsCount++
		}

		if dateFlagsCount > 1 {
			printError(fmt.Errorf("only one date flag can be used at a time (--date, --yesterday, --tomorrow)"))
			return
		}

		// Establecer fecha según flags
		if dateFlag != "" {
			date = dateFlag
		} else if yesterdayFlag {
			date = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		} else if tomorrowFlag {
			date = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		}

		// Agregar tarea
		taskManager := core.NewTaskManager()
		if err := taskManager.AddTask(description, hours, category, date); err != nil {
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
			statusIcon := harvest.GetStatusIcon(task.Status)
			fmt.Printf("  [%d] %s %s - %.1fh (%s) %s\n", task.ID, icon, task.Description, task.Hours, task.Category, statusIcon)
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
			statusIcon := harvest.GetStatusIcon(task.Status)
			fmt.Printf("  [%d] %s %s (%.1fh, %s) %s\n", task.ID, icon, task.Description, task.Hours, task.Category, statusIcon)
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
		if err := taskManager.AddTask(description, hours, "tech", ""); err != nil {
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
		if err := taskManager.AddTask(description, hours, "meeting", ""); err != nil {
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
		if err := taskManager.AddTask(description, hours, "qa", ""); err != nil {
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
		if err := taskManager.AddTask("Daily Standup", dailyHours, "daily", ""); err != nil {
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
	Short: "Generate detailed reports",
	Long: `Generate detailed reports of tasks and time tracking.

You can specify date ranges and filters to customize the report.

Examples:
  harvest report
  harvest report --date 2025-07-21
  harvest report --week
  harvest report --month
  harvest report --date 2025-07-21 --category tech
  harvest report --status completed
  harvest report --harvest (legacy format for Harvest app)`,
	Run: func(cmd *cobra.Command, args []string) {
		dateFlag, _ := cmd.Flags().GetString("date")
		weekFlag, _ := cmd.Flags().GetBool("week")
		monthFlag, _ := cmd.Flags().GetBool("month")
		categoryFlag, _ := cmd.Flags().GetString("category")
		statusFlag, _ := cmd.Flags().GetString("status")
		harvestFlag, _ := cmd.Flags().GetBool("harvest")

		// Validar que solo se use un flag de período
		periodFlagsCount := 0
		if dateFlag != "" {
			periodFlagsCount++
		}
		if weekFlag {
			periodFlagsCount++
		}
		if monthFlag {
			periodFlagsCount++
		}

		if periodFlagsCount > 1 {
			printError(fmt.Errorf("only one period flag can be used at a time (--date, --week, --month)"))
			return
		}

		if harvestFlag {
			// Formato legacy para Harvest
			taskManager := core.NewTaskManagerSQLite()
			defer taskManager.Close()
			generateHarvestReport(taskManager)
		} else {
			// Nuevo formato detallado
			performDetailedReport(dateFlag, weekFlag, monthFlag, categoryFlag, statusFlag)
		}
	},
}

// generateHarvestReport genera el reporte legacy para Harvest
func generateHarvestReport(taskManager *core.TaskManagerSQLite) {
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

// performDetailedReport ejecuta la generación del reporte detallado
func performDetailedReport(date string, week bool, month bool, category string, status string) {
	taskManager := core.NewTaskManagerSQLite()
	defer taskManager.Close()

	var tasks []harvest.Task
	var err error

	// Determinar período del reporte
	if date != "" {
		// Reporte de fecha específica
		tasks, err = taskManager.GetTasksByDate(date)
		if err != nil {
			printError(fmt.Errorf("could not load tasks for date %s: %v", date, err))
			return
		}
		generateDateReport(date, tasks, category, status)
	} else if week {
		// Reporte semanal
		startDate := getWeekStart()
		endDate := getWeekEnd()
		tasks, err = taskManager.SearchTasks("", category, status, "")
		if err != nil {
			printError(fmt.Errorf("could not load tasks: %v", err))
			return
		}
		// Filtrar por semana
		var weekTasks []harvest.Task
		for _, task := range tasks {
			if task.Date >= startDate && task.Date <= endDate {
				weekTasks = append(weekTasks, task)
			}
		}
		generateWeekReport(startDate, endDate, weekTasks, category, status)
	} else if month {
		// Reporte mensual
		startDate := getMonthStart()
		endDate := getMonthEnd()
		tasks, err = taskManager.SearchTasks("", category, status, "")
		if err != nil {
			printError(fmt.Errorf("could not load tasks: %v", err))
			return
		}
		// Filtrar por mes
		var monthTasks []harvest.Task
		for _, task := range tasks {
			if task.Date >= startDate && task.Date <= endDate {
				monthTasks = append(monthTasks, task)
			}
		}
		generateMonthReport(startDate, endDate, monthTasks, category, status)
	} else {
		// Reporte de hoy por defecto
		today := time.Now().Format("2006-01-02")
		tasks, err = taskManager.GetTasksByDate(today)
		if err != nil {
			printError(fmt.Errorf("could not load today's tasks: %v", err))
			return
		}
		generateDateReport(today, tasks, category, status)
	}
}

// generateDateReport genera reporte para una fecha específica
func generateDateReport(date string, tasks []harvest.Task, category string, status string) {
	fmt.Printf("📊 Report for %s\n", date)
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	if len(tasks) == 0 {
		fmt.Println("📝 No tasks found for this date.")
		return
	}

	// Aplicar filtros
	var filteredTasks []harvest.Task
	for _, task := range tasks {
		if category != "" && task.Category != category {
			continue
		}
		if status != "" && task.Status != status {
			continue
		}
		filteredTasks = append(filteredTasks, task)
	}

	if len(filteredTasks) == 0 {
		fmt.Println("📝 No tasks match the specified filters.")
		return
	}

	// Mostrar tareas
	fmt.Printf("📋 Tasks (%d):\n", len(filteredTasks))
	for _, task := range filteredTasks {
		fmt.Printf("[%d] %s %s (%.1fh, %s) %s\n",
			task.ID,
			harvest.GetIcon(task.Category),
			task.Description,
			task.Hours,
			task.Category,
			harvest.GetStatusIcon(task.Status))
	}

	// Estadísticas
	totalHours := 0.0
	completedHours := 0.0
	pendingHours := 0.0
	categoryStats := make(map[string]float64)

	for _, task := range filteredTasks {
		totalHours += task.Hours
		categoryStats[task.Category] += task.Hours

		if task.Status == harvest.StatusCompleted {
			completedHours += task.Hours
		} else {
			pendingHours += task.Hours
		}
	}

	fmt.Printf("\n📈 Statistics:\n")
	fmt.Printf("Total hours: %.1fh\n", totalHours)
	fmt.Printf("Completed: %.1fh\n", completedHours)
	fmt.Printf("Pending: %.1fh\n", pendingHours)

	if len(categoryStats) > 1 {
		fmt.Printf("\n📊 By category:\n")
		for category, hours := range categoryStats {
			fmt.Printf("  %s: %.1fh\n", category, hours)
		}
	}
}

// generateWeekReport genera reporte semanal
func generateWeekReport(startDate, endDate string, tasks []harvest.Task, category string, status string) {
	fmt.Printf("📊 Weekly Report (%s to %s)\n", startDate, endDate)
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	if len(tasks) == 0 {
		fmt.Println("📝 No tasks found for this week.")
		return
	}

	// Agrupar por fecha
	dateGroups := make(map[string][]harvest.Task)
	for _, task := range tasks {
		dateGroups[task.Date] = append(dateGroups[task.Date], task)
	}

	// Mostrar por día
	for date := startDate; date <= endDate; date = addDays(date, 1) {
		if dayTasks, exists := dateGroups[date]; exists {
			fmt.Printf("\n📅 %s:\n", date)
			totalDayHours := 0.0
			for _, task := range dayTasks {
				fmt.Printf("  [%d] %s %s (%.1fh, %s) %s\n",
					task.ID,
					harvest.GetIcon(task.Category),
					task.Description,
					task.Hours,
					task.Category,
					harvest.GetStatusIcon(task.Status))
				totalDayHours += task.Hours
			}
			fmt.Printf("  Total: %.1fh\n", totalDayHours)
		}
	}

	// Estadísticas semanales
	totalHours := 0.0
	completedHours := 0.0
	categoryStats := make(map[string]float64)

	for _, task := range tasks {
		totalHours += task.Hours
		categoryStats[task.Category] += task.Hours
		if task.Status == harvest.StatusCompleted {
			completedHours += task.Hours
		}
	}

	fmt.Printf("\n📈 Weekly Summary:\n")
	fmt.Printf("Total hours: %.1fh\n", totalHours)
	fmt.Printf("Completed: %.1fh\n", completedHours)
	fmt.Printf("Completion rate: %.1f%%\n", (completedHours/totalHours)*100)

	if len(categoryStats) > 1 {
		fmt.Printf("\n📊 By category:\n")
		for category, hours := range categoryStats {
			fmt.Printf("  %s: %.1fh\n", category, hours)
		}
	}
}

// generateMonthReport genera reporte mensual
func generateMonthReport(startDate, endDate string, tasks []harvest.Task, category string, status string) {
	fmt.Printf("📊 Monthly Report (%s to %s)\n", startDate, endDate)
	fmt.Printf("%s\n", strings.Repeat("=", 50))

	if len(tasks) == 0 {
		fmt.Println("📝 No tasks found for this month.")
		return
	}

	// Estadísticas mensuales
	totalHours := 0.0
	completedHours := 0.0
	categoryStats := make(map[string]float64)
	statusStats := make(map[string]int)

	for _, task := range tasks {
		totalHours += task.Hours
		categoryStats[task.Category] += task.Hours
		statusStats[task.Status]++

		if task.Status == harvest.StatusCompleted {
			completedHours += task.Hours
		}
	}

	fmt.Printf("📈 Monthly Summary:\n")
	fmt.Printf("Total hours: %.1fh\n", totalHours)
	fmt.Printf("Completed: %.1fh\n", completedHours)
	fmt.Printf("Completion rate: %.1f%%\n", (completedHours/totalHours)*100)
	fmt.Printf("Total tasks: %d\n", len(tasks))

	fmt.Printf("\n📊 By category:\n")
	for category, hours := range categoryStats {
		fmt.Printf("  %s: %.1fh\n", category, hours)
	}

	fmt.Printf("\n📊 By status:\n")
	for status, count := range statusStats {
		fmt.Printf("  %s: %d tasks\n", status, count)
	}
}

// Funciones auxiliares para fechas
func getWeekStart() string {
	now := time.Now()
	weekday := now.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	} else {
		weekday--
	}
	weekStart := now.AddDate(0, 0, -int(weekday))
	return weekStart.Format("2006-01-02")
}

func getWeekEnd() string {
	weekStart, _ := time.Parse("2006-01-02", getWeekStart())
	weekEnd := weekStart.AddDate(0, 0, 6)
	return weekEnd.Format("2006-01-02")
}

func getMonthStart() string {
	now := time.Now()
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
	return monthStart.Format("2006-01-02")
}

func getMonthEnd() string {
	now := time.Now()
	monthEnd := time.Date(now.Year(), now.Month()+1, 0, 0, 0, 0, 0, now.Location())
	return monthEnd.Format("2006-01-02")
}

func addDays(date string, days int) string {
	t, _ := time.Parse("2006-01-02", date)
	newDate := t.AddDate(0, 0, days)
	return newDate.Format("2006-01-02")
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

// listCmd es el comando para listar tareas
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List tasks for a specific date (default: today)",
	Long: `List all tasks for a given date (default: today).

Examples:
  harvest list
  harvest list --date 2025-07-20
`,
	Run: func(cmd *cobra.Command, args []string) {
		taskManager := core.NewTaskManagerSQLite()
		defer taskManager.Close()
		date, _ := cmd.Flags().GetString("date")
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}
		tasks, err := taskManager.GetTasksByDate(date)
		if err != nil {
			printError(err)
			return
		}
		fmt.Printf("\n📅 Tasks for %s:\n", date)
		if len(tasks) == 0 {
			fmt.Println("  No tasks found.")
			return
		}
		for _, task := range tasks {
			icon := harvest.GetIcon(task.Category)
			statusIcon := harvest.GetStatusIcon(task.Status)
			fmt.Printf("  [%d] %s %s (%.1fh, %s) %s\n", task.ID, icon, task.Description, task.Hours, task.Category, statusIcon)
		}
	},
}

// editCmd es el comando para editar tareas
var editCmd = &cobra.Command{
	Use:   "edit [id]",
	Short: "Edit a task by ID",
	Long: `Edit an existing task by its ID.

Examples:
  harvest edit 1 --description "New description"
  harvest edit 2 --hours 3.5
  harvest edit 3 --category tech
  harvest edit 1 --description "New desc" --hours 2.0 --category meeting
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			printError(fmt.Errorf("invalid task ID: %s", idStr))
			return
		}

		taskManager := core.NewTaskManager()

		// Obtener la tarea actual para mostrar información
		task, err := taskManager.GetTaskByID(id)
		if err != nil {
			printError(err)
			return
		}

		// Mostrar información actual
		icon := harvest.GetIcon(task.Category)
		fmt.Printf("✏️  Editing task:\n")
		fmt.Printf("[%d] %s %s (%.1fh, %s)\n\n", task.ID, icon, task.Description, task.Hours, task.Category)

		// Obtener nuevos valores de los flags
		description, _ := cmd.Flags().GetString("description")
		hoursStr, _ := cmd.Flags().GetString("hours")
		category, _ := cmd.Flags().GetString("category")

		// Parsear horas si se proporcionó
		var hours float64
		if hoursStr != "" {
			hours, err = parseHours(hoursStr)
			if err != nil {
				printError(fmt.Errorf("invalid hours: %s", hoursStr))
				return
			}
		}

		// Actualizar la tarea
		if err := taskManager.UpdateTask(id, description, hours, category); err != nil {
			printError(err)
			return
		}

		printSuccess(fmt.Sprintf("Task %d updated successfully", id))

		// Mostrar la tarea actualizada
		updatedTask, _ := taskManager.GetTaskByID(id)
		if updatedTask != nil {
			icon := harvest.GetIcon(updatedTask.Category)
			fmt.Printf("Updated: [%d] %s %s (%.1fh, %s)\n",
				updatedTask.ID, icon, updatedTask.Description, updatedTask.Hours, updatedTask.Category)
		}
	},
}

// deleteCmd es el comando para eliminar tareas
var deleteCmd = &cobra.Command{
	Use:   "delete [id]",
	Short: "Delete a task by ID",
	Long: `Delete an existing task by its ID.

Examples:
  harvest delete 1
  harvest delete 2 --force
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			printError(fmt.Errorf("invalid task ID: %s", idStr))
			return
		}

		taskManager := core.NewTaskManager()

		// Obtener la tarea para mostrar información
		task, err := taskManager.GetTaskByID(id)
		if err != nil {
			printError(err)
			return
		}

		// Mostrar información de la tarea a eliminar
		icon := harvest.GetIcon(task.Category)
		fmt.Printf("🗑️  Deleting task:\n")
		fmt.Printf("[%d] %s %s (%.1fh, %s)\n\n", task.ID, icon, task.Description, task.Hours, task.Category)

		// Verificar si se debe forzar la eliminación
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Print("Are you sure you want to delete this task? (y/N): ")
			var response string
			fmt.Scanln(&response)

			if strings.ToLower(strings.TrimSpace(response)) != "y" &&
				strings.ToLower(strings.TrimSpace(response)) != "yes" {
				printInfo("Deletion cancelled")
				return
			}
		}

		// Eliminar la tarea
		if err := taskManager.DeleteTask(id); err != nil {
			printError(err)
			return
		}

		printSuccess(fmt.Sprintf("Task %d deleted successfully", id))
	},
}

// completeCmd es el comando para marcar tareas como completadas
var completeCmd = &cobra.Command{
	Use:   "complete [id]",
	Short: "Mark a task as completed",
	Long: `Mark an existing task as completed by its ID.

Examples:
  harvest complete 1
  harvest complete 2 --force
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idStr := args[0]
		id, err := strconv.Atoi(idStr)
		if err != nil {
			printError(fmt.Errorf("invalid task ID: %s", idStr))
			return
		}

		taskManager := core.NewTaskManager()

		// Obtener la tarea para mostrar información
		task, err := taskManager.GetTaskByID(id)
		if err != nil {
			printError(err)
			return
		}

		// Verificar si ya está completada
		if task.Status == harvest.StatusCompleted {
			printInfo(fmt.Sprintf("Task %d is already completed", id))
			return
		}

		// Mostrar información de la tarea a completar
		icon := harvest.GetIcon(task.Category)
		statusIcon := harvest.GetStatusIcon(task.Status)
		fmt.Printf("✅ Completing task:\n")
		fmt.Printf("[%d] %s %s (%.1fh, %s) %s\n\n", task.ID, icon, task.Description, task.Hours, task.Category, statusIcon)

		// Verificar si se debe forzar la completación
		force, _ := cmd.Flags().GetBool("force")

		if !force {
			fmt.Print("Are you sure you want to mark this task as completed? (y/N): ")
			var response string
			fmt.Scanln(&response)

			if strings.ToLower(strings.TrimSpace(response)) != "y" &&
				strings.ToLower(strings.TrimSpace(response)) != "yes" {
				printInfo("Completion cancelled")
				return
			}
		}

		// Marcar como completada
		if err := taskManager.CompleteTask(id); err != nil {
			printError(err)
			return
		}

		printSuccess(fmt.Sprintf("Task %d marked as completed", id))
	},
}

// searchCmd es el comando para buscar tareas
var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search tasks by text, category, status, or date",
	Long: `Search tasks using various criteria.

Examples:
  harvest search "bug"
  harvest search "meeting" --category meeting
  harvest search "" --status completed
  harvest search "development" --category tech --status pending
  harvest search "test" --date 2025-07-21
`,
	Args: cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		query := ""
		if len(args) > 0 {
			query = args[0]
		}

		// Obtener filtros de los flags
		category, _ := cmd.Flags().GetString("category")
		status, _ := cmd.Flags().GetString("status")
		date, _ := cmd.Flags().GetString("date")

		taskManager := core.NewTaskManagerSQLite()
		defer taskManager.Close()
		tasks, err := taskManager.SearchTasks(query, category, status, date)
		if err != nil {
			printError(err)
			return
		}

		// Construir mensaje de búsqueda
		var searchTerms []string
		if query != "" {
			searchTerms = append(searchTerms, fmt.Sprintf("text: '%s'", query))
		}
		if category != "" {
			searchTerms = append(searchTerms, fmt.Sprintf("category: '%s'", category))
		}
		if status != "" {
			searchTerms = append(searchTerms, fmt.Sprintf("status: '%s'", status))
		}
		if date != "" {
			searchTerms = append(searchTerms, fmt.Sprintf("date: '%s'", date))
		}

		searchDescription := "all tasks"
		if len(searchTerms) > 0 {
			searchDescription = strings.Join(searchTerms, ", ")
		}

		fmt.Printf("🔍 Search results for %s:\n", searchDescription)

		if len(tasks) == 0 {
			fmt.Println("  No tasks found.")
			return
		}

		// Agrupar por fecha para mejor visualización
		currentDate := ""
		for _, task := range tasks {
			if task.Date != currentDate {
				currentDate = task.Date
				fmt.Printf("\n📅 %s:\n", currentDate)
			}

			icon := harvest.GetIcon(task.Category)
			statusIcon := harvest.GetStatusIcon(task.Status)
			fmt.Printf("  [%d] %s %s (%.1fh, %s) %s\n",
				task.ID, icon, task.Description, task.Hours, task.Category, statusIcon)
		}

		fmt.Printf("\n📊 Found %d task(s)\n", len(tasks))
	},
}
