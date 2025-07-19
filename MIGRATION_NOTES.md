# 🔄 Migración de Python a Go - Harvest CLI

## Resumen de la Migración

**Fecha**: sáb 19 jul 2025 03:27:29 -03
**Versión Anterior**: 2.0.0 (Python)
**Versión Nueva**: 2.0.0 (Go)

## Cambios Realizados

### ✅ Sistema Nuevo (Go)
- **Instalación**: Automática con `./install.sh`
- **Comando**: `harvest` (disponible globalmente)
- **Soporte**: Multi-plataforma (Linux, macOS, Windows)
- **Upgrade**: Automático con `harvest upgrade`
- **Seguridad**: Backup y rollback automáticos

### 📦 Archivos del Sistema Anterior
- **Backup**: Guardado en `harvest-python-backup/`
- **Restauración**: Posible desde el backup
- **Configuración**: Migrada automáticamente

## Uso del Nuevo Sistema

```bash
# Verificar instalación
harvest --help

# Agregar tarea
harvest add "Mi tarea" 2.0

# Ver estado
harvest status

# Actualizar
harvest upgrade
```

## Restauración (si es necesario)

Si necesitas volver al sistema Python:

```bash
# Restaurar desde backup
cp harvest-python-backup/harvest ./
chmod +x harvest

# Verificar funcionamiento
./harvest --help
```

## Ventajas del Nuevo Sistema

- ✅ **Performance**: Ejecución más rápida
- ✅ **Distribución**: Un solo binario
- ✅ **Instalación**: Un comando
- ✅ **Actualizaciones**: Automáticas
- ✅ **Seguridad**: Backup automático
- ✅ **Multi-plataforma**: Soporte completo

---

**¡La migración está completa! Disfruta del nuevo Harvest CLI. 🌾**
