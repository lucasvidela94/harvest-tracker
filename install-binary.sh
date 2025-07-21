#!/bin/bash

# Script de instalación para workflow CLI usando binario pre-compilado
# No requiere Go instalado

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para imprimir mensajes
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Función para detectar el sistema operativo
detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux";;
        Darwin*)    echo "darwin";;
        CYGWIN*)    echo "windows";;
        MINGW*)     echo "windows";;
        *)          echo "unknown";;
    esac
}

# Función para detectar la arquitectura
detect_arch() {
    case "$(uname -m)" in
        x86_64)     echo "amd64";;
        aarch64)    echo "arm64";;
        armv7l)     echo "arm";;
        i386)       echo "386";;
        *)          echo "unknown";;
    esac
}

# Función para obtener la ruta de instalación
get_install_path() {
    local os=$(detect_os)
    local home_dir="$HOME"
    
    case "$os" in
        linux|darwin)
            # Para Linux/macOS, instalar en ~/.local/bin
            echo "$home_dir/.local/bin"
            ;;
        windows)
            # Para Windows, instalar en %USERPROFILE%\bin
            echo "$home_dir/bin"
            ;;
        *)
            echo "$home_dir/.local/bin"
            ;;
    esac
}

# Función para verificar si el directorio está en PATH
check_path() {
    local install_dir="$1"
    local path_dirs=($(echo "$PATH" | tr ':' ' '))
    
    for dir in "${path_dirs[@]}"; do
        if [[ "$dir" == "$install_dir" ]]; then
            return 0
        fi
    done
    return 1
}

# Función para agregar directorio al PATH
add_to_path() {
    local install_dir="$1"
    local os=$(detect_os)
    local shell_rc=""
    
    case "$os" in
        linux|darwin)
            # Priorizar zshrc sobre bashrc para usuarios de zsh
            if [[ -f "$HOME/.zshrc" ]]; then
                shell_rc="$HOME/.zshrc"
            elif [[ -f "$HOME/.bashrc" ]]; then
                shell_rc="$HOME/.bashrc"
            elif [[ -f "$HOME/.profile" ]]; then
                shell_rc="$HOME/.profile"
            else
                shell_rc="$HOME/.bashrc"
            fi
            
            # Verificar si ya está en el archivo
            if ! grep -q "$install_dir" "$shell_rc" 2>/dev/null; then
                print_info "Agregando $install_dir al PATH en $shell_rc"
                echo "" >> "$shell_rc"
                echo "# workflow CLI" >> "$shell_rc"
                echo "export PATH=\"$install_dir:\$PATH\"" >> "$shell_rc"
                print_success "PATH actualizado. Ejecuta 'source $shell_rc' o reinicia tu terminal"
            else
                print_success "Directorio ya está en el PATH"
            fi
            ;;
        windows)
            print_warning "Para Windows, agrega manualmente $install_dir al PATH"
            ;;
    esac
}

# Función principal
main() {
    echo "🌾 workflow CLI - Instalación con Binario Pre-compilado"
    echo "=================================================="
    
    local os=$(detect_os)
    local arch=$(detect_arch)
    
    print_info "Sistema detectado: $os-$arch"
    
    # Verificar que estamos en el directorio correcto
    if [[ ! -d "releases" ]]; then
        print_error "No se encontró el directorio 'releases'. Ejecuta desde el directorio del proyecto."
        exit 1
    fi
    
    # Determinar el archivo de binario correcto
    local binary_file=""
    case "$os" in
        linux)
            if [[ "$arch" == "amd64" ]]; then
                binary_file="releases/workflow-v2.0.0-linux-amd64.tar.gz"
            elif [[ "$arch" == "arm64" ]]; then
                binary_file="releases/workflow-v2.0.0-linux-arm64.tar.gz"
            else
                print_error "Arquitectura no soportada: $arch"
                exit 1
            fi
            ;;
        darwin)
            if [[ "$arch" == "amd64" ]]; then
                binary_file="releases/workflow-v2.0.0-darwin-amd64.tar.gz"
            elif [[ "$arch" == "arm64" ]]; then
                binary_file="releases/workflow-v2.0.0-darwin-arm64.tar.gz"
            else
                print_error "Arquitectura no soportada: $arch"
                exit 1
            fi
            ;;
        windows)
            binary_file="releases/workflow-v2.0.0-windows-amd64.zip"
            ;;
        *)
            print_error "Sistema operativo no soportado: $os"
            exit 1
            ;;
    esac
    
    # Verificar que el archivo existe
    if [[ ! -f "$binary_file" ]]; then
        print_error "No se encontró el binario: $binary_file"
        exit 1
    fi
    
    print_info "📦 Usando binario: $binary_file"
    
    # Obtener ruta de instalación
    local install_dir=$(get_install_path)
    local install_path="$install_dir/workflow"
    
    # Crear directorio de instalación si no existe
    if [[ ! -d "$install_dir" ]]; then
        print_info "Creando directorio de instalación: $install_dir"
        mkdir -p "$install_dir"
    fi
    
    # Hacer backup del binario existente si existe
    if [[ -f "$install_path" ]]; then
        print_info "Haciendo backup del binario existente..."
        mv "$install_path" "$install_path.bak"
    fi
    
    # Extraer e instalar el binario
    print_info "📦 Extrayendo e instalando binario..."
    
    if [[ "$os" == "windows" ]]; then
        # Para Windows, usar unzip
        unzip -o "$binary_file" -d "$install_dir"
    else
        # Para Linux/macOS, usar tar
        # Extraer a un directorio temporal
        local temp_dir=$(mktemp -d)
        tar -xzf "$binary_file" -C "$temp_dir"
        
        # Mover el binario desde el subdirectorio
        local extracted_dir=$(find "$temp_dir" -name "workflow-*" -type d | head -1)
        if [[ -n "$extracted_dir" && -f "$extracted_dir/workflow" ]]; then
            cp "$extracted_dir/workflow" "$install_path"
            rm -rf "$temp_dir"
        else
            print_error "No se pudo encontrar el binario en el archivo extraído"
            exit 1
        fi
    fi
    
    # Hacer ejecutable
    chmod +x "$install_path"
    
    # Verificar instalación
    if [[ -f "$install_path" ]]; then
        print_success "workflow CLI instalado exitosamente!"
    else
        print_error "Error en la instalación"
        exit 1
    fi
    
    # Verificar PATH
    if check_path "$install_dir"; then
        print_success "Directorio ya está en el PATH"
    else
        print_warning "Directorio no está en el PATH"
        add_to_path "$install_dir"
    fi
    
    # Mostrar información final
    echo ""
    print_success "🎉 Instalación completada!"
    echo ""
    print_info "Para usar workflow CLI:"
    echo "  workflow --help"
    echo ""
    print_info "Ejemplos de uso:"
    echo "  workflow add 'Tarea de ejemplo' 2.0"
    echo "  workflow status"
    echo "  workflow report"
    echo "  workflow upgrade"
    echo ""
    
    if [[ "$os" != "windows" ]]; then
        print_warning "Si 'workflow' no funciona, ejecuta:"
        echo "  source ~/.bashrc  # o ~/.zshrc"
        echo "  # O reinicia tu terminal"
    fi
}

# Ejecutar instalación
main "$@" 