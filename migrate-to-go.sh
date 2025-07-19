#!/bin/bash

# Script de migración de Python a Go para Harvest CLI
# Maneja la transición completa del sistema anterior al nuevo

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
PYTHON_SCRIPT="harvest"
GO_DIR="harvest-go"
BACKUP_DIR="harvest-python-backup"
CURRENT_VERSION=$(cat VERSION 2>/dev/null || echo "unknown")

# Función para hacer backup del sistema Python
backup_python_system() {
    print_info "💾 Creando backup del sistema Python..."
    
    # Crear directorio de backup
    mkdir -p "$BACKUP_DIR"
    
    # Backup del script principal
    if [[ -f "$PYTHON_SCRIPT" ]]; then
        cp "$PYTHON_SCRIPT" "$BACKUP_DIR/"
        print_success "Script Python respaldado"
    fi
    
    # Backup de scripts de instalación
    if [[ -f "install.sh" ]]; then
        cp "install.sh" "$BACKUP_DIR/"
    fi
    
    if [[ -f "uninstall.sh" ]]; then
        cp "uninstall.sh" "$BACKUP_DIR/"
    fi
    
    if [[ -f "release.sh" ]]; then
        cp "release.sh" "$BACKUP_DIR/"
    fi
    
    # Backup de documentación
    if [[ -f "README.md" ]]; then
        cp "README.md" "$BACKUP_DIR/"
    fi
    
    if [[ -f "LICENSE" ]]; then
        cp "LICENSE" "$BACKUP_DIR/"
    fi
    
    if [[ -f "VERSION" ]]; then
        cp "VERSION" "$BACKUP_DIR/"
    fi
    
    # Backup de archivos de configuración
    if [[ -d ".harvest" ]]; then
        cp -r ".harvest" "$BACKUP_DIR/"
        print_success "Configuración respaldada"
    fi
    
    # Crear archivo de información del backup
    cat > "$BACKUP_DIR/backup-info.txt" << EOF
Harvest Python System Backup
Created: $(date)
Original Version: $CURRENT_VERSION
Backup Contents:
- Script principal: $PYTHON_SCRIPT
- Scripts de instalación
- Documentación
- Configuración de usuario
- Archivos de datos

Para restaurar:
1. Copiar archivos de vuelta al directorio principal
2. Ejecutar: chmod +x harvest
3. Verificar funcionamiento: ./harvest --help
EOF
    
    print_success "Backup completado en $BACKUP_DIR"
}

# Función para instalar el sistema Go
install_go_system() {
    print_info "🚀 Instalando sistema Go..."
    
    if [[ ! -d "$GO_DIR" ]]; then
        print_error "Directorio $GO_DIR no encontrado"
        exit 1
    fi
    
    # Cambiar al directorio Go
    cd "$GO_DIR"
    
    # Instalar usando el script automático
    if [[ -f "install.sh" ]]; then
        print_info "Ejecutando instalación automática..."
        ./install.sh
    else
        print_error "Script de instalación no encontrado"
        exit 1
    fi
    
    # Volver al directorio principal
    cd ..
    
    print_success "Sistema Go instalado"
}

# Función para limpiar archivos Python
cleanup_python_files() {
    print_warning "🧹 Limpiando archivos del sistema Python..."
    
    # Lista de archivos a remover (con confirmación)
    FILES_TO_REMOVE=(
        "$PYTHON_SCRIPT"
        "install.sh"
        "uninstall.sh"
        "release.sh"
        "week"
        "finish"
    )
    
    print_warning "Los siguientes archivos serán removidos:"
    for file in "${FILES_TO_REMOVE[@]}"; do
        if [[ -f "$file" ]]; then
            echo "  - $file"
        fi
    done
    
    echo ""
    print_warning "¿Estás seguro de que quieres continuar? (y/N)"
    read -r response
    if [[ ! "$response" =~ ^[Yy]$ ]]; then
        print_info "Limpieza cancelada"
        return
    fi
    
    # Remover archivos
    for file in "${FILES_TO_REMOVE[@]}"; do
        if [[ -f "$file" ]]; then
            rm "$file"
            print_success "Removido: $file"
        fi
    done
    
    print_success "Limpieza completada"
}

