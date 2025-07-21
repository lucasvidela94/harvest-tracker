# 🌾 Harvest CLI

Una herramienta de línea de comandos moderna y eficiente para gestionar tareas y reportes de tiempo, diseñada para integrarse con Harvest.

> **✨ Proyecto completamente migrado a Go** - Mejor rendimiento, mantenibilidad y distribución multiplataforma.

## 🚀 Instalación Rápida

### Instalación Automática (Recomendada)

```bash
# Clonar el repositorio
git clone https://github.com/lucasvidela94/harvest-tracker.git
cd harvest-tracker

# Instalar usando el script automático
./install.sh
```

### Instalación Manual

```bash
# Compilar e instalar
make install-script

# O manualmente
make build
make install
```

## 📋 Uso

Una vez instalado, puedes usar `harvest` desde cualquier lugar:

```bash
# Ver ayuda
harvest --help

# Agregar una tarea
harvest add "Desarrollar nueva funcionalidad" 4.0

# Ver estado actual
harvest status

# Generar reporte para Harvest
harvest report

# Actualizar a la última versión
harvest upgrade
```

## 🚀 Flujo de Trabajo Diario - Un Día en la Vida de un Dev

### 🌅 **Mañana (9:00 AM) - Planificación del Día**

```bash
# Ver qué tareas quedaron pendientes de ayer
harvest list --date 2025-07-20

# Agregar el daily standup
harvest daily

# Agregar tareas planificadas para hoy
harvest add "Revisar PRs pendientes" 1.0
harvest add "Desarrollar feature de login" 4.0
harvest add "Reunión de planning semanal" 1.5

# Ver el estado inicial del día
harvest status
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
harvest complete 1

# Agregar una tarea que surgió (bug fix urgente)
harvest add "Fix bug crítico en producción" 2.0

# Ver estado actualizado
harvest status
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
harvest list

# Editar la tarea del bug para ser más específica
harvest edit 4 --description "Fix bug crítico en API de usuarios" --hours 1.5

# Buscar tareas similares para referencia
harvest search "bug"
```

### 🌆 **Tarde (3:00 PM) - Progreso y Nuevas Tareas**

```bash
# Completar el bug fix
harvest complete 4

# Agregar tarea que surgió durante el desarrollo
harvest add "Documentar nueva API" 1.0

# Duplicar tarea de mañana para mañana (recurrente)
harvest duplicate 1 --tomorrow

# Ver progreso del día
harvest status
```

### 🌙 **Fin de Día (5:30 PM) - Cierre y Reporte**

