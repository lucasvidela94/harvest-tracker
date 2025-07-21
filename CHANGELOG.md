# Changelog
## [2.0.1] - 2025-07-21

### Fixed
- **SQLite Integration**: Fixed `list` and `search` commands to use SQLite instead of JSON
- **Production Testing**: Improved test reliability and debugging
- **Release Process**: Cleaned up release generation to avoid "dirty" tags

### Changed
- **Version**: Bumped to 2.0.1 for stable release
- **Testing**: Enhanced production test with better verification

## [2.0.0] - 2025-07-21

### Added
- **Migración completa a SQLite**: Base de datos robusta con índices optimizados
- **Comando `migrate`**: Migración automática de JSON a SQLite con backup
- **Comando `duplicate`**: Duplicar tareas con opciones de fecha (`--date`, `--yesterday`, `--tomorrow`)
- **Comando `edit`**: Editar tareas existentes por ID con confirmación visual
- **Comando `delete`**: Eliminar tareas con confirmación
- **Comando `complete`**: Marcar tareas como completadas
- **Comando `list`**: Listar tareas con IDs visibles y estados
- **Comando `search`**: Búsqueda avanzada por texto, categoría, estado y fecha
- **Comando `report` mejorado**: Reportes detallados por día, semana y mes con filtros
- **Comando `export`**: Exportación a CSV y JSON con filtros avanzados
- **Estados de tareas**: Sistema completo de estados (pendiente, en progreso, completada, pausada)
- **Flags de fecha**: `--date`, `--yesterday`, `--tomorrow` para agregar tareas en fechas específicas
- **Filtros avanzados**: Por categoría, estado, fecha y texto
- **Índices de base de datos**: Optimización para búsquedas rápidas
- **Sistema de backup**: Backup automático antes de migración
- **Compatibilidad**: Mantiene formato legacy para Harvest app

### Changed
- **Arquitectura**: Migración de JSON a SQLite para mejor rendimiento y escalabilidad
- **Interfaz**: IDs visibles en todos los comandos para facilitar edición
- **Reportes**: Formato detallado con estadísticas y agrupación por fecha
- **Búsqueda**: Búsqueda semántica con múltiples criterios
- **Almacenamiento**: Base de datos SQLite con transacciones e integridad de datos

### Fixed
- **Rendimiento**: Búsquedas y filtros significativamente más rápidos
- **Escalabilidad**: Soporte para miles de tareas sin degradación
- **Integridad**: Transacciones SQLite para evitar pérdida de datos
- **Usabilidad**: IDs visibles para facilitar la gestión de tareas

## [1.1.0] - 2025-01-08

### Added
- Sistema de actualización automática con `harvest --upgrade`
- Verificación de versiones desde GitHub releases
- Backup automático de datos del usuario antes de actualizar
- Restauración automática de datos después de la actualización
- Descarga automática de la última versión desde el repositorio
- Preservación de configuración y datos del usuario durante actualizaciones

### Changed
- Mejorada la gestión de versiones con verificación automática
- Agregado soporte para actualizaciones sin pérdida de datos

### Fixed
- Mejorada la robustez del sistema de instalación

## [1.0.2] - 2025-07-18

### Added
- Nuevas características en esta versión

### Changed
- Cambios en funcionalidades existentes

### Fixed
- Correcciones de bugs



All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-01-07

### Added
- Sistema completo de tracking de tareas para Harvest
- Comando principal `harvest` con subcomandos: add, status, report, daily, finish
- Comando `finish` para completar horas restantes del día
- Comando `week` para generar reportes semanales
- Soporte para categorías: technical, meetings, qa, daily
- Validación de entrada de horas (soporte para lenguaje natural: "two hours", "half", etc.)
- Configuración centralizada en `~/.harvest/config.json`
- Almacenamiento de datos en `~/.harvest/tasks.json`
- Script de instalación universal que detecta shell (bash/zsh)
- Script de desinstalación
- Documentación completa en README.md
- Soporte para daily standups configurables
- Integración con clipboard para copiar reportes

### Features
- Detección automática de shell (bash/zsh)
- Configuración de alias automática
- Instalación de dependencias automática
- Persistencia de datos local
- Reportes formateados para Harvest
- Barra de progreso visual
- Soporte para múltiples zonas horarias 