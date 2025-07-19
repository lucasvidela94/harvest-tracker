#!/bin/bash

# Script de instalación para Harvest CLI
# Instala el binario en el PATH del sistema

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

# Función para agregar al PATH
add_to_path() {
    local install_dir="$1"
    local os=$(detect_os)
    local shell_rc=""
    
    case "$os" in
        linux|darwin)
            if [[ -f "$HOME/.bashrc" ]]; then
                shell_rc="$HOME/.bashrc"
            elif [[ -f "$HOME/.zshrc" ]]; then
                shell_rc="$HOME/.zshrc"
            elif [[ -f "$HOME/.profile" ]]; then
                shell_rc="$HOME/.profile"
            fi
            ;;
        windows)
            # En Windows, el PATH se maneja diferente
            print_warning "Para Windows, agrega manualmente $install_dir al PATH"
            return
            ;;
    esac
    
    if [[ -n "$shell_rc" ]]; then
        if ! grep -q "$install_dir" "$shell_rc"; then
            echo "" >> "$shell_rc"
            echo "# Harvest CLI" >> "$shell_rc"
            echo "export PATH=\"$install_dir:\$PATH\"" >> "$shell_rc"
            print_success "Agregado al PATH en $shell_rc"
            print_info "Ejecuta 'source $shell_rc' o reinicia tu terminal"
        else
            print_info "Ya está en el PATH"
        fi
    else
        print_warning "No se pudo encontrar archivo de configuración del shell"
        print_info "Agrega manualmente $install_dir al PATH"
    fi
}

# Función principal de instalación
main() {
    print_info "🚀 Instalando Harvest CLI..."
    
    # Detectar sistema
    local os=$(detect_os)
    local arch=$(detect_arch)
    
    print_info "Sistema detectado: $os ($arch)"
    
    if [[ "$os" == "unknown" || "$arch" == "unknown" ]]; then
        print_error "Sistema operativo o arquitectura no soportada"
        exit 1
    fi
    
    # Verificar que estamos en el directorio correcto
    if [[ ! -f "go.mod" ]]; then
        print_error "No se encontró go.mod. Ejecuta desde el directorio del proyecto."
        exit 1
    fi
    
    # Compilar el proyecto
    print_info "🔨 Compilando Harvest CLI..."
    if ! go build -o harvest ./cmd/harvest; then
        print_error "Error al compilar el proyecto"
        exit 1
    fi
    
    # Obtener ruta de instalación
    local install_dir=$(get_install_path)
    local install_path="$install_dir/harvest"
    
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
    
    # Instalar el binario
    print_info "📦 Instalando en: $install_path"
    if ! cp harvest "$install_path"; then
        print_error "Error al copiar el binario"
        exit 1
    fi
    
    # Hacer ejecutable
    chmod +x "$install_path"
    
    # Verificar instalación
    if [[ -f "$install_path" ]]; then
        print_success "Harvest CLI instalado exitosamente!"
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
    print_info "Para usar Harvest CLI:"
    echo "  harvest --help"
    echo ""
    print_info "Ejemplos de uso:"
    echo "  harvest add 'Tarea de ejemplo' 2.0"
    echo "  harvest status"
    echo "  harvest report"
    echo "  harvest upgrade"
    echo ""
    
    if [[ "$os" != "windows" ]]; then
        print_warning "Si 'harvest' no funciona, ejecuta:"
        echo "  source ~/.bashrc  # o ~/.zshrc"
        echo "  # O reinicia tu terminal"
    fi
}

# Ejecutar instalación
main "$@" 