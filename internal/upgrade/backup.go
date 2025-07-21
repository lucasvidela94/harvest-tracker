package upgrade

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// BackupManager maneja las operaciones de backup
type BackupManager struct {
	backupDir string
	homeDir   string
}

// NewBackupManager crea un nuevo gestor de backup
func NewBackupManager() *BackupManager {
	homeDir, _ := os.UserHomeDir()
	backupDir := filepath.Join(homeDir, ".workflow", "backup")

	return &BackupManager{
		backupDir: backupDir,
		homeDir:   homeDir,
	}
}

// CreateBackup crea un backup completo de los datos del usuario
func (bm *BackupManager) CreateBackup() error {
	// Crear directorio de backup si no existe
	if err := os.MkdirAll(bm.backupDir, 0755); err != nil {
		return fmt.Errorf("could not create backup directory: %v", err)
	}

	// Crear timestamp para el backup
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := filepath.Join(bm.backupDir, fmt.Sprintf("backup_%s", timestamp))

	if err := os.MkdirAll(backupPath, 0755); err != nil {
		return fmt.Errorf("could not create backup path: %v", err)
	}

	// Lista de archivos a hacer backup
	filesToBackup := []string{
		"config.json",
		"tasks.json",
	}

	// Crear backup de cada archivo
	for _, filename := range filesToBackup {
		sourcePath := filepath.Join(bm.homeDir, ".workflow", filename)
		destPath := filepath.Join(backupPath, filename)

		if err := bm.backupFile(sourcePath, destPath); err != nil {
			return fmt.Errorf("could not backup %s: %v", filename, err)
		}
	}

	// Crear archivo de metadata del backup
	if err := bm.createBackupMetadata(backupPath, timestamp); err != nil {
		return fmt.Errorf("could not create backup metadata: %v", err)
	}

	// Crear symlink al backup más reciente
	latestLink := filepath.Join(bm.backupDir, "latest")
	os.Remove(latestLink) // Remover link anterior si existe
	if err := os.Symlink(backupPath, latestLink); err != nil {
		// Si symlink falla, crear un archivo con la ruta
		latestFile := filepath.Join(bm.backupDir, "latest.txt")
		if err := os.WriteFile(latestFile, []byte(backupPath), 0644); err != nil {
			return fmt.Errorf("could not create latest backup reference: %v", err)
		}
	}

	return nil
}

// backupFile hace backup de un archivo individual
func (bm *BackupManager) backupFile(source, dest string) error {
	// Verificar si el archivo fuente existe
	if _, err := os.Stat(source); os.IsNotExist(err) {
		// Si no existe, crear un archivo vacío en el backup
		return os.WriteFile(dest, []byte("{}"), 0644)
	}

	// Abrir archivo fuente
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("could not open source file: %v", err)
	}
	defer sourceFile.Close()

	// Crear archivo destino
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create destination file: %v", err)
	}
	defer destFile.Close()

	// Copiar contenido
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	return nil
}

// createBackupMetadata crea un archivo de metadata del backup
func (bm *BackupManager) createBackupMetadata(backupPath, timestamp string) error {
	metadata := fmt.Sprintf(`{
  "timestamp": "%s",
  "backup_path": "%s",
  "files": [
    "config.json",
    "tasks.json"
  ],
  "version": "1.0",
  "created_by": "workflow-cli-go"
}`, timestamp, backupPath)

	metadataPath := filepath.Join(backupPath, "backup.json")
	return os.WriteFile(metadataPath, []byte(metadata), 0644)
}

// RestoreBackup restaura los datos desde el backup más reciente
func (bm *BackupManager) RestoreBackup() error {
	// Obtener ruta del backup más reciente
	latestBackup, err := bm.GetLatestBackupPath()
	if err != nil {
		return fmt.Errorf("could not find latest backup: %v", err)
	}

	// Lista de archivos a restaurar
	filesToRestore := []string{
		"config.json",
		"tasks.json",
	}

	// Restaurar cada archivo
	for _, filename := range filesToRestore {
		sourcePath := filepath.Join(latestBackup, filename)
		destPath := filepath.Join(bm.homeDir, ".workflow", filename)

		if err := bm.restoreFile(sourcePath, destPath); err != nil {
			return fmt.Errorf("could not restore %s: %v", filename, err)
		}
	}

	return nil
}

