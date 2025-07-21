# 🌾 Harvest CLI - Scope y Funcionalidades

## 📋 Funcionalidades Actuales

### 🎯 Comandos Principales

#### Gestión de Tareas
- **`harvest add <descripción> <horas> [categoría]`** - Agregar nueva tarea
  - Ejemplo: `harvest add "Fix bug" 2.0`
  - Ejemplo: `harvest add "Development" 3.5 tech`

- **`harvest tech <descripción> <horas>`** - Agregar tarea técnica
  - Ejemplo: `harvest tech "Fix bug" 2.0`
  - Categoría automática: `tech` 💻

- **`harvest meeting <descripción> <horas>`** - Agregar reunión
  - Ejemplo: `harvest meeting "Team sync" 1.0`
  - Categoría automática: `meeting` 🤝

- **`harvest qa <descripción> <horas>`** - Agregar tarea de QA
  - Ejemplo: `harvest qa "Testing" 1.5`
  - Categoría automática: `qa` 🧪

- **`harvest daily`** - Agregar daily standup automático
  - Duración configurable (por defecto: 0.25h)
  - Categoría automática: `daily` 📢

#### Información y Reportes
- **`harvest status`** - Ver estado actual de tareas
  - Muestra fecha actual
  - Horas trabajadas vs objetivo
  - Horas restantes o overtime
  - Lista de tareas con iconos
  - Barra de progreso visual

- **`harvest report`** - Generar reporte para Harvest
  - Formato: "Descripción - X.Xh"
  - Listo para copiar y pegar en Harvest
  - Copia automática al portapapeles (Linux)

#### Sistema y Mantenimiento
- **`harvest upgrade`** - Actualizar a la última versión
  - Verificación de actualizaciones disponibles
  - Backup automático antes de actualizar
  - Migración Python → Go
  - Sistema de rollback

- **`harvest rollback`** - Gestionar rollbacks
  - Verificar disponibilidad de rollback
  - Información de backups
  - Logs de actividad de rollback

- **`harvest version`** - Mostrar información de versión
  - Versión actual
  - Información de build

### ⚙️ Configuración

#### Archivos de Configuración
- **`~/.harvest/config.json`** - Configuración general
  - `daily_hours_target`: Objetivo de horas diarias (8.0h por defecto)
  - `daily_standup_hours`: Horas del daily standup (0.25h por defecto)
  - `data_file`: Ruta del archivo de datos
  - `user_name`: Nombre del usuario
  - `company`: Empresa
  - `timezone`: Zona horaria

- **`~/.harvest/tasks.json`** - Datos de tareas
  - Almacenamiento en formato JSON
  - Compatibilidad con versiones anteriores

#### Categorías de Tareas
- **`tech`** 💻 - Tareas técnicas/desarrollo
- **`meeting`** 🤝 - Reuniones
- **`qa`** 🧪 - Testing/QA
- **`doc`** 📚 - Documentación
- **`planning`** 📋 - Planificación
- **`research`** 🔍 - Investigación
- **`review`** 👀 - Revisión de código
- **`deploy`** 🚀 - Despliegue
- **`daily`** 📢 - Daily standup
- **`general`** 📝 - Tareas generales

### 🛡️ Características de Seguridad

#### Sistema de Backup
- Backup automático antes de actualizaciones
- Backup del binario ejecutable
- Backup de datos de configuración
- Verificación de integridad de backups

#### Sistema de Rollback
- Rollback automático en caso de fallo
- Gestión de versiones anteriores
- Logs detallados de operaciones
- Protección contra pérdida de datos

### 🖥️ Compatibilidad

#### Plataformas Soportadas
- **Linux**: amd64, arm64
- **macOS**: amd64, arm64
- **Windows**: amd64

#### Características Técnicas
- Escrito en Go para alto rendimiento
- Instalación automática con scripts
- Distribución multiplataforma
- Sistema de actualizaciones integrado

## 🚨 Problemas del Flujo Diario Actual

### 🗄️ Limitaciones del Almacenamiento JSON
- **Sin índices**: Búsquedas lentas en archivos grandes
- **Sin consultas complejas**: No se pueden hacer filtros avanzados
- **Sin transacciones**: Riesgo de corrupción de datos
- **Sin relaciones**: No se pueden relacionar tareas con proyectos/clientes
- **Escalabilidad limitada**: El archivo crece indefinidamente
- **Sin backup incremental**: Solo backup completo del archivo

