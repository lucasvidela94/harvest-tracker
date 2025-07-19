package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lucasvidela94/harvest-cli/pkg/harvest"
)

// TaskManager maneja las operaciones de tareas
type TaskManager struct {
	configManager *ConfigManager
	dataFile      string
}

// NewTaskManager crea un nuevo gestor de tareas
func NewTaskManager() *TaskManager {
	configManager := NewConfigManager()
	if err := configManager.Load(); err != nil {
		// Si no puede cargar configuración, usar valores por defecto
		fmt.Printf("⚠️  Warning: Could not load config: %v\n", err)
	}

	return &TaskManager{
		configManager: configManager,
		dataFile:      configManager.GetDataFile(),
	}
}

// LoadTasks carga las tareas desde el archivo JSON
func (tm *TaskManager) LoadTasks() ([]harvest.Task, error) {
	// Crear directorio si no existe
	dir := filepath.Dir(tm.dataFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("could not create data directory: %v", err)
	}

	// Si el archivo no existe, crear uno vacío
	if _, err := os.Stat(tm.dataFile); os.IsNotExist(err) {
		emptyTasks := []harvest.Task{}
		if err := tm.SaveTasks(emptyTasks); err != nil {
			return nil, fmt.Errorf("could not create empty tasks file: %v", err)
		}
		return emptyTasks, nil
	}

	// Leer archivo existente
	file, err := os.Open(tm.dataFile)
	if err != nil {
		return nil, fmt.Errorf("could not open tasks file: %v", err)
	}
	defer file.Close()

	// Usar un decoder personalizado para manejar diferentes formatos de fecha
	decoder := json.NewDecoder(file)
	var tasks []harvest.Task

	// Decodificar como array de mapas primero
	var rawTasks []map[string]interface{}
	if err := decoder.Decode(&rawTasks); err != nil {
		return nil, fmt.Errorf("could not decode tasks file: %v", err)
	}

	// Convertir cada tarea
	for _, rawTask := range rawTasks {
		task := harvest.Task{}

		// ID
		if id, ok := rawTask["id"].(float64); ok {
			task.ID = int(id)
		}

		// Description
		if desc, ok := rawTask["description"].(string); ok {
			task.Description = desc
		}

		// Hours
		if hours, ok := rawTask["hours"].(float64); ok {
			task.Hours = hours
		}

		// Category
		if cat, ok := rawTask["category"].(string); ok {
			task.Category = cat
		}

		// Date
		if date, ok := rawTask["date"].(string); ok {
			task.Date = date
		}

		// CreatedAt - manejar diferentes formatos
		if createdAt, ok := rawTask["created_at"].(string); ok {
			// Intentar diferentes formatos de fecha
			formats := []string{
				"2006-01-02T15:04:05.999999",
				"2006-01-02T15:04:05Z07:00",
				"2006-01-02T15:04:05",
				"2006-01-02 15:04:05",
			}

			var parsedTime time.Time
			var parseErr error

			for _, format := range formats {
				if parsedTime, parseErr = time.Parse(format, createdAt); parseErr == nil {
					break
				}
			}

			if parseErr != nil {
				// Si no se puede parsear, usar tiempo actual
				task.CreatedAt = time.Now()
			} else {
				task.CreatedAt = parsedTime
			}
		} else {
			task.CreatedAt = time.Now()
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// SaveTasks guarda las tareas en el archivo JSON
func (tm *TaskManager) SaveTasks(tasks []harvest.Task) error {
	// Crear directorio si no existe
	dir := filepath.Dir(tm.dataFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create data directory: %v", err)
	}

	file, err := os.Create(tm.dataFile)
	if err != nil {
		return fmt.Errorf("could not create tasks file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(tasks); err != nil {
		return fmt.Errorf("could not encode tasks: %v", err)
	}

	return nil
}

// AddTask agrega una nueva tarea
func (tm *TaskManager) AddTask(description string, hours float64, category string) error {
	tasks, err := tm.LoadTasks()
	if err != nil {
		return fmt.Errorf("could not load tasks: %v", err)
	}

	// Generar nuevo ID
	newID := 1
	if len(tasks) > 0 {
		maxID := 0
		for _, task := range tasks {
			if task.ID > maxID {
				maxID = task.ID
			}
		}
		newID = maxID + 1
	}

	// Crear nueva tarea
	newTask := harvest.Task{
		ID:          newID,
		Description: description,
		Hours:       hours,
		Category:    category,
		Date:        time.Now().Format("2006-01-02"),
		CreatedAt:   time.Now(),
	}

	// Agregar a la lista
	tasks = append(tasks, newTask)

	// Guardar
	if err := tm.SaveTasks(tasks); err != nil {
		return fmt.Errorf("could not save tasks: %v", err)
	}

	return nil
}

// GetTodayTasks obtiene las tareas del día actual
func (tm *TaskManager) GetTodayTasks() ([]harvest.Task, error) {
	tasks, err := tm.LoadTasks()
	if err != nil {
		return nil, err
	}

	today := time.Now().Format("2006-01-02")
	var todayTasks []harvest.Task

	for _, task := range tasks {
		if task.Date == today {
			todayTasks = append(todayTasks, task)
		}
	}

	return todayTasks, nil
}

// GetTotalHours calcula el total de horas de una lista de tareas
func (tm *TaskManager) GetTotalHours(tasks []harvest.Task) float64 {
	total := 0.0
	for _, task := range tasks {
		total += task.Hours
	}
	return total
}

// GetDailyHoursTarget obtiene el objetivo de horas diarias
func (tm *TaskManager) GetDailyHoursTarget() float64 {
	return tm.configManager.GetDailyHoursTarget()
}

// GetDailyStandupHours obtiene las horas del daily standup
func (tm *TaskManager) GetDailyStandupHours() float64 {
	return tm.configManager.GetDailyStandupHours()
}
