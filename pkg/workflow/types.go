package workflow

import (
	"time"
)

// Task representa una tarea individual
type Task struct {
	ID          int       `json:"id"`
	Description string    `json:"description"`
	Hours       float64   `json:"hours"`
	Category    string    `json:"category"`
	Date        string    `json:"date"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// Estados de tareas
const (
	StatusPending    = "pending"
	StatusInProgress = "in_progress"
	StatusCompleted  = "completed"
	StatusPaused     = "paused"
)

// StatusIcon mapea estados a iconos
var StatusIcon = map[string]string{
	StatusPending:    "⏳",
	StatusInProgress: "🔄",
	StatusCompleted:  "✅",
	StatusPaused:     "⏸️",
}

// GetStatusIcon devuelve el icono para un estado
func GetStatusIcon(status string) string {
	if icon, exists := StatusIcon[status]; exists {
		return icon
	}
	return StatusIcon[StatusPending] // default
}

// Config representa la configuración del usuario
type Config struct {
	DailyHoursTarget  float64 `json:"daily_hours_target"`
	DailyStandupHours float64 `json:"daily_standup_hours"`
	DataFile          string  `json:"data_file"`
	UserName          string  `json:"user_name"`
	Company           string  `json:"company"`
	Timezone          string  `json:"timezone"`
}

// CategoryIcon mapea categorías a iconos
var CategoryIcon = map[string]string{
	"tech":     "💻",
	"meeting":  "🤝",
	"qa":       "🧪",
	"doc":      "📚",
	"planning": "📋",
	"research": "🔍",
	"review":   "👀",
	"deploy":   "🚀",
	"daily":    "📢",
	"general":  "📝",
}

// GetIcon devuelve el icono para una categoría
func GetIcon(category string) string {
	if icon, exists := CategoryIcon[category]; exists {
		return icon
	}
	return "📝" // default
}
