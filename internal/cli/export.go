package cli

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/lucasvidela94/workflow-cli/internal/core"
	"github.com/lucasvidela94/workflow-cli/pkg/workflow"
	"github.com/spf13/cobra"
)

// exportCmd es el comando para exportar datos
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export tasks to different formats",
	Long: `Export tasks to different file formats.

Supported formats:
- CSV: Comma-separated values
- JSON: JavaScript Object Notation

You can specify date ranges and filters to customize the export.

Examples:
  workflow export --format csv
  workflow export --format csv --date 2025-07-21
  workflow export --format csv --week
  workflow export --format csv --category tech
  workflow export --format json --status completed
`,
	Run: func(cmd *cobra.Command, args []string) {
		formatFlag, _ := cmd.Flags().GetString("format")
		dateFlag, _ := cmd.Flags().GetString("date")
		weekFlag, _ := cmd.Flags().GetBool("week")
		monthFlag, _ := cmd.Flags().GetBool("month")
		categoryFlag, _ := cmd.Flags().GetString("category")
		statusFlag, _ := cmd.Flags().GetString("status")
		outputFlag, _ := cmd.Flags().GetString("output")

		// Validar formato
		if formatFlag != "csv" && formatFlag != "json" {
			printError(fmt.Errorf("unsupported format: %s. Supported formats: csv, json", formatFlag))
			return
		}

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

		performExport(formatFlag, dateFlag, weekFlag, monthFlag, categoryFlag, statusFlag, outputFlag)
	},
}

// performExport ejecuta la exportación
func performExport(format, date string, week, month bool, category, status, output string) {
	taskManager := core.NewTaskManagerSQLite()
	defer taskManager.Close()

	var tasks []workflow.Task
	var err error

	// Determinar período de exportación
	if date != "" {
		// Exportar fecha específica
		tasks, err = taskManager.GetTasksByDate(date)
		if err != nil {
			printError(fmt.Errorf("could not load tasks for date %s: %v", date, err))
			return
		}
	} else if week {
		// Exportar semana
		startDate := getWeekStart()
		endDate := getWeekEnd()
		tasks, err = taskManager.SearchTasks("", category, status, "")
		if err != nil {
			printError(fmt.Errorf("could not load tasks: %v", err))
			return
		}
		// Filtrar por semana
		var weekTasks []workflow.Task
		for _, task := range tasks {
			if task.Date >= startDate && task.Date <= endDate {
				weekTasks = append(weekTasks, task)
			}
		}
		tasks = weekTasks
	} else if month {
		// Exportar mes
		startDate := getMonthStart()
		endDate := getMonthEnd()
		tasks, err = taskManager.SearchTasks("", category, status, "")
		if err != nil {
			printError(fmt.Errorf("could not load tasks: %v", err))
			return
		}
		// Filtrar por mes
		var monthTasks []workflow.Task
		for _, task := range tasks {
			if task.Date >= startDate && task.Date <= endDate {
				monthTasks = append(monthTasks, task)
			}
		}
		tasks = monthTasks
	} else {
		// Exportar todas las tareas
		tasks, err = taskManager.SearchTasks("", category, status, "")
		if err != nil {
			printError(fmt.Errorf("could not load tasks: %v", err))
			return
		}
	}

	// Aplicar filtros adicionales
	var filteredTasks []workflow.Task
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
		printInfo("No tasks found matching the specified criteria.")
		return
	}

	// Determinar nombre de archivo
	if output == "" {
		timestamp := time.Now().Format("20060102-150405")
		output = fmt.Sprintf("workflow-export-%s.%s", timestamp, format)
	}

	// Exportar según formato
	switch format {
	case "csv":
		err = exportToCSV(filteredTasks, output)
	case "json":
		err = exportToJSON(filteredTasks, output)
	}

	if err != nil {
		printError(fmt.Errorf("could not export to %s: %v", format, err))
		return
	}

	// Obtener ruta absoluta
	absPath, _ := filepath.Abs(output)
	printSuccess(fmt.Sprintf("Exported %d tasks to %s", len(filteredTasks), absPath))
}

// exportToCSV exporta las tareas a formato CSV
func exportToCSV(tasks []workflow.Task, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Escribir encabezados
	headers := []string{"ID", "Description", "Hours", "Category", "Date", "Status", "Created At"}
	if err := writer.Write(headers); err != nil {
		return err
	}

	// Escribir datos
	for _, task := range tasks {
		row := []string{
			strconv.Itoa(task.ID),
			task.Description,
			strconv.FormatFloat(task.Hours, 'f', 1, 64),
			task.Category,
			task.Date,
			task.Status,
			task.CreatedAt.Format("2006-01-02 15:04:05"),
		}
		if err := writer.Write(row); err != nil {
			return err
		}
	}

	return nil
}

// exportToJSON exporta las tareas a formato JSON
func exportToJSON(tasks []workflow.Task, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Escribir inicio del JSON
	file.WriteString("[\n")

	for i, task := range tasks {
		// Escribir tarea
		file.WriteString("  {\n")
		file.WriteString(fmt.Sprintf("    \"id\": %d,\n", task.ID))
		file.WriteString(fmt.Sprintf("    \"description\": \"%s\",\n", task.Description))
		file.WriteString(fmt.Sprintf("    \"hours\": %.1f,\n", task.Hours))
		file.WriteString(fmt.Sprintf("    \"category\": \"%s\",\n", task.Category))
		file.WriteString(fmt.Sprintf("    \"date\": \"%s\",\n", task.Date))
		file.WriteString(fmt.Sprintf("    \"status\": \"%s\",\n", task.Status))
		file.WriteString(fmt.Sprintf("    \"created_at\": \"%s\"\n", task.CreatedAt.Format("2006-01-02 15:04:05")))

		if i < len(tasks)-1 {
			file.WriteString("  },\n")
		} else {
			file.WriteString("  }\n")
		}
	}

	// Escribir fin del JSON
	file.WriteString("]\n")

	return nil
}

func init() {
	// Flags para export
	exportCmd.Flags().String("format", "csv", "Export format (csv, json)")
	exportCmd.Flags().String("date", "", "Export tasks for specific date (format: YYYY-MM-DD)")
	exportCmd.Flags().Bool("week", false, "Export weekly tasks")
	exportCmd.Flags().Bool("month", false, "Export monthly tasks")
	exportCmd.Flags().String("category", "", "Filter by category")
	exportCmd.Flags().String("status", "", "Filter by status (pending, in_progress, completed, paused)")
	exportCmd.Flags().String("output", "", "Output filename (default: workflow-export-YYYYMMDD-HHMMSS.format)")
}
