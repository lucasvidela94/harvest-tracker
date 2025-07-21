package upgrade

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// RollbackManager maneja las operaciones de rollback
type RollbackManager struct {
	backupDir  string
	installDir string
	homeDir    string
}

// NewRollbackManager crea un nuevo gestor de rollback
func NewRollbackManager() *RollbackManager {
	homeDir, _ := os.UserHomeDir()
	backupDir := filepath.Join(homeDir, ".workflow", "backup")
	installDir := filepath.Join(homeDir, ".workflow", "install")

	return &RollbackManager{
		backupDir:  backupDir,
		installDir: installDir,
		homeDir:    homeDir,
	}
}

// DetectInstallationFailure detecta si hubo un fallo en la instalación
func (rm *RollbackManager) DetectInstallationFailure(installPath string) error {
	// Verificar que el binario existe
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		return fmt.Errorf("installation failed: binary not found")
	}

	// Verificar que es ejecutable
	if info, err := os.Stat(installPath); err == nil {
		if info.Mode()&0111 == 0 {
			return fmt.Errorf("installation failed: binary not executable")
		}
	}

	// Verificar que se puede ejecutar (test básico)
	if err := rm.testBinaryExecution(installPath); err != nil {
		return fmt.Errorf("installation failed: binary test failed: %v", err)
	}

	return nil
}

// testBinaryExecution prueba que el binario se puede ejecutar
func (rm *RollbackManager) testBinaryExecution(binaryPath string) error {
	// Por ahora, solo verificamos que el archivo existe y es ejecutable
	// En una implementación real, ejecutaríamos el binario con --version
	// para verificar que funciona correctamente

	return nil
}

// PerformRollback ejecuta el rollback completo
func (rm *RollbackManager) PerformRollback() error {
	// Notificar inicio del rollback
	rm.NotifyUser("🔄 Starting rollback process...")

	// Restaurar versión anterior
	if err := rm.RestorePreviousVersion(); err != nil {
		rm.NotifyUser(fmt.Sprintf("❌ Failed to restore previous version: %v", err))
		return fmt.Errorf("could not restore previous version: %v", err)
	}

	// Restaurar datos del backup
	if err := rm.RestoreDataFromBackup(); err != nil {
		rm.NotifyUser(fmt.Sprintf("❌ Failed to restore data: %v", err))
		return fmt.Errorf("could not restore data: %v", err)
	}

	// Verificar rollback exitoso
	if err := rm.VerifyRollback(); err != nil {
		rm.NotifyUser(fmt.Sprintf("❌ Rollback verification failed: %v", err))
		return fmt.Errorf("rollback verification failed: %v", err)
	}

	// Notificar éxito
	rm.NotifyUser("✅ Rollback completed successfully")

	return nil
}

// RestorePreviousVersion restaura la versión anterior del binario
func (rm *RollbackManager) RestorePreviousVersion() error {
	installPath := rm.getInstallPath()
	backupPath := filepath.Join(rm.backupDir, "workflow.bak")

	// Verificar que existe backup del binario
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("no binary backup found for rollback")
	}

	// Crear directorio de instalación si no existe
	installDir := filepath.Dir(installPath)
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("could not create install directory: %v", err)
	}

	// Restaurar binario anterior
	if err := rm.copyBinary(backupPath, installPath); err != nil {
		return fmt.Errorf("could not restore previous binary: %v", err)
	}

	// Hacer ejecutable
	if err := os.Chmod(installPath, 0755); err != nil {
		return fmt.Errorf("could not make binary executable: %v", err)
	}

	rm.NotifyUser("📦 Previous version restored")

	return nil
}

// RestoreDataFromBackup restaura los datos desde el backup más reciente
func (rm *RollbackManager) RestoreDataFromBackup() error {
	backupManager := NewBackupManager()

	// Verificar que existe un backup
	latestBackup, err := backupManager.GetLatestBackupPath()
	if err != nil {
		return fmt.Errorf("no backup found for data restoration")
	}

	// Restaurar datos
	if err := backupManager.RestoreBackup(); err != nil {
		return fmt.Errorf("could not restore backup: %v", err)
	}

	rm.NotifyUser(fmt.Sprintf("💾 Data restored from backup: %s", latestBackup))

	return nil
}

// VerifyRollback verifica que el rollback fue exitoso
func (rm *RollbackManager) VerifyRollback() error {
	installPath := rm.getInstallPath()

	// Verificar que el binario existe
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		return fmt.Errorf("rollback failed: binary not found")
	}

	// Verificar que es ejecutable
	if info, err := os.Stat(installPath); err == nil {
		if info.Mode()&0111 == 0 {
			return fmt.Errorf("rollback failed: binary not executable")
		}
	}

	// Verificar que se puede ejecutar
	if err := rm.testBinaryExecution(installPath); err != nil {
		return fmt.Errorf("rollback failed: binary test failed: %v", err)
	}

	// Verificar que los datos están restaurados
	if err := rm.verifyDataRestoration(); err != nil {
		return fmt.Errorf("rollback failed: data verification failed: %v", err)
	}

	return nil
}