```bash
# Completar tareas pendientes
harvest complete 2
harvest complete 3

# Ver reporte final del día
harvest report

# Generar reporte para Harvest (formato legacy)
harvest report --harvest
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
harvest report --week

# Exportar datos de la semana para análisis
harvest export --format csv --week --output semana-actual.csv

# Buscar tareas técnicas de la semana
harvest search --category tech --week

# Ver tareas completadas vs pendientes
harvest report --status completed --week
harvest report --status pending --week
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
- No necesitas abrir Harvest durante el día
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
- Formato legacy para copiar a Harvest
- Estadísticas y métricas automáticas

### ✅ **Flexibilidad Total**
- Agregar tareas para fechas pasadas/futuras
- Editar tareas existentes sin perder tiempo
- Reorganizar y ajustar según evoluciona el día
- Migración automática de datos antiguos

## 🛠️ Comandos Disponibles

### 📝 Gestión de Tareas
- `harvest add <descripción> <horas>` - Agregar nueva tarea
- `harvest add --date 2025-07-20 <descripción> <horas>` - Agregar tarea para fecha específica
- `harvest add --yesterday <descripción> <horas>` - Agregar tarea para ayer
- `harvest add --tomorrow <descripción> <horas>` - Agregar tarea para mañana
- `harvest tech <descripción> <horas>` - Agregar tarea técnica
- `harvest meeting <descripción> <horas>` - Agregar reunión
- `harvest qa <descripción> <horas>` - Agregar tarea de QA
- `harvest daily` - Agregar daily standup (automático)

### ✏️ Edición y Gestión
- `harvest edit <id> --description "nueva descripción"` - Editar tarea existente
- `harvest edit <id> --hours 2.5` - Cambiar horas de tarea
- `harvest edit <id> --category tech` - Cambiar categoría
- `harvest delete <id>` - Eliminar tarea
- `harvest duplicate <id>` - Duplicar tarea
- `harvest duplicate <id> --tomorrow` - Duplicar tarea para mañana
- `harvest complete <id>` - Marcar tarea como completada

### 📊 Información y Reportes
- `harvest status` - Ver estado actual de tareas
- `harvest list` - Listar tareas con IDs visibles
- `harvest list --date 2025-07-20` - Listar tareas de fecha específica
- `harvest report` - Reporte detallado de hoy
- `harvest report --week` - Reporte semanal
- `harvest report --month` - Reporte mensual
- `harvest report --date 2025-07-20` - Reporte de fecha específica
- `harvest report --category tech` - Reporte filtrado por categoría
- `harvest report --status completed` - Reporte de tareas completadas
- `harvest report --harvest` - Formato legacy para Harvest app

### 🔍 Búsqueda y Filtros
- `harvest search "texto"` - Buscar tareas por texto
- `harvest search --category tech` - Buscar por categoría
- `harvest search --status pending` - Buscar por estado
- `harvest search --date 2025-07-20` - Buscar por fecha

### 📤 Exportación
- `harvest export --format csv` - Exportar a CSV
- `harvest export --format json` - Exportar a JSON
- `harvest export --week --format csv` - Exportar semana a CSV
- `harvest export --category tech --format csv` - Exportar tareas técnicas

### 🔄 Migración y Sistema
- `harvest migrate` - Migrar datos de JSON a SQLite
- `harvest migrate --dry-run` - Simular migración
- `harvest migrate --backup-only` - Solo crear backup
- `harvest upgrade` - Actualizar a la última versión
- `harvest rollback` - Gestionar rollbacks

## ⚙️ Configuración

El CLI se configura automáticamente en `~/.harvest/`:

- `config.json` - Configuración general
- `tasks.db` - Base de datos SQLite con todas las tareas
- `tasks.json.backup.*` - Backups automáticos de datos JSON (si migraste)

### Migración de Datos

Si tienes datos en el formato JSON anterior, la migración es automática:

```bash
# Migrar datos existentes a SQLite
harvest migrate

# Simular migración sin cambios
harvest migrate --dry-run

# Solo crear backup
harvest migrate --backup-only
```

## 🔄 Actualizaciones

El sistema incluye un sistema de upgrade automático:

```bash
# Verificar actualizaciones
harvest upgrade
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

### El comando `harvest` no funciona

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
harvest --version
harvest --help
```

## 🔧 Desarrollo

### Compilar desde código fuente

```bash
# Clonar repositorio
git clone https://github.com/lucasvidela94/harvest-tracker.git
cd harvest-tracker

# Instalar dependencias
go mod tidy

# Compilar
go build -o harvest ./cmd/harvest

# Ejecutar
./harvest --help
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
harvest/
├── cmd/harvest/          # Punto de entrada principal
├── internal/             # Lógica interna del proyecto
│   ├── cli/             # Comandos CLI
│   ├── core/            # Lógica principal
│   └── upgrade/         # Sistema de upgrade
├── pkg/harvest/         # Tipos y utilidades
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
└── harvest              # Ejecutable compilado
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
- [ ] Integración directa con API de Harvest
- [ ] Interfaz web para gestión de tareas
- [ ] Sincronización en tiempo real
- [ ] Analytics avanzados y métricas
- [ ] Integración con otros sistemas de gestión de tiempo
- [ ] Timer integrado para tracking en tiempo real
- [ ] Recordatorios y notificaciones
- [ ] Integración con Jira, GitHub Issues

---

**¡Disfruta usando Harvest CLI! 🌾** 