#!/bin/bash

# Script para crear releases distribuibles de workflow CLI

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

# Variables
VERSION=${1:-$(git describe --tags --always --dirty 2>/dev/null || echo "dev")}
RELEASE_DIR="releases"
DIST_DIR="$RELEASE_DIR/workflow-$VERSION"

# Plataformas soportadas
PLATFORMS=(
    "linux/amd64"
    "linux/arm64"
    "darwin/amd64"
    "darwin/arm64"
    "windows/amd64"
)

# Función para limpiar
cleanup() {
    print_info "🧹 Limpiando archivos temporales..."
    rm -rf "$DIST_DIR"
}

# Función para crear release
create_release() {
    print_info "🚀 Creando release v$VERSION..."
    
    # Crear directorio de release
    mkdir -p "$DIST_DIR"
    
    # Compilar para todas las plataformas
    print_info "🔨 Compilando para todas las plataformas..."
    
    for platform in "${PLATFORMS[@]}"; do
        IFS='/' read -r GOOS GOARCH <<< "$platform"
        BINARY_NAME="workflow"
        
        if [[ "$GOOS" == "windows" ]]; then
            BINARY_NAME="workflow.exe"
        fi
        
        OUTPUT="$DIST_DIR/workflow-$GOOS-$GOARCH"
        if [[ "$GOOS" == "windows" ]]; then
            OUTPUT="$DIST_DIR/workflow-$GOOS-$GOARCH.exe"
        fi
        
        print_info "Compilando para $GOOS/$GOARCH..."
        
        GOOS=$GOOS GOARCH=$GOARCH go build \
            -ldflags "-X main.Version=$VERSION" \
            -o "$OUTPUT" \
            ./cmd/workflow
        
        # Hacer ejecutable (excepto Windows)
        if [[ "$GOOS" != "windows" ]]; then
            chmod +x "$OUTPUT"
        fi
    done
    
    # Copiar archivos adicionales
    print_info "📁 Copiando archivos adicionales..."
    
    # Scripts de instalación
    if [[ -f "install.sh" ]]; then
        cp install.sh "$DIST_DIR/"
        chmod +x "$DIST_DIR/install.sh"
    fi
    
    if [[ -f "uninstall.sh" ]]; then
        cp uninstall.sh "$DIST_DIR/"
        chmod +x "$DIST_DIR/uninstall.sh"
    fi
    
    # README
    if [[ -f "README.md" ]]; then
        cp README.md "$DIST_DIR/"
    fi
    
    # LICENSE
    if [[ -f "LICENSE" ]]; then
        cp LICENSE "$DIST_DIR/"
    fi
    
    # Crear archivos tar.gz para cada plataforma
    print_info "📦 Creando archivos de distribución..."
    
    cd "$RELEASE_DIR"
    
    for platform in "${PLATFORMS[@]}"; do
        IFS='/' read -r GOOS GOARCH <<< "$platform"
        PLATFORM_DIR="workflow-$VERSION-$GOOS-$GOARCH"
        
        # Crear directorio específico para la plataforma
        mkdir -p "$PLATFORM_DIR"
        
        # Copiar binario específico
        if [[ "$GOOS" == "windows" ]]; then
            cp "workflow-$VERSION/workflow-$GOOS-$GOARCH.exe" "$PLATFORM_DIR/workflow.exe"
        else
            cp "workflow-$VERSION/workflow-$GOOS-$GOARCH" "$PLATFORM_DIR/workflow"
            chmod +x "$PLATFORM_DIR/workflow"
        fi
        
        # Copiar archivos adicionales
        cp workflow-$VERSION/install.sh "$PLATFORM_DIR/" 2>/dev/null || true
        cp workflow-$VERSION/uninstall.sh "$PLATFORM_DIR/" 2>/dev/null || true
        cp workflow-$VERSION/README.md "$PLATFORM_DIR/" 2>/dev/null || true
        cp workflow-$VERSION/LICENSE "$PLATFORM_DIR/" 2>/dev/null || true
        
        # Crear archivo tar.gz
        tar -czf "$PLATFORM_DIR.tar.gz" "$PLATFORM_DIR"
        
        # Limpiar directorio temporal
        rm -rf "$PLATFORM_DIR"
        
        print_success "Creado $PLATFORM_DIR.tar.gz"
    done
    
    # Crear archivo ZIP para Windows
    if command -v zip >/dev/null 2>&1; then
        for platform in "${PLATFORMS[@]}"; do
            IFS='/' read -r GOOS GOARCH <<< "$platform"
            
            if [[ "$GOOS" == "windows" ]]; then
                PLATFORM_DIR="workflow-$VERSION-$GOOS-$GOARCH"
                
                # Crear directorio específico para la plataforma
                mkdir -p "$PLATFORM_DIR"
                
                # Copiar binario específico
                cp "workflow-$VERSION/workflow-$GOOS-$GOARCH.exe" "$PLATFORM_DIR/workflow.exe"
                
                # Copiar archivos adicionales
                cp workflow-$VERSION/install.sh "$PLATFORM_DIR/" 2>/dev/null || true
                cp workflow-$VERSION/uninstall.sh "$PLATFORM_DIR/" 2>/dev/null || true
                cp workflow-$VERSION/README.md "$PLATFORM_DIR/" 2>/dev/null || true
                cp workflow-$VERSION/LICENSE "$PLATFORM_DIR/" 2>/dev/null || true
                
                # Crear archivo ZIP
                zip -r "$PLATFORM_DIR.zip" "$PLATFORM_DIR"
                
                # Limpiar directorio temporal
                rm -rf "$PLATFORM_DIR"
                
                print_success "Creado $PLATFORM_DIR.zip"
            fi
        done
    fi
    
    cd ..
    
    # Crear checksums
    print_info "🔐 Generando checksums..."
    
    cd "$RELEASE_DIR"
    
    # SHA256 checksums
    for file in *.tar.gz *.zip; do
        if [[ -f "$file" ]]; then
            sha256sum "$file" >> "workflow-$VERSION-checksums.txt"
        fi
    done
    
    cd ..
    
    print_success "🎉 Release v$VERSION creado exitosamente!"
}

