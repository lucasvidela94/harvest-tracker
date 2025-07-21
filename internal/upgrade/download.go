package upgrade

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// DownloadManager maneja las operaciones de descarga
type DownloadManager struct {
	downloadDir string
	homeDir     string
	client      *http.Client
}

// NewDownloadManager crea un nuevo gestor de descarga
func NewDownloadManager() *DownloadManager {
	homeDir, _ := os.UserHomeDir()
	downloadDir := filepath.Join(homeDir, ".workflow", "downloads")

	// Crear cliente HTTP con timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	return &DownloadManager{
		downloadDir: downloadDir,
		homeDir:     homeDir,
		client:      client,
	}
}

// DownloadLatestVersion descarga la última versión disponible
func (dm *DownloadManager) DownloadLatestVersion(version string) (string, error) {
	// Crear directorio de descargas si no existe
	if err := os.MkdirAll(dm.downloadDir, 0755); err != nil {
		return "", fmt.Errorf("could not create download directory: %v", err)
	}

	// Obtener URL de descarga
	downloadURL, err := dm.GetDownloadURL(version)
	if err != nil {
		return "", fmt.Errorf("could not get download URL: %v", err)
	}

	// Obtener nombre del archivo
	filename := dm.GetFilename(version)
	downloadPath := filepath.Join(dm.downloadDir, filename)

	// Descargar archivo
	if err := dm.downloadFile(downloadURL, downloadPath); err != nil {
		return "", fmt.Errorf("could not download file: %v", err)
	}

	// Verificar descarga
	if err := dm.VerifyDownload(downloadPath, version); err != nil {
		// Si la verificación falla, limpiar archivo descargado
		os.Remove(downloadPath)
		return "", fmt.Errorf("download verification failed: %v", err)
	}

	return downloadPath, nil
}

// GetDownloadURL construye la URL de descarga para la versión
func (dm *DownloadManager) GetDownloadURL(version string) (string, error) {
	// Determinar arquitectura y sistema operativo
	arch := runtime.GOARCH
	os := runtime.GOOS

	// Mapear arquitecturas
	archMap := map[string]string{
		"amd64": "x86_64",
		"arm64": "aarch64",
		"386":   "i386",
	}

	// Mapear sistemas operativos
	osMap := map[string]string{
		"linux":   "linux",
		"darwin":  "darwin",
		"windows": "windows",
	}

	// Obtener nombres mapeados
	archName, ok := archMap[arch]
	if !ok {
		return "", fmt.Errorf("unsupported architecture: %s", arch)
	}

	osName, ok := osMap[os]
	if !ok {
		return "", fmt.Errorf("unsupported operating system: %s", os)
	}

	// Construir URL de descarga
	// Formato: https://github.com/{owner}/{repo}/releases/download/v{version}/workflow-{version}-{os}-{arch}.tar.gz
	downloadURL := fmt.Sprintf(
		"https://github.com/%s/%s/releases/download/v%s/workflow-%s-%s-%s.tar.gz",
		RepoOwner,
		RepoName,
		version,
		version,
		osName,
		archName,
	)

	return downloadURL, nil
}

// GetFilename genera el nombre del archivo de descarga
func (dm *DownloadManager) GetFilename(version string) string {
	arch := runtime.GOARCH
	os := runtime.GOOS

	// Mapear arquitecturas
	archMap := map[string]string{
		"amd64": "x86_64",
		"arm64": "aarch64",
		"386":   "i386",
	}

	// Mapear sistemas operativos
	osMap := map[string]string{
		"linux":   "linux",
		"darwin":  "darwin",
		"windows": "windows",
	}

	archName := archMap[arch]
	osName := osMap[os]

	return fmt.Sprintf("workflow-%s-%s-%s.tar.gz", version, osName, archName)
}