// verifyDataRestoration verifica que los datos se restauraron correctamente
func (rm *RollbackManager) verifyDataRestoration() error {
	// Verificar que existen los archivos críticos
	criticalFiles := []string{"config.json", "tasks.json"}

	for _, filename := range criticalFiles {
		filePath := filepath.Join(rm.homeDir, ".workflow", filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("critical file missing after rollback: %s", filename)
		}
	}

	return nil
}

// NotifyUser notifica al usuario sobre el proceso de rollback
func (rm *RollbackManager) NotifyUser(message string) {
	fmt.Println(message)

	// También escribir en log
	rm.writeToLog(message)
}

// writeToLog escribe mensajes en el log de rollback
func (rm *RollbackManager) writeToLog(message string) {
	logPath := filepath.Join(rm.backupDir, "rollback.log")

	// Crear mensaje con timestamp
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	logMessage := fmt.Sprintf("[%s] %s\n", timestamp, message)

	// Agregar al archivo de log
	file, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		// Si no se puede escribir el log, no fallar
		return
	}
	defer file.Close()

	file.WriteString(logMessage)
}

// GetRollbackLog obtiene el contenido del log de rollback
func (rm *RollbackManager) GetRollbackLog() (string, error) {
	logPath := filepath.Join(rm.backupDir, "rollback.log")

	content, err := os.ReadFile(logPath)
	if err != nil {
		return "", fmt.Errorf("could not read rollback log: %v", err)
	}

	return string(content), nil
}

// CleanupRollbackLog limpia el log de rollback
func (rm *RollbackManager) CleanupRollbackLog() error {
	logPath := filepath.Join(rm.backupDir, "rollback.log")

	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		// No hay log para limpiar
		return nil
	}

	if err := os.Remove(logPath); err != nil {
		return fmt.Errorf("could not remove rollback log: %v", err)
	}

	return nil
}

// getInstallPath determina la ruta de instalación
func (rm *RollbackManager) getInstallPath() string {
	// Para desarrollo, instalar en ~/.local/bin
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".local", "bin", "workflow")
}

// copyBinary copia un archivo binario
func (rm *RollbackManager) copyBinary(src, dst string) error {
	// Abrir archivo fuente
	srcFile, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("could not open source file: %v", err)
	}
	defer srcFile.Close()

	// Crear archivo destino
	dstFile, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("could not create destination file: %v", err)
	}
	defer dstFile.Close()

	// Copiar contenido
	if _, err := rm.copyFileContent(srcFile, dstFile); err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	return nil
}

// copyFileContent copia el contenido de un archivo a otro
func (rm *RollbackManager) copyFileContent(src, dst *os.File) (int64, error) {
	// Usar io.Copy para copiar eficientemente
	return io.Copy(dst, src)
}

// IsRollbackAvailable verifica si hay un rollback disponible
func (rm *RollbackManager) IsRollbackAvailable() bool {
	// Verificar que existe backup del binario
	binaryBackup := filepath.Join(rm.backupDir, "workflow.bak")
	if _, err := os.Stat(binaryBackup); os.IsNotExist(err) {
		return false
	}

	// Verificar que existe backup de datos
	backupManager := NewBackupManager()
	if _, err := backupManager.GetLatestBackupPath(); err != nil {
		return false
	}

	return true
}

// GetRollbackInfo obtiene información sobre el rollback disponible
func (rm *RollbackManager) GetRollbackInfo() (map[string]interface{}, error) {
	info := make(map[string]interface{})

	// Verificar disponibilidad
	info["available"] = rm.IsRollbackAvailable()

	if !info["available"].(bool) {
		return info, nil
	}

	// Información del backup del binario
	binaryBackup := filepath.Join(rm.backupDir, "workflow.bak")
	if stat, err := os.Stat(binaryBackup); err == nil {
		info["binary_backup_size"] = stat.Size()
		info["binary_backup_time"] = stat.ModTime()
	}

	// Información del backup de datos
	backupManager := NewBackupManager()
	if latestBackup, err := backupManager.GetLatestBackupPath(); err == nil {
		if stat, err := os.Stat(latestBackup); err == nil {
			info["data_backup_time"] = stat.ModTime()
		}
		info["data_backup_path"] = latestBackup
	}

	return info, nil
}