### 🆔 Gestión de IDs
- **IDs no visibles**: El usuario no ve los IDs de las tareas
- **Sin forma de referenciar**: No hay forma natural de identificar tareas
- **Sin listado con IDs**: El comando `status` no muestra IDs
- **Sin búsqueda por ID**: No se puede buscar una tarea específica
- **Sin confirmación visual**: No hay forma de verificar qué tarea se está editando

### 📅 Gestión de Fechas
- **Solo fecha actual**: No se pueden agregar tareas para fechas pasadas
- **Sin flag `--date`**: No hay forma de especificar una fecha específica
- **Sin retroactividad**: Si te olvidas de registrar el lunes, no puedes hacerlo el martes
- **Sin planificación**: No se pueden agregar tareas para fechas futuras

### ✏️ Edición y Modificación
- **Sin edición**: No hay forma de modificar tareas existentes
- **Sin eliminación**: No se pueden eliminar tareas incorrectas
- **Sin ajuste de horas**: Si la hora estimada fue incorrecta, no se puede cambiar
- **Sin corrección de descripción**: Errores tipográficos no se pueden corregir

### 📊 Gestión de Estado
- **Sin estados**: No hay estados como "pendiente", "en progreso", "completada"
- **Sin tracking**: No se puede marcar tareas como completadas
- **Sin progreso**: No hay forma de ver qué tareas están terminadas
- **Sin resumen**: No hay resumen de tareas completadas vs pendientes

### 🔍 Búsqueda y Filtrado
- **Sin búsqueda**: No se pueden buscar tareas por texto
- **Sin filtros**: No hay filtros por fecha, categoría, estado
- **Sin historial**: No hay forma de ver tareas de días anteriores
- **Sin exportación selectiva**: No se pueden exportar tareas específicas

## 🔧 Soluciones para el Flujo Diario

### 🗄️ Migración a SQLite

#### Ventajas de SQLite
- **Consultas complejas**: Filtros, ordenamiento, agrupación
- **Índices**: Búsquedas rápidas por fecha, categoría, texto
- **Transacciones**: Integridad de datos garantizada
- **Escalabilidad**: Manejo eficiente de miles de tareas
- **Backup incremental**: Solo cambios, no archivo completo
- **Relaciones**: Futura integración con proyectos/clientes

#### Estructura de Base de Datos
```sql
-- Tabla de tareas
CREATE TABLE tasks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    description TEXT NOT NULL,
    hours REAL NOT NULL,
    category TEXT DEFAULT 'general',
    date TEXT NOT NULL,
    status TEXT DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Índices para búsquedas rápidas
CREATE INDEX idx_tasks_date ON tasks(date);
CREATE INDEX idx_tasks_category ON tasks(category);
CREATE INDEX idx_tasks_status ON tasks(status);
CREATE INDEX idx_tasks_description ON tasks(description);
```

#### Consultas SQL Posibles
```sql
-- Tareas de hoy con estado
SELECT id, description, hours, category, status 
FROM tasks 
WHERE date = '2025-07-21' 
ORDER BY created_at;

-- Tareas pendientes de la semana
SELECT id, description, hours, category, date 
FROM tasks 
WHERE status = 'pending' 
AND date BETWEEN '2025-07-21' AND '2025-07-25';

-- Búsqueda por texto
SELECT id, description, hours, category, date 
FROM tasks 
WHERE description LIKE '%bug%' 
ORDER BY date DESC;

-- Resumen por categoría
SELECT category, COUNT(*) as count, SUM(hours) as total_hours 
FROM tasks 
WHERE date = '2025-07-21' 
GROUP BY category;

-- Tareas más largas
SELECT id, description, hours, category 
FROM tasks 
WHERE hours > 4.0 
ORDER BY hours DESC;
```

### 🆔 Gestión Inteligente de IDs

#### Comando `list` con IDs
```bash
harvest list
# Output:
# 📅 Today's Tasks (2025-07-21):
# [1] 📝 Tarea 1 - 2.0h (general)
# [2] 📝 Tarea 2 - 1.5h (general)
# [3] 💻 Desarrollo feature X - 4.0h (tech)
```

#### Comando `status` mejorado
```bash
harvest status
# Output:
# 📅 Today (2025-07-21): 7.5h / 8.0h
# 📈 Remaining: 0.5h
# 
# 📝 Today's tasks:
# [1] 📝 Tarea 1 - 2.0h (general) ✅
# [2] 📝 Tarea 2 - 1.5h (general) 🔄
# [3] 💻 Desarrollo feature X - 4.0h (tech) ⏳
```

