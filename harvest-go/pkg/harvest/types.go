package harvest

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
	CreatedAt   time.Time `json:"created_at"`
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
