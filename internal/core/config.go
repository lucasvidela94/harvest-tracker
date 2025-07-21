package core

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/lucasvidela94/workflow-cli/pkg/workflow"
)

// ConfigManager maneja la configuración del usuario
type ConfigManager struct {
	configPath string
	config     *workflow.Config
}

// NewConfigManager crea un nuevo gestor de configuración
func NewConfigManager() *ConfigManager {
	homeDir, _ := os.UserHomeDir()
	configPath := filepath.Join(homeDir, ".workflow", "config.json")

	return &ConfigManager{
		configPath: configPath,
		config:     getDefaultConfig(homeDir),
	}
}

// getDefaultConfig devuelve la configuración por defecto
func getDefaultConfig(homeDir string) *workflow.Config {
	return &workflow.Config{
		DailyHoursTarget:  8.0,
		DailyStandupHours: 0.25,
		DataFile:          filepath.Join(homeDir, ".workflow", "tasks.json"),
		UserName:          os.Getenv("USER"),
		Company:           "",
		Timezone:          "UTC",
	}
}

// Load carga la configuración desde el archivo
func (cm *ConfigManager) Load() error {
	if _, err := os.Stat(cm.configPath); os.IsNotExist(err) {
		// Crear directorio si no existe
		dir := filepath.Dir(cm.configPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return err
		}

		// Guardar configuración por defecto
		return cm.Save()
	}

	file, err := os.Open(cm.configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewDecoder(file).Decode(cm.config)
}

// Save guarda la configuración en el archivo
func (cm *ConfigManager) Save() error {
	file, err := os.Create(cm.configPath)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(cm.config)
}

// Get devuelve la configuración actual
func (cm *ConfigManager) Get() *workflow.Config {
	return cm.config
}

// GetDataFile devuelve la ruta del archivo de datos
func (cm *ConfigManager) GetDataFile() string {
	return cm.config.DataFile
}

// GetDailyHoursTarget devuelve el objetivo de horas diarias
func (cm *ConfigManager) GetDailyHoursTarget() float64 {
	return cm.config.DailyHoursTarget
}

// GetDailyStandupHours devuelve las horas del daily standup
func (cm *ConfigManager) GetDailyStandupHours() float64 {
	return cm.config.DailyStandupHours
}
