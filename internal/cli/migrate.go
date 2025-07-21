package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lucasvidela94/workflow-cli/internal/core"
	"github.com/lucasvidela94/workflow-cli/pkg/workflow"
	"github.com/spf13/cobra"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate data from JSON to SQLite",
	Long: `Migrate existing data from JSON format to SQLite database.

This command will:
1. Detect existing JSON data
2. Create a backup of the JSON file
3. Migrate all tasks to SQLite
4. Verify the migration was successful

Examples:
  workflow migrate
  workflow migrate --dry-run
  workflow migrate --backup-only
`,
	Run: func(cmd *cobra.Command, args []string) {
		dryRun, _ := cmd.Flags().GetBool("dry-run")
		backupOnly, _ := cmd.Flags().GetBool("backup-only")

		if dryRun {
			performDryRunMigration()
		} else if backupOnly {
			performBackupOnly()
		} else {
			performMigration()
		}
	},
}

func performDryRunMigration() {
	fmt.Println("🔍 Performing dry run migration...")

	jsonPath := getJSONDataPath()
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		printInfo("No JSON data file found. Nothing to migrate.")
		return
	}

	jsonManager := core.NewTaskManager()
	tasks, err := jsonManager.LoadTasks()
	if err != nil {
		printError(fmt.Errorf("could not load JSON data: %v", err))
		return
	}

	fmt.Printf("📊 Found %d tasks in JSON format\n", len(tasks))
	fmt.Println("✅ Dry run completed successfully!")
	fmt.Println("📝 To perform actual migration, run: workflow migrate")
}

// performBackupOnly crea solo un backup del archivo JSON
func performBackupOnly() {
	fmt.Println("💾 Creating backup of JSON data...")

	jsonPath := getJSONDataPath()
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		printInfo("No JSON data file found. Nothing to backup.")
		return
	}

	backupPath := jsonPath + ".backup." + time.Now().Format("20060102-150405")

	// Copiar archivo
	if err := copyFile(jsonPath, backupPath); err != nil {
		printError(fmt.Errorf("could not create backup: %v", err))
		return
	}

	printSuccess(fmt.Sprintf("Backup created: %s", backupPath))
}

// performMigration ejecuta la migración completa
func performMigration() {
	fmt.Println("🔄 Starting migration from JSON to SQLite...")

	// Verificar si existe archivo JSON
	jsonPath := getJSONDataPath()
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		printInfo("No JSON data file found. Nothing to migrate.")
		return
	}

	// Crear backup antes de migrar
	fmt.Println("💾 Creating backup...")
	backupPath := jsonPath + ".backup." + time.Now().Format("20060102-150405")
	if err := copyFile(jsonPath, backupPath); err != nil {
		printError(fmt.Errorf("could not create backup: %v", err))
		return
	}
	printSuccess(fmt.Sprintf("Backup created: %s", backupPath))

	// Cargar datos JSON
	fmt.Println("📥 Loading JSON data...")
	jsonManager := core.NewTaskManager()
	tasks, err := jsonManager.LoadTasks()
	if err != nil {
		printError(fmt.Errorf("could not load JSON data: %v", err))
		return
	}

	fmt.Printf("📊 Found %d tasks to migrate\n", len(tasks))

	// Inicializar SQLite
	fmt.Println("🗄️  Initializing SQLite database...")
	sqliteManager := core.NewTaskManagerSQLite()
	defer sqliteManager.Close()

	// Migrar cada tarea
	fmt.Println("🔄 Migrating tasks...")
	migratedCount := 0
	for _, task := range tasks {
		// Crear nueva tarea en SQLite
		newTask := &workflow.Task{
			Description: task.Description,
			Hours:       task.Hours,
			Category:    task.Category,
			Date:        task.Date,
			Status:      task.Status,
			CreatedAt:   task.CreatedAt,
		}

		if err := sqliteManager.SaveTaskToDatabase(newTask); err != nil {
			printError(fmt.Errorf("could not migrate task %d: %v", task.ID, err))
			continue
		}

		migratedCount++
	}

	// Verificar migración
	fmt.Println("✅ Verifying migration...")
	sqliteTasks, err := sqliteManager.LoadTasks()
	if err != nil {
		printError(fmt.Errorf("could not verify migration: %v", err))
		return
	}

	if len(sqliteTasks) != len(tasks) {
		printError(fmt.Errorf("migration verification failed: expected %d tasks, got %d", len(tasks), len(sqliteTasks)))
		return
	}

	printSuccess(fmt.Sprintf("Migration completed successfully! %d tasks migrated.", migratedCount))
	fmt.Printf("📁 JSON backup: %s\n", backupPath)
	fmt.Printf("🗄️  SQLite database: %s\n", sqliteManager.GetDatabasePath())
	fmt.Println("📝 You can now safely delete the JSON file if everything works correctly.")
}

func getJSONDataPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".workflow", "tasks.json")
}

func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	return os.WriteFile(dst, input, 0644)
}

func init() {
	migrateCmd.Flags().Bool("dry-run", false, "Simulate migration without making changes")
	migrateCmd.Flags().Bool("backup-only", false, "Only create backup, don't migrate")
}