// restoreFile restaura un archivo individual
func (bm *BackupManager) restoreFile(source, dest string) error {
	// Verificar si el archivo fuente existe
	if _, err := os.Stat(source); os.IsNotExist(err) {
		// Si no existe en el backup, no restaurar
		return nil
	}

	// Abrir archivo fuente
	sourceFile, err := os.Open(source)
	if err != nil {
		return fmt.Errorf("could not open source file: %v", err)
	}
	defer sourceFile.Close()

	// Crear directorio destino si no existe
	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return fmt.Errorf("could not create destination directory: %v", err)
	}

	// Crear archivo destino
	destFile, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("could not create destination file: %v", err)
	}
	defer destFile.Close()

	// Copiar contenido
	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	return nil
}

// GetLatestBackupPath obtiene la ruta del backup más reciente
func (bm *BackupManager) GetLatestBackupPath() (string, error) {
	// Intentar leer symlink
	latestLink := filepath.Join(bm.backupDir, "latest")
	if linkPath, err := os.Readlink(latestLink); err == nil {
		return linkPath, nil
	}

	// Si no hay symlink, intentar leer archivo de texto
	latestFile := filepath.Join(bm.backupDir, "latest.txt")
	if content, err := os.ReadFile(latestFile); err == nil {
		return string(content), nil
	}

	// Si no hay referencia, buscar el backup más reciente por timestamp
	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return "", fmt.Errorf("could not read backup directory: %v", err)
	}

	var latestBackup string
	var latestTime time.Time

	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) > 7 && entry.Name()[:7] == "backup_" {
			// Extraer timestamp del nombre
			timestampStr := entry.Name()[7:] // Remover "backup_"
			if backupTime, err := time.Parse("2006-01-02_15-04-05", timestampStr); err == nil {
				if backupTime.After(latestTime) {
					latestTime = backupTime
					latestBackup = filepath.Join(bm.backupDir, entry.Name())
				}
			}
		}
	}

	if latestBackup == "" {
		return "", fmt.Errorf("no backup found")
	}

	return latestBackup, nil
}

// VerifyBackup verifica la integridad del backup
func (bm *BackupManager) VerifyBackup() error {
	latestBackup, err := bm.GetLatestBackupPath()
	if err != nil {
		return fmt.Errorf("could not find backup to verify: %v", err)
	}

	// Verificar que el directorio existe
	if _, err := os.Stat(latestBackup); os.IsNotExist(err) {
		return fmt.Errorf("backup directory does not exist: %s", latestBackup)
	}

	// Verificar archivos críticos
	criticalFiles := []string{"config.json", "tasks.json"}
	for _, filename := range criticalFiles {
		filePath := filepath.Join(latestBackup, filename)
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			return fmt.Errorf("critical file missing in backup: %s", filename)
		}
	}

	// Verificar metadata
	metadataPath := filepath.Join(latestBackup, "backup.json")
	if _, err := os.Stat(metadataPath); os.IsNotExist(err) {
		return fmt.Errorf("backup metadata missing")
	}

	return nil
}

// GetBackupPath devuelve la ruta del directorio de backup
func (bm *BackupManager) GetBackupPath() string {
	return bm.backupDir
}

// ListBackups lista todos los backups disponibles
func (bm *BackupManager) ListBackups() ([]string, error) {
	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return nil, fmt.Errorf("could not read backup directory: %v", err)
	}

	var backups []string
	for _, entry := range entries {
		if entry.IsDir() && len(entry.Name()) > 7 && entry.Name()[:7] == "backup_" {
			backups = append(backups, entry.Name())
		}
	}

	return backups, nil
}

// CleanOldBackups limpia backups antiguos (mantiene solo los últimos 5)
func (bm *BackupManager) CleanOldBackups() error {
	backups, err := bm.ListBackups()
	if err != nil {
		return fmt.Errorf("could not list backups: %v", err)
	}

	// Si hay menos de 5 backups, no hacer nada
	if len(backups) <= 5 {
		return nil
	}

	// Ordenar backups por timestamp (más antiguos primero)
	// Por simplicidad, asumimos que los nombres están ordenados cronológicamente
	// En una implementación real, deberíamos parsear los timestamps

	// Eliminar backups antiguos (mantener solo los últimos 5)
	for i := 0; i < len(backups)-5; i++ {
		oldBackupPath := filepath.Join(bm.backupDir, backups[i])
		if err := os.RemoveAll(oldBackupPath); err != nil {
			return fmt.Errorf("could not remove old backup %s: %v", backups[i], err)
		}
	}

	return nil
}
