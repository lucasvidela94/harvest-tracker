# 🔄 Workflow CLI

**Una herramienta de línea de comandos enterprise para el seguimiento de tareas y gestión de workflows productivos.**

[![Release](https://img.shields.io/github/v/release/lucasvidela94/workflow-cli)](https://github.com/lucasvidela94/workflow-cli/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/lucasvidela94/workflow-cli)](https://goreportcard.com/report/github.com/lucasvidela94/workflow-cli)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

> **✨ Proyecto completamente migrado a Go** - Mejor rendimiento, mantenibilidad y distribución multiplataforma.

## 🚀 **Instalación Enterprise**

### **Instalación Automática (Recomendada)**
```bash
# Instalar la última versión automáticamente
curl -fsSL https://raw.githubusercontent.com/lucasvidela94/workflow-cli/main/install-latest.sh | bash
```

### **Instalación Manual**
```bash
# Descargar para tu plataforma
wget https://github.com/lucasvidela94/workflow-cli/releases/latest/download/workflow-$(curl -s https://api.github.com/repos/lucasvidela94/workflow-cli/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/').tar.gz

# Extraer e instalar
tar -xzf workflow-*.tar.gz
sudo mv workflow-*/workflow /usr/local/bin/
```

### **Verificar Instalación**
```bash
workflow version
workflow --help
```

## 🏢 **Características Enterprise**

### **✅ Funcionalidades Avanzadas**
- **SQLite Database**: Base de datos robusta y escalable
- **Búsqueda Avanzada**: Por texto, categoría, estado y fecha
- **Gestión Completa**: Editar, eliminar, duplicar tareas
- **Exportación**: CSV y JSON con filtros avanzados
- **Estados de Tareas**: Pendiente, en progreso, completada, pausada
- **Migración Automática**: De JSON a SQLite con backup

### **✅ Distribución Profesional**
- **One-liner Installation**: Instalación automática multiplataforma
- **Auto-Update**: Sistema de actualización automática
- **Checksums**: Verificación de integridad de archivos
- **Multiplatform**: Linux, macOS, Windows (amd64, arm64)

### **✅ Integración Enterprise**
- **CI/CD Ready**: Fácil integración en pipelines
- **Docker Support**: Contenedores listos para producción
- **API Integration**: Preparado para integraciones futuras

## 🔄 **Auto-Update**

workflow CLI se actualiza automáticamente:

```bash
# Verificar actualizaciones
workflow check-update

# Actualizar a la última versión
workflow upgrade
```

## 🚀 Instalación desde Código Fuente

### Opción 1: Instalación con Binario Pre-compilado (Desarrollo)

**No requiere Go instalado - ¡Más fácil!**

```bash
# Clonar el repositorio
git clone https://github.com/lucasvidela94/workflow-cli.git
cd workflow-cli

# Instalar usando binario pre-compilado
./install-binary.sh
```

### Opción 2: Instalación desde Código Fuente

**Requiere Go 1.24.5+ instalado**

```bash
# Clonar el repositorio
git clone https://github.com/lucasvidela94/workflow-cli.git
cd workflow-cli

# Instalar usando el script automático
./install.sh
```

### Opción 3: Instalación Manual

```bash
# Compilar e instalar
make install-script

# O manualmente
make build
make install
```

## 📋 Uso

Una vez instalado, puedes usar `workflow` desde cualquier lugar:

```bash
# Ver ayuda completa
workflow --help

# Agregar tareas
workflow add "Desarrollar nueva funcionalidad" 4.0
workflow tech "API development" 3.5
workflow meeting "Sprint planning" 1.5
workflow qa "Testing new features" 2.0
workflow daily

# Ver estado y progreso
workflow status

# Gestionar tareas
workflow list                    # Listar todas las tareas
workflow list --date 2025-07-21  # Tareas de fecha específica
workflow search "bug"            # Buscar tareas
workflow edit 1 --hours 3.0      # Editar tarea por ID
workflow delete 2                # Eliminar tarea
workflow complete 3              # Marcar como completada

# Reportes y exportación
workflow report                  # Reporte de productividad
workflow export --format csv     # Exportar a CSV
workflow export --format json    # Exportar a JSON

# Actualización
workflow check-update            # Verificar actualizaciones
workflow upgrade                 # Actualizar automáticamente
```

## 🚀 Flujo de Trabajo Diario - Un Día en la Vida de un Dev

### 🌅 **Mañana (9:00 AM) - Planificación del Día**

```bash
# Ver qué tareas quedaron pendientes de ayer
workflow list --date 2025-07-20

# Agregar el daily standup
workflow daily

# Agregar tareas planificadas para hoy
workflow add "Revisar PRs pendientes" 1.0
workflow add "Desarrollar feature de login" 4.0
workflow add "Reunión de planning semanal" 1.5

# Ver el estado inicial del día
workflow status
```

**Output:**
```
📅 Today (2025-07-21): 6.5h / 8.0h
📈 Remaining: 1.5h
  [1] 📝 Revisar PRs pendientes (1.0h, general) ⏳
  [2] 💻 Desarrollar feature de login (4.0h, tech) ⏳
  [3] 📝 Reunión de planning semanal (1.5h, meeting) ⏳
📊 [█████████████░░░░░░░] 81.3%
```

### ☕ **Media Mañana (11:00 AM) - Ajustes y Progreso**

```bash
# Completar la revisión de PRs
workflow complete 1

# Agregar una tarea que surgió (bug fix urgente)
workflow add "Fix bug crítico en producción" 2.0

# Ver estado actualizado
workflow status
```

**Output:**
```
📅 Today (2025-07-21): 8.5h / 8.0h
📈 Overtime: 0.5h
  [1] 📝 Revisar PRs pendientes (1.0h, general) ✅
  [2] 💻 Desarrollar feature de login (4.0h, tech) ⏳
  [3] 📝 Reunión de planning semanal (1.5h, meeting) ⏳
  [4] 📝 Fix bug crítico en producción (2.0h, general) ⏳
📊 [██████████████████░░] 106.3%
```

### 🍽️ **Almuerzo (1:00 PM) - Revisión y Ajustes**

```bash
# Ver qué tareas tenemos y reorganizar
workflow list

# Editar la tarea del bug para ser más específica
workflow edit 4 --description "Fix bug crítico en API de usuarios" --hours 1.5

# Buscar tareas similares para referencia
workflow search "bug"
```

### 🌆 **Tarde (3:00 PM) - Progreso y Nuevas Tareas**

```bash
# Completar el bug fix
workflow complete 4

# Agregar tarea que surgió durante el desarrollo
workflow add "Documentar nueva API" 1.0

# Duplicar tarea de mañana para mañana (recurrente)
workflow duplicate 1 --tomorrow

# Ver progreso del día
workflow status
```

### 🌙 **Fin de Día (5:30 PM) - Cierre y Reporte**

```bash
# Completar tareas pendientes
workflow complete 2
workflow complete 3

# Ver reporte final del día
workflow report

# Generar reporte para workflow (formato legacy)
workflow report --workflow
```

**Output del reporte final:**
```
📊 Report for 2025-07-21
==================================================
📋 Tasks (5):
[1] 📝 Revisar PRs pendientes (1.0h, general) ✅
[2] 💻 Desarrollar feature de login (4.0h, tech) ✅
[3] 📝 Reunión de planning semanal (1.5h, meeting) ✅
[4] 📝 Fix bug crítico en API de usuarios (1.5h, general) ✅
[5] 📝 Documentar nueva API (1.0h, general) ⏳

📈 Statistics:
Total hours: 9.0h
Completed: 8.0h
Pending: 1.0h

📊 By category:
  general: 3.5h
  tech: 4.0h
  meeting: 1.5h
```

### 📅 **Viernes - Revisión Semanal**

```bash
# Ver reporte de toda la semana
workflow report --week

# Exportar datos de la semana para análisis
workflow export --format csv --week --output semana-actual.csv

# Buscar tareas técnicas de la semana
workflow search --category tech --week

# Ver tareas completadas vs pendientes
workflow report --status completed --week
workflow report --status pending --week
```

**Output del reporte semanal:**
```
📊 Weekly Report (2025-07-21 to 2025-07-27)
==================================================

📅 2025-07-21:
  [1] 📝 Revisar PRs pendientes (1.0h, general) ✅
  [2] 💻 Desarrollar feature de login (4.0h, tech) ✅
  [3] 📝 Reunión de planning semanal (1.5h, meeting) ✅
  [4] 📝 Fix bug crítico en API de usuarios (1.5h, general) ✅
  [5] 📝 Documentar nueva API (1.0h, general) ⏳
  Total: 9.0h

📈 Weekly Summary:
Total hours: 38.5h
Completed: 35.0h
Completion rate: 90.9%

📊 By category:
  tech: 20.0h
  general: 12.5h
  meeting: 6.0h
```

## 🎯 **Beneficios del Flujo Optimizado**

### ✅ **Sin Interrupciones**
- No necesitas abrir workflow durante el día
- Registro de tareas en tiempo real desde la terminal
- Flujo natural que se integra con tu trabajo

### ✅ **Gestión Inteligente**
- IDs visibles para edición rápida
- Estados de tareas para seguimiento
- Búsqueda y filtros avanzados
- Duplicación de tareas recurrentes

### ✅ **Reportes Automáticos**
- Reportes detallados por día, semana y mes
- Exportación a CSV/JSON para análisis
- Formato legacy para copiar a workflow
- Estadísticas y métricas automáticas

### ✅ **Flexibilidad Total**
- Agregar tareas para fechas pasadas/futuras
- Editar tareas existentes sin perder tiempo
- Reorganizar y ajustar según evoluciona el día
- Migración automática de datos antiguos

## 🛠️ Comandos Disponibles

### 📝 Gestión de Tareas
- `workflow add <descripción> <horas>` - Agregar nueva tarea
- `workflow add --date 2025-07-20 <descripción> <horas>` - Agregar tarea para fecha específica
- `workflow add --yesterday <descripción> <horas>` - Agregar tarea para ayer
- `workflow add --tomorrow <descripción> <horas>` - Agregar tarea para mañana
- `workflow tech <descripción> <horas>` - Agregar tarea técnica
- `workflow meeting <descripción> <horas>` - Agregar reunión
- `workflow qa <descripción> <horas>` - Agregar tarea de QA
- `workflow daily` - Agregar daily standup (automático)

### ✏️ Edición y Gestión
- `workflow edit <id> --description "nueva descripción"` - Editar tarea existente
- `workflow edit <id> --hours 2.5` - Cambiar horas de tarea
- `workflow edit <id> --category tech` - Cambiar categoría
- `workflow delete <id>` - Eliminar tarea
- `workflow duplicate <id>` - Duplicar tarea
- `workflow duplicate <id> --tomorrow` - Duplicar tarea para mañana
- `workflow complete <id>` - Marcar tarea como completada

### 📊 Información y Reportes
- `workflow status` - Ver estado actual de tareas
- `workflow list` - Listar tareas con IDs visibles
- `workflow list --date 2025-07-20` - Listar tareas de fecha específica
- `workflow report` - Reporte detallado de hoy
- `workflow report --week` - Reporte semanal
- `workflow report --month` - Reporte mensual
- `workflow report --date 2025-07-20` - Reporte de fecha específica
- `workflow report --category tech` - Reporte filtrado por categoría
- `workflow report --status completed` - Reporte de tareas completadas
- `workflow report --workflow` - Formato legacy para workflow app

### 🔍 Búsqueda y Filtros
- `workflow search "texto"` - Buscar tareas por texto
- `workflow search --category tech` - Buscar por categoría
- `workflow search --status pending` - Buscar por estado
- `workflow search --date 2025-07-20` - Buscar por fecha

### 📤 Exportación
- `workflow export --format csv` - Exportar a CSV
- `workflow export --format json` - Exportar a JSON
- `workflow export --week --format csv` - Exportar semana a CSV
- `workflow export --category tech --format csv` - Exportar tareas técnicas

### 🔄 Migración y Sistema
- `workflow migrate` - Migrar datos de JSON a SQLite
- `workflow migrate --dry-run` - Simular migración
- `workflow migrate --backup-only` - Solo crear backup
- `workflow upgrade` - Actualizar a la última versión
- `workflow rollback` - Gestionar rollbacks

## ⚙️ Configuración

El CLI se configura automáticamente en `~/.workflow/`:

- `config.json` - Configuración general
- `tasks.db` - Base de datos SQLite con todas las tareas
- `tasks.json.backup.*` - Backups automáticos de datos JSON (si migraste)

### Migración de Datos

Si tienes datos en el formato JSON anterior, la migración es automática:

```bash
# Migrar datos existentes a SQLite
workflow migrate

# Simular migración sin cambios
workflow migrate --dry-run

# Solo crear backup
workflow migrate --backup-only
```

## 🔄 Actualizaciones

El sistema incluye un sistema de upgrade automático:

```bash
# Verificar actualizaciones
workflow upgrade
```

## 🛡️ Seguridad

- **Backup automático** antes de cualquier cambio
- **Verificación de integridad** en cada paso
- **Rollback automático** en caso de fallo
- **Logs detallados** para auditoría

## 🖥️ Plataformas Soportadas

- **Linux**: amd64, arm64
- **macOS**: amd64, arm64
- **Windows**: amd64

## 🗑️ Desinstalación

```bash
# Desinstalación automática
./uninstall.sh

# O manualmente
make uninstall-script
```

## 🐛 Solución de Problemas

### El comando `workflow` no funciona

```bash
# Verificar instalación
make check

# Si está instalado pero no en PATH
export PATH="$HOME/.local/bin:$PATH"

# O agregar permanentemente a tu shell
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

### Verificar instalación

```bash
# Verificar que funciona
workflow --version
workflow --help
```

## 🔧 Desarrollo

### Compilar desde código fuente

```bash
# Clonar repositorio
git clone https://github.com/lucasvidela94/workflow-cli.git
cd workflow-cli

# Instalar dependencias
go mod tidy

# Compilar
go build -o workflow ./cmd/workflow

# Ejecutar
./workflow --help
```

### Comandos de desarrollo

```bash
# Compilar
make build

# Compilar para todas las plataformas
make build-all

# Ejecutar tests
make test

# Verificar código
make code-check

# Modo desarrollo
make dev
```

## 📁 Estructura del Proyecto

```
workflow/
├── cmd/workflow/          # Punto de entrada principal
├── internal/             # Lógica interna del proyecto
│   ├── cli/             # Comandos CLI
│   ├── core/            # Lógica principal
│   └── upgrade/         # Sistema de upgrade
├── pkg/workflow/         # Tipos y utilidades
├── build/               # Archivos de build para múltiples plataformas
├── releases/            # Releases compilados
├── install.sh           # Script de instalación
├── uninstall.sh         # Script de desinstalación
├── release.sh           # Script de release
├── Makefile             # Comandos de build y desarrollo
├── go.mod               # Dependencias Go
├── go.sum               # Checksums de dependencias
├── README.md            # Este archivo
├── CHANGELOG.md         # Historial de cambios
├── LICENSE              # Licencia del proyecto
├── VERSION              # Versión actual
└── workflow              # Ejecutable compilado
```

## 🎯 Características Principales

- **⚡ Alto Rendimiento**: Escrito en Go con base de datos SQLite optimizada
- **🔧 Fácil Instalación**: Scripts automáticos de instalación
- **🔄 Actualizaciones Automáticas**: Sistema de upgrade integrado
- **🛡️ Seguridad**: Backup automático y migración segura de datos
- **📱 Multiplataforma**: Soporte para Linux, macOS y Windows
- **📊 Reportes Avanzados**: Reportes detallados por día, semana y mes
- **🔍 Búsqueda Inteligente**: Búsqueda semántica con múltiples filtros
- **✏️ Edición en Tiempo Real**: Editar tareas sin interrumpir el flujo
- **📤 Exportación Flexible**: Exportar a CSV y JSON con filtros
- **🔄 Migración Automática**: Migración transparente de JSON a SQLite
- **📈 Estados de Tareas**: Sistema completo de estados (pendiente, en progreso, completada)
- **🎯 IDs Visibles**: Identificación fácil de tareas para edición rápida

## 🤝 Contribuir

1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crea un Pull Request

## 📄 Licencia

Este proyecto está bajo la Licencia MIT. Ver el archivo `LICENSE` para más detalles.

## 🆘 Soporte

Si tienes problemas o preguntas:

1. Revisa la sección de [Solución de Problemas](#-solución-de-problemas)
2. Abre un issue en GitHub
3. Contacta al equipo de desarrollo

## 📈 Roadmap

### ✅ Completado en v2.0.0
- [x] Migración a SQLite con índices optimizados
- [x] Comandos de edición y gestión avanzada
- [x] Reportes detallados por día, semana y mes
- [x] Búsqueda semántica con múltiples filtros
- [x] Exportación a CSV y JSON
- [x] Estados de tareas completos
- [x] IDs visibles para edición rápida
- [x] Migración automática de datos

### 🚀 Próximas Funcionalidades
- [ ] Integración directa con API de workflow
- [ ] Interfaz web para gestión de tareas
- [ ] Sincronización en tiempo real
- [ ] Analytics avanzados y métricas
- [ ] Integración con otros sistemas de gestión de tiempo
- [ ] Timer integrado para tracking en tiempo real
- [ ] Recordatorios y notificaciones
- [ ] Integración con Jira, GitHub Issues

---

**¡Disfruta usando workflow CLI! 🌾** 