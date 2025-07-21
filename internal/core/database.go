package core

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/lucasvidela94/workflow-cli/pkg/workflow"
	_ "github.com/mattn/go-sqlite3"
)

// DatabaseManager maneja las operaciones de la base de datos SQLite
type DatabaseManager struct {
	dbPath string
	db     *sql.DB
}

// NewDatabaseManager crea un nuevo gestor de base de datos
func NewDatabaseManager(dataDir string) *DatabaseManager {
	dbPath := filepath.Join(dataDir, "tasks.db")
	return &DatabaseManager{
		dbPath: dbPath,
	}
}

// Init inicializa la base de datos y crea las tablas si no existen
func (dm *DatabaseManager) Init() error {
	// Crear directorio si no existe
	dir := filepath.Dir(dm.dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create database directory: %v", err)
	}

	// Abrir conexión a la base de datos
	db, err := sql.Open("sqlite3", dm.dbPath)
	if err != nil {
		return fmt.Errorf("could not open database: %v", err)
	}
	dm.db = db

	// Crear tabla si no existe
	if err := dm.createTables(); err != nil {
		return fmt.Errorf("could not create tables: %v", err)
	}

	return nil
}

// Close cierra la conexión a la base de datos
func (dm *DatabaseManager) Close() error {
	if dm.db != nil {
		return dm.db.Close()
	}
	return nil
}

// createTables crea las tablas necesarias
func (dm *DatabaseManager) createTables() error {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		hours REAL NOT NULL,
		category TEXT DEFAULT 'general',
		date TEXT NOT NULL,
		status TEXT DEFAULT 'pending',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE INDEX IF NOT EXISTS idx_tasks_date ON tasks(date);
	CREATE INDEX IF NOT EXISTS idx_tasks_category ON tasks(category);
	CREATE INDEX IF NOT EXISTS idx_tasks_status ON tasks(status);
	CREATE INDEX IF NOT EXISTS idx_tasks_description ON tasks(description);
	`

	_, err := dm.db.Exec(query)
	return err
}

// LoadTasks carga todas las tareas desde la base de datos
func (dm *DatabaseManager) LoadTasks() ([]workflow.Task, error) {
	query := `SELECT id, description, hours, category, date, status, created_at FROM tasks ORDER BY date DESC, id DESC`

	rows, err := dm.db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("could not query tasks: %v", err)
	}
	defer rows.Close()

	var tasks []workflow.Task
	for rows.Next() {
		var task workflow.Task
		var createdAtStr string

		err := rows.Scan(&task.ID, &task.Description, &task.Hours, &task.Category, &task.Date, &task.Status, &createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("could not scan task: %v", err)
		}

		// Parsear created_at
		if createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			task.CreatedAt = createdAt
		} else {
			task.CreatedAt = time.Now()
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// SaveTask guarda una nueva tarea en la base de datos
func (dm *DatabaseManager) SaveTask(task *workflow.Task) error {
	query := `
	INSERT INTO tasks (description, hours, category, date, status, created_at)
	VALUES (?, ?, ?, ?, ?, ?)
	`

	result, err := dm.db.Exec(query, task.Description, task.Hours, task.Category, task.Date, task.Status, task.CreatedAt)
	if err != nil {
		return fmt.Errorf("could not insert task: %v", err)
	}

	// Obtener el ID generado
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("could not get last insert id: %v", err)
	}

	task.ID = int(id)
	return nil
}

// UpdateTask actualiza una tarea existente
func (dm *DatabaseManager) UpdateTask(task *workflow.Task) error {
	query := `
	UPDATE tasks 
	SET description = ?, hours = ?, category = ?, date = ?, status = ?, updated_at = CURRENT_TIMESTAMP
	WHERE id = ?
	`

	result, err := dm.db.Exec(query, task.Description, task.Hours, task.Category, task.Date, task.Status, task.ID)
	if err != nil {
		return fmt.Errorf("could not update task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", task.ID)
	}

	return nil
}

// DeleteTask elimina una tarea por ID
func (dm *DatabaseManager) DeleteTask(id int) error {
	query := `DELETE FROM tasks WHERE id = ?`

	result, err := dm.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("could not delete task: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("could not get rows affected: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("task with ID %d not found", id)
	}

	return nil
}

// GetTaskByID obtiene una tarea específica por ID
func (dm *DatabaseManager) GetTaskByID(id int) (*workflow.Task, error) {
	query := `SELECT id, description, hours, category, date, status, created_at FROM tasks WHERE id = ?`

	var task workflow.Task
	var createdAtStr string

	err := dm.db.QueryRow(query, id).Scan(&task.ID, &task.Description, &task.Hours, &task.Category, &task.Date, &task.Status, &createdAtStr)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("task with ID %d not found", id)
		}
		return nil, fmt.Errorf("could not scan task: %v", err)
	}

	// Parsear created_at
	if createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
		task.CreatedAt = createdAt
	} else {
		task.CreatedAt = time.Now()
	}

	return &task, nil
}

// GetTasksByDate obtiene tareas de una fecha específica
func (dm *DatabaseManager) GetTasksByDate(date string) ([]workflow.Task, error) {
	query := `SELECT id, description, hours, category, date, status, created_at FROM tasks WHERE date = ? ORDER BY id`

	rows, err := dm.db.Query(query, date)
	if err != nil {
		return nil, fmt.Errorf("could not query tasks by date: %v", err)
	}
	defer rows.Close()

	var tasks []workflow.Task
	for rows.Next() {
		var task workflow.Task
		var createdAtStr string

		err := rows.Scan(&task.ID, &task.Description, &task.Hours, &task.Category, &task.Date, &task.Status, &createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("could not scan task: %v", err)
		}

		// Parsear created_at
		if createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			task.CreatedAt = createdAt
		} else {
			task.CreatedAt = time.Now()
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// SearchTasks busca tareas según criterios específicos
func (dm *DatabaseManager) SearchTasks(query string, category string, status string, date string) ([]workflow.Task, error) {
	baseQuery := `SELECT id, description, hours, category, date, status, created_at FROM tasks WHERE 1=1`
	var args []interface{}
	var conditions []string

	if query != "" {
		conditions = append(conditions, "description LIKE ?")
		args = append(args, "%"+query+"%")
	}

	if category != "" {
		conditions = append(conditions, "category = ?")
		args = append(args, category)
	}

	if status != "" {
		conditions = append(conditions, "status = ?")
		args = append(args, status)
	}

	if date != "" {
		conditions = append(conditions, "date = ?")
		args = append(args, date)
	}

	if len(conditions) > 0 {
		baseQuery += " AND " + conditions[0]
		for i := 1; i < len(conditions); i++ {
			baseQuery += " AND " + conditions[i]
		}
	}

	baseQuery += " ORDER BY date DESC, id DESC"

	rows, err := dm.db.Query(baseQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("could not search tasks: %v", err)
	}
	defer rows.Close()

	var tasks []workflow.Task
	for rows.Next() {
		var task workflow.Task
		var createdAtStr string

		err := rows.Scan(&task.ID, &task.Description, &task.Hours, &task.Category, &task.Date, &task.Status, &createdAtStr)
		if err != nil {
			return nil, fmt.Errorf("could not scan task: %v", err)
		}

		// Parsear created_at
		if createdAt, err := time.Parse("2006-01-02 15:04:05", createdAtStr); err == nil {
			task.CreatedAt = createdAt
		} else {
			task.CreatedAt = time.Now()
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// GetDatabasePath devuelve la ruta de la base de datos
func (dm *DatabaseManager) GetDatabasePath() string {
	return dm.dbPath
}
