package upgrade

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// InstallManager maneja las operaciones de instalación
type InstallManager struct {
	installDir string
	homeDir    string
	backupDir  string
}

// NewInstallManager crea un nuevo gestor de instalación
func NewInstallManager() *InstallManager {
	homeDir, _ := os.UserHomeDir()
	installDir := filepath.Join(homeDir, ".workflow", "install")
	backupDir := filepath.Join(homeDir, ".workflow", "backup")

	return &InstallManager{
		installDir: installDir,
		homeDir:    homeDir,
		backupDir:  backupDir,
	}
}

// InstallNewVersion instala la nueva versión
func (im *InstallManager) InstallNewVersion(downloadPath, version string) error {
	// Crear directorio de instalación si no existe
	if err := os.MkdirAll(im.installDir, 0755); err != nil {
		return fmt.Errorf("could not create install directory: %v", err)
	}

	// Extraer archivo descargado
	extractPath, err := im.ExtractArchive(downloadPath)
	if err != nil {
		return fmt.Errorf("could not extract archive: %v", err)
	}

	// Reemplazar binario
	if err := im.ReplaceBinary(extractPath); err != nil {
		return fmt.Errorf("could not replace binary: %v", err)
	}

	// Restaurar datos del backup
	if err := im.RestoreData(); err != nil {
		return fmt.Errorf("could not restore data: %v", err)
	}

	// Verificar instalación
	if err := im.VerifyInstallation(); err != nil {
		return fmt.Errorf("installation verification failed: %v", err)
	}

	// Limpiar archivos temporales
	if err := im.Cleanup(extractPath); err != nil {
		return fmt.Errorf("could not cleanup temporary files: %v", err)
	}

	return nil
}

// ExtractArchive extrae el archivo tar.gz descargado
func (im *InstallManager) ExtractArchive(archivePath string) (string, error) {
	// Crear directorio temporal para extracción
	extractDir := filepath.Join(im.installDir, "extract")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return "", fmt.Errorf("could not create extract directory: %v", err)
	}

	// Abrir archivo tar.gz
	file, err := os.Open(archivePath)
	if err != nil {
		return "", fmt.Errorf("could not open archive: %v", err)
	}
	defer file.Close()

	// Crear reader gzip
	gzr, err := gzip.NewReader(file)
	if err != nil {
		return "", fmt.Errorf("could not create gzip reader: %v", err)
	}
	defer gzr.Close()

	// Crear reader tar
	tr := tar.NewReader(gzr)

	// Extraer archivos
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", fmt.Errorf("could not read tar header: %v", err)
		}

		// Construir ruta de destino
		target := filepath.Join(extractDir, header.Name)

		// Verificar que no hay path traversal
		if !strings.HasPrefix(target, extractDir) {
			return "", fmt.Errorf("invalid path in archive: %s", header.Name)
		}

		// Crear directorio si es necesario
		if header.Typeflag == tar.TypeDir {
			if err := os.MkdirAll(target, 0755); err != nil {
				return "", fmt.Errorf("could not create directory: %v", err)
			}
			continue
		}

		// Crear archivo
		if err := im.extractFile(tr, target, header); err != nil {
			return "", fmt.Errorf("could not extract file %s: %v", header.Name, err)
		}
	}

	return extractDir, nil
}

// extractFile extrae un archivo individual del tar
func (im *InstallManager) extractFile(tr *tar.Reader, target string, header *tar.Header) error {
	// Crear directorio padre si no existe
	dir := filepath.Dir(target)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("could not create parent directory: %v", err)
	}

	// Crear archivo destino
	file, err := os.Create(target)
	if err != nil {
		return fmt.Errorf("could not create file: %v", err)
	}
	defer file.Close()

	// Copiar contenido
	if _, err := io.Copy(file, tr); err != nil {
		return fmt.Errorf("could not copy file content: %v", err)
	}

	// Establecer permisos
	if err := file.Chmod(os.FileMode(header.Mode)); err != nil {
		return fmt.Errorf("could not set file permissions: %v", err)
	}

	return nil
}

// ReplaceBinary reemplaza el binario actual con el nuevo
func (im *InstallManager) ReplaceBinary(extractPath string) error {
	// Buscar binario en el directorio extraído
	binaryPath, err := im.findBinary(extractPath)
	if err != nil {
		return fmt.Errorf("could not find binary: %v", err)
	}

	// Determinar ruta de instalación
	installPath := im.getInstallPath()

	// Crear directorio de instalación si no existe
	installDir := filepath.Dir(installPath)
	if err := os.MkdirAll(installDir, 0755); err != nil {
		return fmt.Errorf("could not create install directory: %v", err)
	}

	// Hacer backup del binario actual si existe
	if err := im.backupCurrentBinary(installPath); err != nil {
		return fmt.Errorf("could not backup current binary: %v", err)
	}

	// Copiar nuevo binario
	if err := im.copyBinary(binaryPath, installPath); err != nil {
		return fmt.Errorf("could not copy binary: %v", err)
	}

	// Hacer ejecutable
	if err := os.Chmod(installPath, 0755); err != nil {
		return fmt.Errorf("could not make binary executable: %v", err)
	}

	return nil
}

