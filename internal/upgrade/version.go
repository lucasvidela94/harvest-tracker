package upgrade

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/google/go-github/v62/github"
)

const (
	// RepoOwner es el propietario del repositorio
	RepoOwner = "lucasvidela94"
	// RepoName es el nombre del repositorio
	RepoName = "workflow-cli"
	// CurrentVersion es la versión actual de Go
	CurrentVersion = "1.0.1"
)

// VersionManager maneja la detección y comparación de versiones
type VersionManager struct {
	client *github.Client
}

// NewVersionManager crea un nuevo gestor de versiones
func NewVersionManager() *VersionManager {
	return &VersionManager{
		client: github.NewClient(nil),
	}
}

// DetectPythonInstallation detecta si existe una instalación de Python
func (vm *VersionManager) DetectPythonInstallation() (bool, string, error) {
	// Buscar el script de Python en el PATH
	pythonPath, err := exec.LookPath("workflow")
	if err != nil {
		return false, "", nil // No hay instalación Python
	}

	// Verificar que es el script de Python (no el binario de Go)
	if strings.Contains(pythonPath, "scripts/workflow") || strings.HasSuffix(pythonPath, "workflow") {
		// Verificar si es un script Python (no un binario)
		info, err := os.Stat(pythonPath)
		if err != nil {
			return false, "", nil
		}

		// Si es un archivo regular (no un binario ejecutable), probablemente es Python
		if info.Mode().IsRegular() && (info.Mode()&0111) == 0 {
			// Intentar ejecutar para obtener versión
			cmd := exec.Command("python3", pythonPath, "--version")
			output, err := cmd.Output()
			if err != nil {
				// Si falla, intentar con python
				cmd = exec.Command("python", pythonPath, "--version")
				output, err = cmd.Output()
				if err != nil {
					return false, "", nil // No se puede ejecutar, asumir que no es Python válido
				}
			}

			// Extraer versión del output
			version := strings.TrimSpace(string(output))
			return true, version, nil
		}
	}

	return false, "", nil
}

// GetCurrentVersion obtiene la versión actual
func (vm *VersionManager) GetCurrentVersion() (string, error) {
	// Primero intentar detectar versión Python
	hasPython, pythonVersion, err := vm.DetectPythonInstallation()
	if err != nil {
		return "", err
	}

	if hasPython {
		return pythonVersion, nil
	}

	// Si no hay Python, usar versión Go
	return CurrentVersion, nil
}

// GetLatestVersion obtiene la última versión disponible en GitHub
func (vm *VersionManager) GetLatestVersion() (string, error) {
	ctx := context.Background()

	// Obtener el último release
	release, _, err := vm.client.Repositories.GetLatestRelease(ctx, RepoOwner, RepoName)
	if err != nil {
		return "", fmt.Errorf("could not get latest release: %v", err)
	}

	if release.TagName == nil {
		return "", fmt.Errorf("no tag name found in latest release")
	}

	// Limpiar el tag (remover 'v' si existe)
	version := strings.TrimPrefix(*release.TagName, "v")
	return version, nil
}

// CompareVersions compara dos versiones y determina si hay actualización disponible
func (vm *VersionManager) CompareVersions(current, latest string) (bool, error) {
	// Si la versión actual es "unknown" o vacía, considerar que hay actualización
	if current == "" || current == "unknown" {
		return true, nil
	}

	// Comparación simple de versiones semánticas
	currentClean := strings.TrimPrefix(current, "v")
	latestClean := strings.TrimPrefix(latest, "v")

	// Si las versiones son iguales, no hay actualización
	if currentClean == latestClean {
		return false, nil
	}

	// Para una comparación más robusta, podríamos usar una librería de versiones
	// Por ahora, asumimos que si son diferentes, hay actualización disponible
	return true, nil
}

// GetBinaryName obtiene el nombre del binario para la plataforma actual
func (vm *VersionManager) GetBinaryName() string {
	os := runtime.GOOS
	arch := runtime.GOARCH

	switch os {
	case "linux":
		if arch == "amd64" {
			return "workflow-linux-amd64"
		} else if arch == "arm64" {
			return "workflow-linux-arm64"
		}
	case "darwin":
		if arch == "amd64" {
			return "workflow-darwin-amd64"
		} else if arch == "arm64" {
			return "workflow-darwin-arm64"
		}
	case "windows":
		if arch == "amd64" {
			return "workflow-windows-amd64.exe"
		}
	}

	// Fallback
	return "workflow"
}

// GetInstallationPath obtiene la ruta de instalación actual
func (vm *VersionManager) GetInstallationPath() (string, error) {
	// Buscar en el PATH
	path, err := exec.LookPath("workflow")
	if err != nil {
		return "", fmt.Errorf("workflow not found in PATH: %v", err)
	}

	return path, nil
}

// IsGoInstallation verifica si la instalación actual es de Go
func (vm *VersionManager) IsGoInstallation() (bool, error) {
	path, err := vm.GetInstallationPath()
	if err != nil {
		return false, err
	}

	// Verificar si es un binario ejecutable (no un script)
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}

	// Si es un archivo ejecutable, probablemente es Go
	return info.Mode().IsRegular() && (info.Mode()&0111) != 0, nil
}

// GetUpgradeInfo obtiene información completa sobre el upgrade
func (vm *VersionManager) GetUpgradeInfo() (*UpgradeInfo, error) {
	currentVersion, err := vm.GetCurrentVersion()
	if err != nil {
		return nil, fmt.Errorf("could not get current version: %v", err)
	}

	latestVersion, err := vm.GetLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("could not get latest version: %v", err)
	}

	hasUpdate, err := vm.CompareVersions(currentVersion, latestVersion)
	if err != nil {
		return nil, fmt.Errorf("could not compare versions: %v", err)
	}

	isGo, err := vm.IsGoInstallation()
	if err != nil {
		isGo = false // Asumir que no es Go si hay error
	}

	return &UpgradeInfo{
		CurrentVersion:   currentVersion,
		LatestVersion:    latestVersion,
		HasUpdate:        hasUpdate,
		IsGoInstallation: isGo,
		BinaryName:       vm.GetBinaryName(),
	}, nil
}

// UpgradeInfo contiene información sobre el upgrade
type UpgradeInfo struct {
	CurrentVersion   string
	LatestVersion    string
	HasUpdate        bool
	IsGoInstallation bool
	BinaryName       string
}