# Función para actualizar documentación
update_documentation() {
    print_info "📝 Actualizando documentación..."
    
    # Mover README del sistema Go al directorio principal
    if [[ -f "$GO_DIR/README.md" ]]; then
        cp "$GO_DIR/README.md" "README.md"
        print_success "README actualizado"
    fi
    
    # Crear archivo de migración
    cat > "MIGRATION_NOTES.md" << EOF
# 🔄 Migración de Python a Go - Harvest CLI

## Resumen de la Migración

**Fecha**: $(date)
**Versión Anterior**: $CURRENT_VERSION (Python)
**Versión Nueva**: 2.0.0 (Go)

## Cambios Realizados

### ✅ Sistema Nuevo (Go)
- **Instalación**: Automática con \`./install.sh\`
- **Comando**: \`harvest\` (disponible globalmente)
- **Soporte**: Multi-plataforma (Linux, macOS, Windows)
- **Upgrade**: Automático con \`harvest upgrade\`
- **Seguridad**: Backup y rollback automáticos

### 📦 Archivos del Sistema Anterior
- **Backup**: Guardado en \`$BACKUP_DIR/\`
- **Restauración**: Posible desde el backup
- **Configuración**: Migrada automáticamente

## Uso del Nuevo Sistema

\`\`\`bash
# Verificar instalación
harvest --help

# Agregar tarea
harvest add "Mi tarea" 2.0

# Ver estado
harvest status

# Actualizar
harvest upgrade
\`\`\`

## Restauración (si es necesario)

Si necesitas volver al sistema Python:

\`\`\`bash
# Restaurar desde backup
cp $BACKUP_DIR/harvest ./
chmod +x harvest

# Verificar funcionamiento
./harvest --help
\`\`\`

## Ventajas del Nuevo Sistema

- ✅ **Performance**: Ejecución más rápida
- ✅ **Distribución**: Un solo binario
- ✅ **Instalación**: Un comando
- ✅ **Actualizaciones**: Automáticas
- ✅ **Seguridad**: Backup automático
- ✅ **Multi-plataforma**: Soporte completo

---

**¡La migración está completa! Disfruta del nuevo Harvest CLI. 🌾**
EOF
    
    print_success "Documentación actualizada"
}

# Función para verificar la migración
verify_migration() {
    print_info "🔍 Verificando migración..."
    
    # Verificar que el comando Go funciona
    if command -v harvest >/dev/null 2>&1; then
        print_success "Comando 'harvest' disponible globalmente"
        harvest version
    else
        print_warning "Comando 'harvest' no está en PATH"
        print_info "Ejecuta: export PATH=\"\$HOME/.local/bin:\$PATH\""
    fi
    
    # Verificar que el backup existe
    if [[ -d "$BACKUP_DIR" ]]; then
        print_success "Backup del sistema Python creado"
        echo "  Ubicación: $BACKUP_DIR"
        echo "  Contenido: $(ls -la "$BACKUP_DIR" | wc -l) archivos"
    else
        print_warning "Backup no encontrado"
    fi
    
    # Verificar que los archivos Python fueron removidos
    if [[ ! -f "$PYTHON_SCRIPT" ]]; then
        print_success "Script Python removido"
    else
        print_warning "Script Python aún existe"
    fi
    
    print_success "Verificación completada"
}

# Función para mostrar ayuda
show_help() {
    echo "🌾 Harvest CLI - Script de Migración Python → Go"
    echo ""
    echo "Uso: $0 [OPCIÓN]"
    echo ""
    echo "Opciones:"
    echo "  --backup-only     Solo crear backup del sistema Python"
    echo "  --install-only    Solo instalar el sistema Go"
    echo "  --cleanup-only    Solo limpiar archivos Python"
    echo "  --verify-only     Solo verificar la migración"
    echo "  --help           Mostrar esta ayuda"
    echo ""
    echo "Sin opciones: Ejecutar migración completa"
    echo ""
    echo "El script realizará:"
    echo "1. Backup del sistema Python"
    echo "2. Instalación del sistema Go"
    echo "3. Limpieza de archivos Python"
    echo "4. Actualización de documentación"
    echo "5. Verificación de la migración"
}

# Función principal
main() {
    case "${1:-}" in
        --backup-only)
            backup_python_system
            ;;
        --install-only)
            install_go_system
            ;;
        --cleanup-only)
            cleanup_python_files
            ;;
        --verify-only)
            verify_migration
            ;;
        --help|-h)
            show_help
            exit 0
            ;;
        "")
            # Migración completa
            print_info "🚀 Iniciando migración completa Python → Go..."
            echo ""
            
            backup_python_system
            echo ""
            
            install_go_system
            echo ""
            
            cleanup_python_files
            echo ""
            
            update_documentation
            echo ""
            
            verify_migration
            echo ""
            
            print_success "🎉 ¡Migración completada exitosamente!"
            echo ""
            print_info "Próximos pasos:"
            echo "1. Probar el nuevo sistema: harvest --help"
            echo "2. Migrar tus tareas: harvest add 'Tarea de prueba' 1.0"
            echo "3. Verificar funcionamiento: harvest status"
            echo ""
            print_info "Si necesitas restaurar el sistema Python:"
            echo "  Revisa $BACKUP_DIR/backup-info.txt"
            ;;
        *)
            print_error "Opción desconocida: $1"
            show_help
            exit 1
            ;;
    esac
}

# Ejecutar función principal
main "$@" 