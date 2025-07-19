# Changelog
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