#### Búsqueda por texto con IDs
```bash
harvest search "feature"
# Output:
# 🔍 Search results for "feature":
# [3] 💻 Desarrollo feature X - 4.0h (tech) - 2025-07-21
# [7] 💻 Fix feature Y - 2.0h (tech) - 2025-07-20
```

#### Confirmación visual en edición
```bash
harvest edit 3
# Output:
# ✏️ Editing task:
# [3] 💻 Desarrollo feature X - 4.0h (tech) - 2025-07-21
# 
# New description (or Enter to keep current): 
# New hours (or Enter to keep current): 
# New category (or Enter to keep current): 
```

### 📅 Gestión de Fechas Avanzada
- **`harvest add --date 2025-07-20 "Tarea del lunes" 3.0`** - Agregar tarea para fecha específica
- **`harvest add --yesterday "Tarea olvidada" 2.0`** - Agregar tarea para ayer
- **`harvest add --tomorrow "Planificación" 1.5`** - Agregar tarea para mañana
- **`harvest add --week`** - Agregar tarea para toda la semana
- **`harvest status --date 2025-07-20`** - Ver estado de fecha específica
- **`harvest status --week`** - Ver resumen de la semana

### ✏️ Edición y Modificación de Tareas
- **`harvest edit <id>`** - Modificar tarea existente (interactivo)
- **`harvest edit <id> --description "Nueva descripción"`** - Cambiar descripción
- **`harvest edit <id> --hours 1.5`** - Ajustar horas
- **`harvest edit <id> --category tech`** - Cambiar categoría
- **`harvest delete <id>`** - Eliminar tarea
- **`harvest delete --date 2025-07-20`** - Eliminar todas las tareas de una fecha
- **`harvest duplicate <id>`** - Duplicar tarea

### 📊 Gestión de Estados
- **`harvest complete <id>`** - Marcar tarea como completada ✅
- **`harvest start <id>`** - Marcar tarea como en progreso 🔄
- **`harvest pause <id>`** - Pausar tarea ⏸️
- **`harvest status --completed`** - Ver solo tareas completadas
- **`harvest status --pending`** - Ver solo tareas pendientes
- **`harvest status --progress`** - Ver tareas en progreso

### 🔍 Búsqueda y Filtrado
- **`harvest search "bug"`** - Buscar tareas por texto
- **`harvest list --date 2025-07-20`** - Listar tareas de fecha específica
- **`harvest list --category tech`** - Listar tareas por categoría
- **`harvest list --week`** - Listar tareas de la semana
- **`harvest list --month`** - Listar tareas del mes
- **`harvest history`** - Ver historial de tareas

### 📋 Reportes Mejorados
- **`harvest report --date 2025-07-20`** - Reporte de fecha específica
- **`harvest report --week`** - Reporte semanal
- **`harvest report --completed`** - Solo tareas completadas
- **`harvest report --pending`** - Solo tareas pendientes
- **`harvest report --format csv`** - Exportar en formato CSV
- **`harvest report --format excel`** - Exportar en formato Excel

### 🔄 Migración de Datos JSON → SQLite

#### Proceso de Migración Automática
```bash
# El sistema detecta automáticamente si hay datos JSON
harvest migrate
# Output:
# 🔄 Detected JSON data file
# 📊 Found 45 tasks in JSON format
# 🗄️ Migrating to SQLite...
# ✅ Migration completed successfully!
# 📁 JSON backup saved to: ~/.harvest/tasks.json.backup
```

#### Comandos de Migración
- **`harvest migrate`** - Migración automática JSON → SQLite
- **`harvest migrate --dry-run`** - Simular migración sin cambios
- **`harvest migrate --backup`** - Crear backup antes de migrar
- **`harvest migrate --rollback`** - Revertir a JSON si hay problemas
- **`harvest export --format json`** - Exportar SQLite a JSON
- **`harvest import --format json`** - Importar JSON a SQLite

#### Estructura de Archivos Post-Migración
```
~/.harvest/
├── config.json          # Configuración (mantiene JSON)
├── tasks.db             # Base de datos SQLite (nuevo)
├── tasks.json.backup    # Backup del JSON original
└── migration.log        # Log de la migración
```

