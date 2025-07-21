package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/lucasvidela94/workflow-cli/pkg/workflow"
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
func (tm *TaskManager) LoadTasks() ([]workflow.Task, error) {
	// Crear directorio si no existe
	dir := filepath.Dir(tm.dataFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("could not create data directory: %v", err)
	}

	// Si el archivo no existe, crear uno vacío
	if _, err := os.Stat(tm.dataFile); os.IsNotExist(err) {
		emptyTasks := []workflow.Task{}
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
	var tasks []workflow.Task

	// Decodificar como array de mapas primero
	var rawTasks []map[string]interface{}
	if err := decoder.Decode(&rawTasks); err != nil {
		return nil, fmt.Errorf("could not decode tasks file: %v", err)
	}

	// Convertir cada tarea
	for _, rawTask := range rawTasks {
		task := workflow.Task{}

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

		// Status (con valor por defecto para tareas existentes)
		if status, ok := rawTask["status"].(string); ok {
			task.Status = status
		} else {
			task.Status = workflow.StatusPending // Valor por defecto para tareas existentes
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
func (tm *TaskManager) SaveTasks(tasks []workflow.Task) error {
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
func (tm *TaskManager) AddTask(description string, hours float64, category string, date string) error {
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

	// Usar fecha proporcionada o fecha actual
	taskDate := date
	if taskDate == "" {
		taskDate = time.Now().Format("2006-01-02")
	}

	// Crear nueva tarea
	newTask := workflow.Task{
		ID:          newID,
		Description: description,
		Hours:       hours,
		Category:    category,
		Date:        taskDate,
		Status:      workflow.StatusPending,
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

// UpdateTask actualiza una tarea existente por ID
func (tm *TaskManager) UpdateTask(id int, description string, hours float64, category string) error {
	tasks, err := tm.LoadTasks()
	if err != nil {
		return fmt.Errorf("could not load tasks: %v", err)
	}

	// Buscar la tarea por ID
	taskIndex := -1
	for i, task := range tasks {
		if task.ID == id {
			taskIndex = i
			break
		}
	}

	if taskIndex == -1 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	// Actualizar la tarea
	if description != "" {
		tasks[taskIndex].Description = description
	}
	if hours > 0 {
		tasks[taskIndex].Hours = hours
	}
	if category != "" {
		tasks[taskIndex].Category = category
	}

	// Guardar cambios
	if err := tm.SaveTasks(tasks); err != nil {
		return fmt.Errorf("could not save tasks: %v", err)
	}

	return nil
}

// DeleteTask elimina una tarea por ID
func (tm *TaskManager) DeleteTask(id int) error {
	tasks, err := tm.LoadTasks()
	if err != nil {
		return fmt.Errorf("could not load tasks: %v", err)
	}

	// Buscar y eliminar la tarea por ID
	taskIndex := -1
	for i, task := range tasks {
		if task.ID == id {
			taskIndex = i
			break
		}
	}

	if taskIndex == -1 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	// Eliminar la tarea
	tasks = append(tasks[:taskIndex], tasks[taskIndex+1:]...)

	// Guardar cambios
	if err := tm.SaveTasks(tasks); err != nil {
		return fmt.Errorf("could not save tasks: %v", err)
	}

	return nil
}

// GetTaskByID obtiene una tarea específica por ID
func (tm *TaskManager) GetTaskByID(id int) (*workflow.Task, error) {
	tasks, err := tm.LoadTasks()
	if err != nil {
		return nil, fmt.Errorf("could not load tasks: %v", err)
	}

	for _, task := range tasks {
		if task.ID == id {
			return &task, nil
		}
	}

	return nil, fmt.Errorf("task with ID %d not found", id)
}

// GetTodayTasks obtiene las tareas del día actual
func (tm *TaskManager) GetTodayTasks() ([]workflow.Task, error) {
	tasks, err := tm.LoadTasks()
	if err != nil {
		return nil, err
	}

	today := time.Now().Format("2006-01-02")
	var todayTasks []workflow.Task

	for _, task := range tasks {
		if task.Date == today {
			todayTasks = append(todayTasks, task)
		}
	}

	return todayTasks, nil
}

// GetTotalHours calcula el total de horas de una lista de tareas
func (tm *TaskManager) GetTotalHours(tasks []workflow.Task) float64 {
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

// GetTasksByDate obtiene las tareas de una fecha específica
func (tm *TaskManager) GetTasksByDate(date string) ([]workflow.Task, error) {
	tasks, err := tm.LoadTasks()
	if err != nil {
		return nil, err
	}

	var filteredTasks []workflow.Task
	for _, task := range tasks {
		if task.Date == date {
			filteredTasks = append(filteredTasks, task)
		}
	}
	return filteredTasks, nil
}

// CompleteTask marca una tarea como completada
func (tm *TaskManager) CompleteTask(id int) error {
	tasks, err := tm.LoadTasks()
	if err != nil {
		return fmt.Errorf("could not load tasks: %v", err)
	}

	// Buscar la tarea por ID
	taskIndex := -1
	for i, task := range tasks {
		if task.ID == id {
			taskIndex = i
			break
		}
	}

	if taskIndex == -1 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	// Marcar como completada
	tasks[taskIndex].Status = workflow.StatusCompleted

	// Guardar cambios
	if err := tm.SaveTasks(tasks); err != nil {
		return fmt.Errorf("could not save tasks: %v", err)
	}

	return nil
}

// UpdateTaskStatus actualiza el estado de una tarea
func (tm *TaskManager) UpdateTaskStatus(id int, status string) error {
	tasks, err := tm.LoadTasks()
	if err != nil {
		return fmt.Errorf("could not load tasks: %v", err)
	}

	// Buscar la tarea por ID
	taskIndex := -1
	for i, task := range tasks {
		if task.ID == id {
			taskIndex = i
			break
		}
	}

	if taskIndex == -1 {
		return fmt.Errorf("task with ID %d not found", id)
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
	tasks[taskIndex].Status = status

	// Guardar cambios
	if err := tm.SaveTasks(tasks); err != nil {
		return fmt.Errorf("could not save tasks: %v", err)
	}

	return nil
}

// SearchTasks busca tareas según criterios específicos
func (tm *TaskManager) SearchTasks(query string, category string, status string, date string) ([]workflow.Task, error) {
	tasks, err := tm.LoadTasks()
	if err != nil {
		return nil, err
	}

	var filteredTasks []workflow.Task

	for _, task := range tasks {
		// Filtro por texto (descripción)
		if query != "" {
			if !strings.Contains(strings.ToLower(task.Description), strings.ToLower(query)) {
				continue
			}
		}

		// Filtro por categoría
		if category != "" {
			if task.Category != category {
				continue
			}
		}

		// Filtro por estado
		if status != "" {
			if task.Status != status {
				continue
			}
		}

		// Filtro por fecha
		if date != "" {
			if task.Date != date {
				continue
			}
		}

		filteredTasks = append(filteredTasks, task)
	}

	// Ordenar por fecha (más reciente primero) y luego por ID
	sort.Slice(filteredTasks, func(i, j int) bool {
		if filteredTasks[i].Date != filteredTasks[j].Date {
			return filteredTasks[i].Date > filteredTasks[j].Date
		}
		return filteredTasks[i].ID > filteredTasks[j].ID
	})

	return filteredTasks, nil
}