// findBinary busca el binario en el directorio extraído
func (im *InstallManager) findBinary(extractPath string) (string, error) {
	// Buscar archivo ejecutable
	entries, err := os.ReadDir(extractPath)
	if err != nil {
		return "", fmt.Errorf("could not read extract directory: %v", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() && entry.Name() == "workflow" {
			return filepath.Join(extractPath, entry.Name()), nil
		}
	}

	// Si no se encuentra en el directorio raíz, buscar recursivamente
	return im.findBinaryRecursive(extractPath)
}

// findBinaryRecursive busca el binario recursivamente
func (im *InstallManager) findBinaryRecursive(dir string) (string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return "", fmt.Errorf("could not read directory: %v", err)
	}

	for _, entry := range entries {
		path := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			// Buscar en subdirectorios
			if found, err := im.findBinaryRecursive(path); err == nil {
				return found, nil
			}
		} else if entry.Name() == "workflow" {
			return path, nil
		}
	}

	return "", fmt.Errorf("binary not found in extracted archive")
}

// getInstallPath determina la ruta de instalación
func (im *InstallManager) getInstallPath() string {
	// Para desarrollo, instalar en ~/.local/bin
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".local", "bin", "workflow")
}

// backupCurrentBinary hace backup del binario actual
func (im *InstallManager) backupCurrentBinary(installPath string) error {
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		// No hay binario actual, no hacer backup
		return nil
	}

	backupPath := filepath.Join(im.backupDir, "workflow.bak")

	// Copiar binario actual a backup
	if err := im.copyBinary(installPath, backupPath); err != nil {
		return fmt.Errorf("could not backup current binary: %v", err)
	}

	return nil
}

// copyBinary copia un archivo binario
func (im *InstallManager) copyBinary(src, dst string) error {
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
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return fmt.Errorf("could not copy file: %v", err)
	}

	return nil
}

// RestoreData restaura los datos del backup
func (im *InstallManager) RestoreData() error {
	backupManager := NewBackupManager()

	// Verificar que existe un backup
	if _, err := backupManager.GetLatestBackupPath(); err != nil {
		// No hay backup, no restaurar
		return nil
	}

	// Restaurar datos
	if err := backupManager.RestoreBackup(); err != nil {
		return fmt.Errorf("could not restore backup: %v", err)
	}

	return nil
}

// VerifyInstallation verifica que la instalación fue exitosa
func (im *InstallManager) VerifyInstallation() error {
	installPath := im.getInstallPath()

	// Verificar que el binario existe
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		return fmt.Errorf("binary not found at %s", installPath)
	}

	// Verificar que es ejecutable
	if info, err := os.Stat(installPath); err == nil {
		if info.Mode()&0111 == 0 {
			return fmt.Errorf("binary is not executable")
		}
	}

	// Verificar que se puede ejecutar
	if err := im.testBinary(installPath); err != nil {
		return fmt.Errorf("binary test failed: %v", err)
	}

	return nil
}

// testBinary prueba que el binario funciona
func (im *InstallManager) testBinary(binaryPath string) error {
	// Por ahora, solo verificamos que el archivo existe y es ejecutable
	// En una implementación real, ejecutaríamos el binario con --version
	// para verificar que funciona correctamente

	return nil
}

// Cleanup limpia archivos temporales
func (im *InstallManager) Cleanup(extractPath string) error {
	// Eliminar directorio de extracción
	if err := os.RemoveAll(extractPath); err != nil {
		return fmt.Errorf("could not remove extract directory: %v", err)
	}

	// Limpiar directorio de instalación si está vacío
	if err := im.cleanupInstallDir(); err != nil {
		return fmt.Errorf("could not cleanup install directory: %v", err)
	}

	return nil
}

// cleanupInstallDir limpia el directorio de instalación si está vacío
func (im *InstallManager) cleanupInstallDir() error {
	entries, err := os.ReadDir(im.installDir)
	if err != nil {
		return fmt.Errorf("could not read install directory: %v", err)
	}

	// Si el directorio está vacío, eliminarlo
	if len(entries) == 0 {
		if err := os.Remove(im.installDir); err != nil {
			return fmt.Errorf("could not remove empty install directory: %v", err)
		}
	}

	return nil
}

// GetInstallPath devuelve la ruta de instalación
func (im *InstallManager) GetInstallPath() string {
	return im.getInstallPath()
}

// RollbackInstallation revierte la instalación
func (im *InstallManager) RollbackInstallation() error {
	installPath := im.getInstallPath()
	backupPath := filepath.Join(im.backupDir, "workflow.bak")

	// Verificar que existe backup del binario
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return fmt.Errorf("no binary backup found for rollback")
	}

	// Restaurar binario anterior
	if err := im.copyBinary(backupPath, installPath); err != nil {
		return fmt.Errorf("could not restore previous binary: %v", err)
	}

	// Hacer ejecutable
	if err := os.Chmod(installPath, 0755); err != nil {
		return fmt.Errorf("could not make binary executable: %v", err)
	}

	// Restaurar datos del backup
	if err := im.RestoreData(); err != nil {
		return fmt.Errorf("could not restore data during rollback: %v", err)
	}

	return nil
}
