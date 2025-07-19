#!/bin/bash

# Script de desinstalación para Harvest CLI
# Remueve el binario y limpia la configuración

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

# Función para obtener la ruta de instalación
get_install_path() {
    local os=$(detect_os)
    local home_dir="$HOME"
    
    case "$os" in
        linux|darwin)
            echo "$home_dir/.local/bin"
            ;;
        windows)
            echo "$home_dir/bin"
            ;;
        *)
            echo "$home_dir/.local/bin"
            ;;
    esac
}

# Función para remover del PATH
remove_from_path() {
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
            fi
            ;;
        windows)
            print_warning "Para Windows, remueve manualmente $install_dir del PATH"
            return
            ;;
    esac
    
    if [[ -n "$shell_rc" ]]; then
        # Crear backup del archivo
        cp "$shell_rc" "$shell_rc.backup"
        
        # Remover líneas relacionadas con Harvest (de forma más segura)
        if grep -q "harvest" "$shell_rc"; then
            # Remover líneas que contengan alias de harvest
            sed -i '/alias harvest=/d' "$shell_rc"
            sed -i '/alias finish=/d' "$shell_rc"
            sed -i '/alias week=/d' "$shell_rc"
            
            print_success "Removido del PATH en $shell_rc"
        else
            print_info "No se encontraron configuraciones de Harvest en $shell_rc"
        fi
    else
        print_warning "No se pudo encontrar archivo de configuración del shell"
        print_info "Remueve manualmente $install_dir del PATH"
    fi
}

# Función para limpiar datos
cleanup_data() {
    local home_dir="$HOME"
    local data_dir="$home_dir/.harvest"
    
    if [[ -d "$data_dir" ]]; then
        print_warning "¿Deseas eliminar todos los datos de Harvest? (y/N)"
        read -r response
        if [[ "$response" =~ ^[Yy]$ ]]; then
            rm -rf "$data_dir"
            print_success "Datos eliminados"
        else
            print_info "Datos preservados en $data_dir"
        fi
    fi
}

# Función principal de desinstalación
main() {
    print_info "🗑️  Desinstalando Harvest CLI..."
    
    # Obtener ruta de instalación
    local install_dir=$(get_install_path)
    local install_path="$install_dir/harvest"
    
    # Verificar si está instalado
    if [[ ! -f "$install_path" ]]; then
        print_warning "Harvest CLI no está instalado en $install_path"
        print_info "Buscando en otras ubicaciones..."
        
        # Buscar en PATH
        if command -v harvest >/dev/null 2>&1; then
            local found_path=$(which harvest)
            print_info "Encontrado en: $found_path"
            install_path="$found_path"
        else
            print_error "Harvest CLI no encontrado en el sistema"
            exit 1
        fi
    fi
    
    # Confirmar desinstalación
    print_warning "¿Estás seguro de que quieres desinstalar Harvest CLI? (y/N)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        print_info "Desinstalación cancelada"
        exit 0
    fi
    
    # Hacer backup del binario
    if [[ -f "$install_path" ]]; then
        print_info "Haciendo backup del binario..."
        cp "$install_path" "$install_path.backup"
    fi
    
    # Remover binario
    if [[ -f "$install_path" ]]; then
        rm "$install_path"
        print_success "Binario removido"
    fi
    
    # Remover del PATH
    remove_from_path "$install_dir"
    
    # Limpiar datos (opcional)
    cleanup_data
    
    # Verificar si el directorio está vacío
    if [[ -d "$install_dir" ]] && [[ -z "$(ls -A "$install_dir")" ]]; then
        print_info "Removiendo directorio vacío: $install_dir"
        rmdir "$install_dir"
    fi
    
    print_success "🎉 Desinstalación completada!"
    echo ""
    print_info "Si cambiaste de opinión, puedes restaurar desde:"
    echo "  $install_path.backup"
}

# Ejecutar desinstalación
main "$@" 