package upgrade

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// GitHubRelease representa la respuesta de la API de GitHub
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// CheckForUpdates verifica si hay una nueva versión disponible
func CheckForUpdates(currentVersion string) (bool, string, error) {
	// Obtener la última versión desde GitHub
	latestVersion, err := getLatestVersion()
	if err != nil {
		return false, "", fmt.Errorf("error obteniendo última versión: %w", err)
	}

	// Comparar versiones
	if latestVersion != currentVersion {
		return true, latestVersion, nil
	}

	return false, "", nil
}

// AutoUpdate descarga e instala la última versión automáticamente
func AutoUpdate(currentVersion string) error {
	fmt.Println("🔄 Verificando actualizaciones...")

	// Verificar si hay actualizaciones
	hasUpdate, latestVersion, err := CheckForUpdates(currentVersion)
	if err != nil {
		return err
	}

	if !hasUpdate {
		fmt.Println("✅ Ya tienes la última versión:", currentVersion)
		return nil
	}

	fmt.Printf("📦 Nueva versión disponible: %s\n", latestVersion)
	fmt.Println("🚀 Descargando actualización...")

	// Descargar la nueva versión
	downloadURL, err := getDownloadURL(latestVersion)
	if err != nil {
		return fmt.Errorf("error obteniendo URL de descarga: %w", err)
	}

	// Crear directorio temporal
	tempDir, err := os.MkdirTemp("", "harvest-update-*")
	if err != nil {
		return fmt.Errorf("error creando directorio temporal: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Descargar archivo
	archivePath := filepath.Join(tempDir, "harvest-update.tar.gz")
	if err := downloadFile(downloadURL, archivePath); err != nil {
		return fmt.Errorf("error descargando archivo: %w", err)
	}

	// Extraer archivo
	if err := extractArchive(archivePath, tempDir); err != nil {
		return fmt.Errorf("error extrayendo archivo: %w", err)
	}

	// Obtener ruta del ejecutable actual
	currentExe, err := os.Executable()
	if err != nil {
		return fmt.Errorf("error obteniendo ruta del ejecutable: %w", err)
	}

	// Crear backup
	backupPath := currentExe + ".backup"
	if err := copyFile(currentExe, backupPath); err != nil {
		return fmt.Errorf("error creando backup: %w", err)
	}

	// Encontrar el nuevo ejecutable
	newExePath := filepath.Join(tempDir, "harvest")
	if runtime.GOOS == "windows" {
		newExePath = filepath.Join(tempDir, "harvest.exe")
	}

	// Reemplazar ejecutable
	if err := copyFile(newExePath, currentExe); err != nil {
		// Restaurar backup en caso de error
		copyFile(backupPath, currentExe)
		return fmt.Errorf("error reemplazando ejecutable: %w", err)
	}

	// Hacer ejecutable
	if err := os.Chmod(currentExe, 0755); err != nil {
		return fmt.Errorf("error estableciendo permisos: %w", err)
	}

	// Limpiar backup
	os.Remove(backupPath)

	fmt.Printf("✅ ¡Actualizado exitosamente a la versión %s!\n", latestVersion)
	return nil
}

// getLatestVersion obtiene la última versión desde GitHub
func getLatestVersion() (string, error) {
	resp, err := http.Get("https://api.github.com/repos/lucasvidela94/harvest-tracker/releases/latest")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", err
	}

	return release.TagName, nil
}

// getDownloadURL obtiene la URL de descarga para la plataforma actual
func getDownloadURL(version string) (string, error) {
	resp, err := http.Get("https://api.github.com/repos/lucasvidela94/harvest-tracker/releases/tags/" + version)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var release GitHubRelease
	if err := json.Unmarshal(body, &release); err != nil {
		return "", err
	}

	// Construir nombre del archivo esperado
	os := runtime.GOOS
	arch := runtime.GOARCH
	if arch == "amd64" {
		arch = "amd64"
	} else if arch == "arm64" {
		arch = "arm64"
	}

	expectedName := fmt.Sprintf("harvest-%s-%s-%s.tar.gz", version, os, arch)

	// Buscar el asset correcto
	for _, asset := range release.Assets {
		if asset.Name == expectedName {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("no se encontró el archivo para %s-%s", os, arch)
}

// downloadFile descarga un archivo desde una URL
func downloadFile(url, filepath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// extractArchive extrae un archivo tar.gz
func extractArchive(archivePath, destDir string) error {
	cmd := exec.Command("tar", "-xzf", archivePath, "-C", destDir)
	return cmd.Run()
}

// copyFile copia un archivo
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