## 📅 Flujos de Trabajo del Día a Día

### 🌅 Flujo Matutino
```bash
# Ver qué hay pendiente de ayer
harvest status --yesterday

# Agregar daily standup
harvest daily

# Planificar tareas del día
harvest add "Revisar emails" 0.5
harvest add "Desarrollo feature X" 4.0 tech
harvest add "Reunión de equipo" 1.0 meeting
```

### 🌆 Flujo Vespertino
```bash
# Ver progreso del día
harvest status

# Marcar tareas completadas
harvest complete 1
harvest complete 3

# Ajustar horas si fue diferente
harvest edit 2 --hours 3.5

# Generar reporte para Harvest
harvest report
```

### 📅 Flujo de Recuperación (Olvidos)
```bash
# Martes: "Me olvidé de registrar el lunes"
harvest add --yesterday "Desarrollo feature A" 3.0 tech
harvest add --yesterday "Reunión de planning" 1.5 meeting
harvest status --yesterday

# Viernes: "Me olvidé de toda la semana"
harvest add --date 2025-07-21 "Tarea del lunes" 2.0
harvest add --date 2025-07-22 "Tarea del martes" 1.5
harvest add --date 2025-07-23 "Tarea del miércoles" 3.0
harvest add --date 2025-07-24 "Tarea del jueves" 2.5
harvest report --week
```

### 🔄 Flujo de Corrección
```bash
# Error tipográfico
harvest edit 1 --description "Fix bug in login system"

# Hora incorrecta
harvest edit 2 --hours 1.5

# Categoría incorrecta
harvest edit 3 --category qa

# Tarea que no va
harvest delete 4

# Duplicar tarea para otro día
harvest duplicate 1 --date 2025-07-22
```

### 📊 Flujo de Análisis
```bash
# Ver productividad de la semana
harvest status --week

# Buscar tareas específicas
harvest search "bug"

# Ver solo tareas técnicas
harvest list --category tech

# Ver tareas completadas
harvest status --completed

# Exportar reporte semanal
harvest report --week --format csv
```

## 🚀 Mejoras Propuestas para Futuras Versiones

### 🔌 Integración con APIs

#### Integración Directa con Harvest
- **API de Harvest**: Conexión directa con la API de Harvest
- **Sincronización automática**: Envío automático de tareas
- **Autenticación OAuth**: Login seguro con Harvest
- **Proyectos y clientes**: Selección de proyectos desde CLI

#### Integración con Otros Sistemas
- **Jira**: Sincronización con tickets de Jira
- **GitHub**: Integración con issues y PRs
- **Slack**: Notificaciones automáticas
- **Trello**: Sincronización con tarjetas

### 📊 Reportes Avanzados

#### Analytics y Métricas
- **Reportes semanales/mensuales**: Análisis de productividad
- **Gráficos de tiempo**: Visualización de distribución de tiempo
- **Métricas de eficiencia**: KPIs de productividad
- **Exportación a Excel/CSV**: Reportes en múltiples formatos

#### Reportes Personalizados
- **Filtros por fecha**: Reportes de rangos específicos
- **Filtros por categoría**: Análisis por tipo de tarea
- **Filtros por proyecto**: Reportes por proyecto/cliente
- **Templates personalizables**: Formatos de reporte configurables

### 🎨 Interfaz de Usuario

#### Interfaz Web
- **Dashboard web**: Interfaz gráfica para gestión
- **Visualización de datos**: Gráficos y estadísticas
- **Gestión de tareas**: Interfaz drag & drop
- **Configuración avanzada**: Panel de configuración web

#### Mejoras en CLI
- **Modo interactivo**: Interfaz conversacional
- **Autocompletado**: Sugerencias inteligentes
- **Aliases personalizables**: Comandos abreviados
- **Temas visuales**: Personalización de colores

### 🔄 Funcionalidades Avanzadas

#### Gestión de Tiempo
- **Timer integrado**: Cronómetro para tareas en tiempo real
- **Pomodoro**: Integración con técnica Pomodoro
- **Recordatorios**: Notificaciones de descansos
- **Tracking automático**: Detección de actividad

#### Colaboración
- **Equipos**: Gestión de equipos y roles
- **Compartir reportes**: Envío automático a stakeholders
- **Aprobaciones**: Flujo de aprobación de reportes
- **Comentarios**: Sistema de comentarios en tareas