// downloadFile descarga un archivo desde una URL
func (dm *DownloadManager) downloadFile(url, destPath string) error {
	// Crear request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("could not create request: %v", err)
	}

	// Agregar User-Agent para evitar bloqueos
	req.Header.Set("User-Agent", "workflow-cli-go/1.0")

	// Realizar request
	resp, err := dm.client.Do(req)
	if err != nil {
		return fmt.Errorf("could not download file: %v", err)
	}
	defer resp.Body.Close()

	// Verificar status code
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status: %d", resp.StatusCode)
	}

	// Crear archivo destino
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("could not create destination file: %v", err)
	}
	defer destFile.Close()

	// Copiar contenido con progreso
	written, err := io.Copy(destFile, resp.Body)
	if err != nil {
		return fmt.Errorf("could not write file: %v", err)
	}

	// Verificar que se escribió algo
	if written == 0 {
		return fmt.Errorf("downloaded file is empty")
	}

	return nil
}

// VerifyDownload verifica la integridad del archivo descargado
func (dm *DownloadManager) VerifyDownload(filePath, version string) error {
	// Verificar que el archivo existe
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("downloaded file does not exist: %s", filePath)
	}

	// Verificar tamaño mínimo (al menos 1MB)
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("could not get file info: %v", err)
	}

	if fileInfo.Size() < 1024*1024 { // 1MB
		return fmt.Errorf("downloaded file is too small: %d bytes", fileInfo.Size())
	}

	// Verificar que es un archivo tar.gz válido
	if !strings.HasSuffix(filePath, ".tar.gz") {
		return fmt.Errorf("downloaded file is not a valid tar.gz archive")
	}

	// Por ahora, solo verificamos el tamaño y extensión
	// En una implementación real, verificaríamos el checksum SHA256
	// que debería estar disponible en el release de GitHub

	return nil
}

// GetDownloadPath devuelve la ruta del directorio de descargas
func (dm *DownloadManager) GetDownloadPath() string {
	return dm.downloadDir
}

// CleanDownloads limpia archivos de descarga antiguos
func (dm *DownloadManager) CleanDownloads() error {
	entries, err := os.ReadDir(dm.downloadDir)
	if err != nil {
		return fmt.Errorf("could not read download directory: %v", err)
	}

	// Eliminar archivos más antiguos que 24 horas
	cutoff := time.Now().Add(-24 * time.Hour)

	for _, entry := range entries {
		if !entry.IsDir() {
			filePath := filepath.Join(dm.downloadDir, entry.Name())

			// Obtener información del archivo
			fileInfo, err := os.Stat(filePath)
			if err != nil {
				continue
			}

			// Si el archivo es más antiguo que 24 horas, eliminarlo
			if fileInfo.ModTime().Before(cutoff) {
				if err := os.Remove(filePath); err != nil {
					return fmt.Errorf("could not remove old download %s: %v", entry.Name(), err)
				}
			}
		}
	}

	return nil
}

// GetDownloadSize obtiene el tamaño del archivo a descargar
func (dm *DownloadManager) GetDownloadSize(url string) (int64, error) {
	// Crear request HEAD para obtener información del archivo
	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return 0, fmt.Errorf("could not create request: %v", err)
	}

	req.Header.Set("User-Agent", "workflow-cli-go/1.0")

	resp, err := dm.client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("could not get file info: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("could not get file info, status: %d", resp.StatusCode)
	}

	// Obtener tamaño del archivo
	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "" {
		return 0, fmt.Errorf("could not determine file size")
	}

	// Parsear tamaño
	var size int64
	if _, err := fmt.Sscanf(contentLength, "%d", &size); err != nil {
		return 0, fmt.Errorf("could not parse file size: %v", err)
	}

	return size, nil
}

// CalculateChecksum calcula el checksum SHA256 de un archivo
func (dm *DownloadManager) CalculateChecksum(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("could not open file: %v", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", fmt.Errorf("could not calculate checksum: %v", err)
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// VerifyChecksum verifica el checksum de un archivo
func (dm *DownloadManager) VerifyChecksum(filePath, expectedChecksum string) error {
	actualChecksum, err := dm.CalculateChecksum(filePath)
	if err != nil {
		return fmt.Errorf("could not calculate checksum: %v", err)
	}

	if actualChecksum != expectedChecksum {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}