# Función para mostrar ayuda
show_help() {
    echo "🌾 workflow CLI - Release Script"
    echo ""
    echo "Uso: $0 [VERSION]"
    echo ""
    echo "Argumentos:"
    echo "  VERSION    Versión a crear (opcional, usa git tag por defecto)"
    echo ""
    echo "Ejemplos:"
    echo "  $0           # Usar versión de git tag"
    echo "  $0 1.2.0     # Crear release v1.2.0"
    echo ""
    echo "El script creará:"
    echo "  - Binarios para todas las plataformas"
    echo "  - Archivos tar.gz para Linux/macOS"
    echo "  - Archivos ZIP para Windows"
    echo "  - Checksums SHA256"
    echo "  - Scripts de instalación"
}

# Función principal
main() {
    # Verificar argumentos
    if [[ "$1" == "-h" || "$1" == "--help" ]]; then
        show_help
        exit 0
    fi
    
    # Verificar que estamos en el directorio correcto
    if [[ ! -f "go.mod" ]]; then
        print_error "No se encontró go.mod. Ejecuta desde el directorio del proyecto."
        exit 1
    fi
    
    # Verificar que Go está instalado
    if ! command -v go >/dev/null 2>&1; then
        print_error "Go no está instalado o no está en PATH"
        exit 1
    fi
    
    # Crear release
    create_release
    
    # Mostrar resumen
    echo ""
    print_success "📊 Resumen del release:"
    echo "  Versión: $VERSION"
    echo "  Directorio: $RELEASE_DIR"
    echo "  Archivos creados:"
    
    cd "$RELEASE_DIR"
    for file in *.tar.gz *.zip *.txt; do
        if [[ -f "$file" ]]; then
            size=$(du -h "$file" | cut -f1)
            echo "    - $file ($size)"
        fi
    done
    cd ..
    
    echo ""
    print_info "🚀 Para publicar el release:"
    echo "1. Sube los archivos a GitHub Releases"
    echo "2. Etiqueta el release como v$VERSION"
    echo "3. Incluye el changelog correspondiente"
}

# Ejecutar con cleanup en caso de error
trap cleanup EXIT

# Ejecutar función principal
main "$@" 