#### Automatización
- **Webhooks**: Integración con sistemas externos
- **Scripts personalizados**: Automatización de flujos
- **Triggers**: Acciones automáticas basadas en eventos
- **Sincronización en tiempo real**: Actualizaciones instantáneas

### 🛠️ Mejoras Técnicas

#### Rendimiento
- **Base de datos local**: Migración de JSON a SQLite
- **Caché inteligente**: Optimización de consultas
- **Compresión de datos**: Reducción de tamaño de archivos
- **Indexación**: Búsqueda rápida de tareas

#### Seguridad
- **Encriptación**: Cifrado de datos sensibles
- **Autenticación local**: Login con contraseña
- **Auditoría**: Logs detallados de todas las operaciones
- **Backup en la nube**: Sincronización con servicios cloud

#### Arquitectura
- **Plugins**: Sistema de plugins extensible
- **APIs internas**: APIs para integraciones
- **Microservicios**: Arquitectura modular
- **Contenedores**: Distribución con Docker

### 📱 Experiencia de Usuario

#### Accesibilidad
- **Modo oscuro**: Tema oscuro para CLI y web
- **Accesibilidad**: Soporte para lectores de pantalla
- **Internacionalización**: Múltiples idiomas
- **Responsive**: Interfaz adaptativa

#### Personalización
- **Templates de tareas**: Plantillas reutilizables
- **Workflows personalizados**: Flujos de trabajo configurables
- **Atajos de teclado**: Navegación rápida
- **Preferencias avanzadas**: Configuración granular

### 🔍 Funcionalidades de Búsqueda y Filtrado

#### Búsqueda Inteligente
- **Búsqueda por texto**: Búsqueda en descripciones
- **Búsqueda por fecha**: Filtros temporales
- **Búsqueda por categoría**: Filtros por tipo
- **Búsqueda por horas**: Filtros por duración

#### Organización
- **Tags**: Sistema de etiquetas
- **Proyectos**: Organización por proyectos
- **Prioridades**: Sistema de prioridades
- **Estados**: Estados de tareas (pendiente, en progreso, completada)

### 📈 Roadmap Sugerido

#### Versión 2.1 - Migración a SQLite ✅ COMPLETADO
- [x] **Migración JSON → SQLite** con backup automático
- [x] **Comando `migrate`** para migración de datos
- [x] **Comando `list`** con IDs visibles
- [x] **Comando `status`** mejorado con IDs y estados
- [x] **Comando `edit`** con confirmación visual
- [x] **Comando `delete`** con confirmación

#### Versión 2.2 - Flujo Diario Básico ✅ COMPLETADO
- [x] **Flag `--date`** para agregar tareas en fechas específicas
- [x] **Flag `--yesterday`** y `--tomorrow` para fechas relativas
- [x] **Comando `complete`** para marcar tareas como completadas
- [x] **Comando `search`** para buscar tareas por texto
- [x] **Estados de tareas** (pendiente, en progreso, completada)
- [x] **Filtros por categoría** en status y list

#### Versión 2.2 - Gestión Avanzada ✅ COMPLETADO
- [x] **Comando `search`** para buscar tareas por texto
- [x] **Filtros por categoría** en status y list
- [x] **Comando `duplicate`** para duplicar tareas
- [x] **Estados de tareas** (pendiente, en progreso, completada)
- [x] **Reportes por fecha** y rangos de fechas
- [x] **Exportación en múltiples formatos** (CSV, JSON)

#### Versión 2.3 - Integración Básica
- [ ] Integración con API de Harvest
- [ ] Autenticación OAuth
- [ ] Sincronización básica de tareas
- [ ] Reportes mejorados

#### Versión 2.2 - Interfaz Web
- [ ] Dashboard web básico
- [ ] Visualización de datos
- [ ] Gestión de tareas web
- [ ] Configuración web

#### Versión 2.3 - Automatización
- [ ] Timer integrado
- [ ] Recordatorios
- [ ] Webhooks básicos
- [ ] Integración con Jira

#### Versión 2.4 - Colaboración
- [ ] Gestión de equipos
- [ ] Compartir reportes
- [ ] Sistema de comentarios
- [ ] Aprobaciones

#### Versión 3.0 - Plataforma Completa
- [ ] Base de datos local
- [ ] Sistema de plugins
- [ ] APIs internas
- [ ] Contenedores Docker

---

**Nota**: Este documento se actualiza regularmente según las necesidades del proyecto y feedback de los usuarios. 