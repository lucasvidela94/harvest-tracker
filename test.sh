#!/bin/bash

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

echo "🧪 Iniciando tests de workflow CLI..."
echo "=================================="

# Test 1: Verificar que workflow funciona
print_info "Test 1: Verificar que workflow funciona"
workflow version
check_exit "workflow version"

workflow --help
check_exit "workflow --help"

# Test 2: Verificar estructura de directorios
print_info "Test 2: Verificar estructura de directorios"
workflow status
check_exit "workflow status (crea directorio .workflow)"

# Test 3: Agregar tareas básicas
print_info "Test 3: Agregar tareas básicas"
workflow add "Test task 1" 2.0
check_exit "workflow add task 1"

workflow add "Test task 2" 1.5
check_exit "workflow add task 2"

workflow add "Test task 3" 3.0
check_exit "workflow add task 3"

# Test 4: Verificar estado
print_info "Test 4: Verificar estado"
workflow status
check_exit "workflow status"

# Test 5: Listar tareas
print_info "Test 5: Listar tareas"
workflow list
check_exit "workflow list"

# Test 6: Editar tarea
print_info "Test 6: Editar tarea"
workflow edit 1 --hours 2.5
check_exit "workflow edit task 1"

workflow edit 1 --description "Test task 1 updated"
check_exit "workflow edit description"

# Test 7: Completar tarea
print_info "Test 7: Completar tarea"
echo "y" | workflow complete 1
check_exit "workflow complete task 1"

# Test 8: Agregar tareas con categorías específicas
print_info "Test 8: Agregar tareas con categorías específicas"
workflow tech "Technical task" 2.0
check_exit "workflow tech"

workflow meeting "Team meeting" 1.0
check_exit "workflow meeting"

workflow qa "QA testing" 1.5
check_exit "workflow qa"

workflow daily
check_exit "workflow daily"

# Test 9: Agregar tareas con fechas específicas
print_info "Test 9: Agregar tareas con fechas específicas"
workflow add "Yesterday task" 1.0 --yesterday
check_exit "workflow add --yesterday"

workflow add "Tomorrow task" 1.0 --tomorrow
check_exit "workflow add --tomorrow"

workflow add "Specific date task" 1.0 --date 2025-07-20
check_exit "workflow add --date"

# Test 10: Búsqueda
print_info "Test 10: Búsqueda"
workflow search "test"
check_exit "workflow search"

workflow search --category tech
check_exit "workflow search --category"

workflow search --status completed
check_exit "workflow search --status"

# Test 11: Reportes
print_info "Test 11: Reportes"
workflow report
check_exit "workflow report"

workflow report --workflow
check_exit "workflow report --workflow"

workflow report --week
check_exit "workflow report --week"

# Test 12: Exportación
print_info "Test 12: Exportación"
workflow export --format csv --output test-export.csv
check_exit "workflow export csv"

workflow export --format json --output test-export.json
check_exit "workflow export json"

# Verificar archivos exportados
check_file "test-export.csv"
check_file "test-export.json"

# Test 13: Duplicar tarea
print_info "Test 13: Duplicar tarea"
workflow duplicate 2
check_exit "workflow duplicate"

workflow duplicate 2 --tomorrow
check_exit "workflow duplicate --tomorrow"

# Test 14: Eliminar tarea
print_info "Test 14: Eliminar tarea"
workflow delete 5 --force
check_exit "workflow delete"

# Test 15: Verificar migración (simular datos JSON viejos)
print_info "Test 15: Verificar migración"
# Crear datos JSON viejos simulados
mkdir -p ~/.workflow
cat > ~/.workflow/tasks.json << 'EOF'
[
  {
    "id": 1,
    "description": "Old JSON task 1",
    "hours": 2.0,
    "category": "general",
    "date": "2025-07-21",
    "created_at": "2025-07-21T10:00:00Z"
  },
  {
    "id": 2,
    "description": "Old JSON task 2",
    "hours": 1.5,
    "category": "tech",
    "date": "2025-07-21",
    "created_at": "2025-07-21T11:00:00Z"
  }
]
EOF

# Ejecutar migración
workflow migrate --dry-run
check_exit "workflow migrate --dry-run"

# Test 16: Verificar directorio de datos
print_info "Test 16: Verificar directorio de datos"
ls -la ~/.workflow/
check_exit "list .workflow directory"

# Test 17: Verificar archivo de configuración
print_info "Test 17: Verificar archivo de configuración"
if [ -f ~/.workflow/config.json ]; then
    print_success "config.json existe"
    cat ~/.workflow/config.json
else
    print_warning "config.json no existe (se creará automáticamente)"
fi

# Test 18: Verificar base de datos SQLite
print_info "Test 18: Verificar base de datos SQLite"
if [ -f ~/.workflow/tasks.db ]; then
    print_success "tasks.db existe"
    # Verificar que podemos leer la base de datos
    sqlite3 ~/.workflow/tasks.db "SELECT COUNT(*) FROM tasks;" 2>/dev/null
    check_exit "SQLite database is readable"
else
    print_warning "tasks.db no existe (se creará automáticamente)"
fi

echo ""
echo "🎉 Tests completados!"
echo "===================="
print_success "Todos los tests básicos pasaron"
echo ""
echo "📊 Resumen de funcionalidades testeadas:"
echo "  ✅ Comandos básicos (add, edit, complete, delete)"
echo "  ✅ Categorías específicas (tech, meeting, qa, daily)"
echo "  ✅ Fechas (yesterday, tomorrow, specific date)"
echo "  ✅ Búsqueda y filtros"
echo "  ✅ Reportes (diario, semanal, formato legacy)"
echo "  ✅ Exportación (CSV, JSON)"
echo "  ✅ Duplicación de tareas"
echo "  ✅ Migración de datos"
echo "  ✅ Estructura de archivos"
echo ""
print_info "Para ejecutar tests específicos, puedes modificar este script"
print_info "o agregar más casos de prueba según necesites." 