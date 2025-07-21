package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/lucasvidela94/harvest-cli/internal/core"
	"github.com/lucasvidela94/harvest-cli/pkg/harvest"
	"github.com/spf13/cobra"
)

// parseTaskID convierte un string a ID de tarea
func parseTaskID(idStr string) int {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return -1
	}
	return id
}

// duplicateCmd es el comando para duplicar tareas
var duplicateCmd = &cobra.Command{
	Use:   "duplicate <task-id>",
	Short: "Duplicate an existing task",
	Long: `Duplicate an existing task with the same description, hours, and category.

You can specify a different date for the duplicated task using flags.

Examples:
  harvest duplicate 1
  harvest duplicate 1 --date 2025-07-22
  harvest duplicate 1 --tomorrow
  harvest duplicate 1 --yesterday
`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		taskID := parseTaskID(args[0])
		if taskID == -1 {
			printError(fmt.Errorf("invalid task ID: %s", args[0]))
			return
		}

		// Obtener flags
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

		// Determinar fecha para la tarea duplicada
		targetDate := ""
		if dateFlag != "" {
			targetDate = dateFlag
		} else if yesterdayFlag {
			targetDate = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		} else if tomorrowFlag {
			targetDate = time.Now().AddDate(0, 0, 1).Format("2006-01-02")
		}

		performDuplicate(taskID, targetDate)
	},
}

// performDuplicate ejecuta la duplicación de la tarea
func performDuplicate(taskID int, targetDate string) {
	taskManager := core.NewTaskManagerSQLite()
	defer taskManager.Close()

	// Obtener la tarea original
	originalTask, err := taskManager.GetTaskByID(taskID)
	if err != nil {
		printError(fmt.Errorf("could not find task %d: %v", taskID, err))
		return
	}

	// Mostrar tarea original
	fmt.Printf("📋 Original task:\n")
	fmt.Printf("[%d] %s %s (%.1fh, %s) %s\n",
		originalTask.ID,
		harvest.GetIcon(originalTask.Category),
		originalTask.Description,
		originalTask.Hours,
		originalTask.Category,
		harvest.GetStatusIcon(originalTask.Status))

	// Crear nueva tarea
	newTask := &harvest.Task{
		Description: originalTask.Description,
		Hours:       originalTask.Hours,
		Category:    originalTask.Category,
		Status:      harvest.StatusPending, // Siempre pendiente al duplicar
		CreatedAt:   time.Now(),
	}

	// Establecer fecha
	if targetDate != "" {
		newTask.Date = targetDate
	} else {
		newTask.Date = time.Now().Format("2006-01-02") // Hoy por defecto
	}

	// Guardar nueva tarea
	if err := taskManager.SaveTaskToDatabase(newTask); err != nil {
		printError(fmt.Errorf("could not duplicate task: %v", err))
		return
	}

	// Mostrar confirmación
	fmt.Printf("\n✅ Task duplicated successfully!\n")
	fmt.Printf("📝 New task:\n")
	fmt.Printf("[%d] %s %s (%.1fh, %s) %s - %s\n",
		newTask.ID,
		harvest.GetIcon(newTask.Category),
		newTask.Description,
		newTask.Hours,
		newTask.Category,
		harvest.GetStatusIcon(newTask.Status),
		newTask.Date)

	// Mostrar estado actualizado
	if newTask.Date == time.Now().Format("2006-01-02") {
		fmt.Printf("\n📊 Today's Status:\n")
		todayTasks, _ := taskManager.GetTodayTasks()
		totalHours := taskManager.GetTotalHours(todayTasks)
		targetHours := taskManager.GetDailyHoursTarget()
		remaining := targetHours - totalHours

		fmt.Printf("📅 Today (%s): %.1fh / %.1fh\n", newTask.Date, totalHours, targetHours)
		if remaining > 0 {
			fmt.Printf("📈 Remaining: %.1fh\n", remaining)
		} else {
			fmt.Printf("📈 Overtime: %.1fh\n", -remaining)
		}
	}
}

func init() {
	// Flags para duplicate
	duplicateCmd.Flags().String("date", "", "Target date for duplicated task (format: YYYY-MM-DD)")
	duplicateCmd.Flags().Bool("yesterday", false, "Duplicate task for yesterday")
	duplicateCmd.Flags().Bool("tomorrow", false, "Duplicate task for tomorrow")
}
