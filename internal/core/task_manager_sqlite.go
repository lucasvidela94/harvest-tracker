package core

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lucasvidela94/workflow-cli/pkg/workflow"
)

// TaskManagerSQLite es una implementación del TaskManager usando SQLite
type TaskManagerSQLite struct {
	configManager *ConfigManager
	dbManager     *DatabaseManager
}

// NewTaskManagerSQLite crea un nuevo gestor de tareas con SQLite
func NewTaskManagerSQLite() *TaskManagerSQLite {
	configManager := NewConfigManager()
	if err := configManager.Load(); err != nil {
		// Si no puede cargar configuración, usar valores por defecto
		fmt.Printf("⚠️  Warning: Could not load config: %v\n", err)
	}

	// Obtener directorio de datos
	homeDir, _ := os.UserHomeDir()
	dataDir := filepath.Join(homeDir, ".workflow")

	dbManager := NewDatabaseManager(dataDir)
	if err := dbManager.Init(); err != nil {
		fmt.Printf("⚠️  Warning: Could not initialize database: %v\n", err)
	}

	return &TaskManagerSQLite{
		configManager: configManager,
		dbManager:     dbManager,
	}
}

// LoadTasks carga las tareas desde la base de datos
func (tm *TaskManagerSQLite) LoadTasks() ([]workflow.Task, error) {
	return tm.dbManager.LoadTasks()
}

// SaveTasks guarda las tareas en la base de datos (compatibilidad con interfaz JSON)
func (tm *TaskManagerSQLite) SaveTasks(tasks []workflow.Task) error {
	// Para SQLite, no necesitamos este método ya que cada tarea se guarda individualmente
	// Pero lo mantenemos para compatibilidad
	return nil
}

// AddTask agrega una nueva tarea
func (tm *TaskManagerSQLite) AddTask(description string, hours float64, category string, date string) error {
	// Usar fecha proporcionada o fecha actual
	taskDate := date
	if taskDate == "" {
		taskDate = time.Now().Format("2006-01-02")
	}

	// Crear nueva tarea
	newTask := &workflow.Task{
		Description: description,
		Hours:       hours,
		Category:    category,
		Date:        taskDate,
		Status:      workflow.StatusPending,
		CreatedAt:   time.Now(),
	}

	// Guardar en la base de datos
	return tm.dbManager.SaveTask(newTask)
}

// UpdateTask actualiza una tarea existente por ID
func (tm *TaskManagerSQLite) UpdateTask(id int, description string, hours float64, category string) error {
	// Obtener la tarea actual
	task, err := tm.dbManager.GetTaskByID(id)
	if err != nil {
		return err
	}

	// Actualizar campos si se proporcionan
	if description != "" {
		task.Description = description
	}
	if hours > 0 {
		task.Hours = hours
	}
	if category != "" {
		task.Category = category
	}

	// Guardar cambios
	return tm.dbManager.UpdateTask(task)
}

// DeleteTask elimina una tarea por ID
func (tm *TaskManagerSQLite) DeleteTask(id int) error {
	return tm.dbManager.DeleteTask(id)
}

// GetTaskByID obtiene una tarea específica por ID
func (tm *TaskManagerSQLite) GetTaskByID(id int) (*workflow.Task, error) {
	return tm.dbManager.GetTaskByID(id)
}

// GetTodayTasks obtiene las tareas del día actual
func (tm *TaskManagerSQLite) GetTodayTasks() ([]workflow.Task, error) {
	today := time.Now().Format("2006-01-02")
	return tm.dbManager.GetTasksByDate(today)
}

// GetTasksByDate obtiene las tareas de una fecha específica
func (tm *TaskManagerSQLite) GetTasksByDate(date string) ([]workflow.Task, error) {
	return tm.dbManager.GetTasksByDate(date)
}

// SearchTasks busca tareas según criterios específicos
func (tm *TaskManagerSQLite) SearchTasks(query string, category string, status string, date string) ([]workflow.Task, error) {
	return tm.dbManager.SearchTasks(query, category, status, date)
}

// CompleteTask marca una tarea como completada
func (tm *TaskManagerSQLite) CompleteTask(id int) error {
	// Obtener la tarea actual
	task, err := tm.dbManager.GetTaskByID(id)
	if err != nil {
		return err
	}

	// Marcar como completada
	task.Status = workflow.StatusCompleted

	// Guardar cambios
	return tm.dbManager.UpdateTask(task)
}

// UpdateTaskStatus actualiza el estado de una tarea
func (tm *TaskManagerSQLite) UpdateTaskStatus(id int, status string) error {
	// Obtener la tarea actual
	task, err := tm.dbManager.GetTaskByID(id)
	if err != nil {
		return err
	}

	// Validar estado
	validStatuses := []string{workflow.StatusPending, workflow.StatusInProgress, workflow.StatusCompleted, workflow.StatusPaused}
	isValid := false
	for _, validStatus := range validStatuses {
		if status == validStatus {
			isValid = true
			break
		}
	}

	if !isValid {
		return fmt.Errorf("invalid status: %s. Valid statuses are: %v", status, validStatuses)
	}

	// Actualizar estado
	task.Status = status

	// Guardar cambios
	return tm.dbManager.UpdateTask(task)
}

// GetTotalHours calcula el total de horas de una lista de tareas
func (tm *TaskManagerSQLite) GetTotalHours(tasks []workflow.Task) float64 {
	total := 0.0
	for _, task := range tasks {
		total += task.Hours
	}
	return total
}

// GetDailyHoursTarget obtiene el objetivo de horas diarias
func (tm *TaskManagerSQLite) GetDailyHoursTarget() float64 {
	return tm.configManager.GetDailyHoursTarget()
}

// GetDailyStandupHours obtiene las horas del daily standup
func (tm *TaskManagerSQLite) GetDailyStandupHours() float64 {
	return tm.configManager.GetDailyStandupHours()
}

// Close cierra la conexión a la base de datos
func (tm *TaskManagerSQLite) Close() error {
	return tm.dbManager.Close()
}

// GetDatabasePath devuelve la ruta de la base de datos
func (tm *TaskManagerSQLite) GetDatabasePath() string {
	return tm.dbManager.GetDatabasePath()
}

// SaveTaskToDatabase guarda una tarea directamente en la base de datos
func (tm *TaskManagerSQLite) SaveTaskToDatabase(task *workflow.Task) error {
	return tm.dbManager.SaveTask(task)
}
