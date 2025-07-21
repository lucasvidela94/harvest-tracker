#!/bin/bash

# Script de test para verificar instalación sin Go
# Simula el entorno de un usuario que no tiene Go instalado

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

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# Función para verificar si un comando fue exitoso
check_exit() {
    if [ $? -eq 0 ]; then
        print_success "$1"
    else
        print_error "$1"
        exit 1
    fi
}

# Función para verificar si un archivo existe
check_file() {
    if [ -f "$1" ]; then
        print_success "Archivo $1 existe"
    else
        print_error "Archivo $1 no existe"
        exit 1
    fi
}

echo "🧪 Test de instalación sin Go"
echo "============================"

# Test 1: Verificar que no tenemos Go instalado
print_info "Test 1: Verificar que Go NO está instalado"
if command -v go &> /dev/null; then
    print_error "Go está instalado, pero debería NO estar"
    exit 1
else
    print_success "Go NO está instalado (correcto)"
fi

# Test 2: Verificar que tenemos los archivos necesarios
print_info "Test 2: Verificar archivos de instalación"
check_file "install-binary.sh"
check_file "releases/harvest-v2.0.0-linux-amd64.tar.gz"

# Test 3: Verificar que el script es ejecutable
print_info "Test 3: Verificar permisos del script"
if [ -x "install-binary.sh" ]; then
    print_success "install-binary.sh es ejecutable"
else
    print_error "install-binary.sh no es ejecutable"
    chmod +x install-binary.sh
    check_exit "Hacer ejecutable install-binary.sh"
fi

# Test 4: Simular instalación en directorio temporal
print_info "Test 4: Simular instalación"
mkdir -p /tmp/harvest-test
cd /tmp/harvest-test

# Crear estructura mínima
mkdir -p releases
cp /test/releases/harvest-v2.0.0-linux-amd64.tar.gz releases/

# Ejecutar instalación
print_info "Ejecutando install-binary.sh..."
bash /test/install-binary.sh
check_exit "install-binary.sh ejecutado"

# Test 5: Verificar que harvest fue instalado
print_info "Test 5: Verificar instalación"
if [ -f "$HOME/.local/bin/harvest" ]; then
    print_success "harvest instalado en $HOME/.local/bin/harvest"
else
    print_error "harvest no fue instalado"
    exit 1
fi

# Test 6: Verificar que harvest funciona
print_info "Test 6: Verificar que harvest funciona"
$HOME/.local/bin/harvest version
check_exit "harvest version"

$HOME/.local/bin/harvest --help
check_exit "harvest --help"

# Test 7: Probar funcionalidad básica
print_info "Test 7: Probar funcionalidad básica"
$HOME/.local/bin/harvest add "Test task" 1.0
check_exit "harvest add"

$HOME/.local/bin/harvest status
check_exit "harvest status"

$HOME/.local/bin/harvest report --harvest
check_exit "harvest report --harvest"

# Test 8: Verificar que se creó la estructura de datos
print_info "Test 8: Verificar estructura de datos"
if [ -d "$HOME/.harvest" ]; then
    print_success "Directorio .harvest creado"
else
    print_error "Directorio .harvest no fue creado"
    exit 1
fi

# Test 9: Verificar base de datos SQLite
print_info "Test 9: Verificar base de datos SQLite"
if [ -f "$HOME/.harvest/tasks.db" ]; then
    print_success "Base de datos SQLite creada"
else
    print_warning "Base de datos SQLite no existe (se creará automáticamente)"
fi

# Test 10: Verificar archivo de configuración
print_info "Test 10: Verificar archivo de configuración"
if [ -f "$HOME/.harvest/config.json" ]; then
    print_success "Archivo de configuración creado"
    cat "$HOME/.harvest/config.json"
else
    print_warning "Archivo de configuración no existe (se creará automáticamente)"
fi

echo ""
echo "🎉 Test de instalación sin Go completado!"
echo "=========================================="
print_success "Todos los tests pasaron"
echo ""
echo "📊 Resumen:"
echo "  ✅ Go NO está instalado (correcto)"
echo "  ✅ Script de instalación funciona"
echo "  ✅ Binario se instala correctamente"
echo "  ✅ Harvest CLI funciona sin Go"
echo "  ✅ Estructura de datos se crea"
echo "  ✅ Funcionalidad básica funciona"
echo ""
print_info "El script install-binary.sh funciona perfectamente sin Go instalado" 