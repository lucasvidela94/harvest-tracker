# 🌾 Harvest Scripts v1.0

Sistema simple y directo para gestionar tareas de Harvest desde la línea de comandos.

## 🚀 Instalación

### Instalación Automática (Recomendada)

```bash
# Instalar scripts
./install.sh
```

El script de instalación:
- ✅ Detecta automáticamente tu shell (zsh, bash, etc.)
- ✅ Crea la configuración en `~/.harvest/`
- ✅ Instala dependencias (pyperclip)
- ✅ Configura aliases permanentes en tu shell
- ✅ Recarga la configuración automáticamente

### Instalación Manual

Si prefieres configurar manualmente:

```bash
# Agregar a ~/.zshrc o ~/.bashrc:
alias harvest='~/scripts/harvest/harvest'
alias finish='~/scripts/harvest/finish'
alias week='~/scripts/harvest/week'

# Crear configuración
mkdir -p ~/.harvest
echo '{"daily_hours_target": 8.0, "daily_standup_hours": 0.25, "data_file": "~/.harvest/tasks.json"}' > ~/.harvest/config.json

# Instalar dependencias
pip3 install --user pyperclip
```

### Desinstalación

```bash
./uninstall.sh
```

## 📋 Comandos Principales

### `harvest` - Comando Principal
```bash
# Agregar tareas por categoría
harvest daily                    # Daily standup (0.25h)
harvest tech "Fix bug" 2.0       # Tarea técnica
harvest meeting "Sync" 1.0       # Reunión
harvest qa "Testing" 1.5         # QA/Testing
harvest add "Task" 1.0 doc       # Tarea genérica con categoría

# Ver estado y reportes
harvest status                   # Estado actual del día
harvest report                   # Generar reporte para Harvest
```

### `finish` - Completar el Día
```bash
finish                     # Modo interactivo
finish 2.0                 # Agregar tarea de 2h automáticamente
finish 1.5 "Final task"    # Agregar tarea específica
```

### `week` - Reportes Semanales
```bash
week                       # Mostrar reporte semanal
week copy                  # Copiar reporte al portapapeles
```

## 🎯 Flujo de Trabajo Diario

### 1. Inicio del día
```bash
harvest daily                    # Agregar daily standup
```

### 2. Durante el día
```bash
harvest tech "Development" 2.5   # Agregar tareas mientras trabajas
harvest meeting "Planning" 1.0   # Reuniones
harvest qa "Testing" 1.0         # Testing
harvest status                   # Verificar progreso
```

### 3. Final del día
```bash
finish                     # Completar horas restantes
harvest report                   # Generar reporte para Harvest
```

## 📊 Categorías Disponibles

- **`tech`** 💻 - Desarrollo, debugging, code review
- **`meeting`** 🤝 - Reuniones, planning, syncs
- **`qa`** 🧪 - Testing, bug fixes, validation
- **`doc`** 📚 - Documentación, research
- **`planning`** 📋 - Sprint planning, roadmap
- **`research`** 🔍 - Investigación, POCs
- **`review`** 👀 - Code review, PRs
- **`deploy`** 🚀 - Deployment, releases
- **`daily`** 📢 - Daily standup (0.25h automático)

## 💡 Características

- ✅ **Comandos secuenciales** - No más modo interactivo molesto
- ✅ **Categorías con iconos** - Fácil identificación visual
- ✅ **Barra de progreso** - Estado visual del día
- ✅ **Sugerencias inteligentes** - Basadas en horas restantes
- ✅ **Copia automática** - Reportes listos para Harvest
- ✅ **Completado inteligente** - `finish` para completar el día

## 📈 Ejemplo de Uso

```bash
# Día típico
harvest daily
harvest tech "Feature development" 3.0
harvest meeting "Team sync" 1.0
harvest status
harvest qa "Bug fixes" 1.5
finish 2.5 "Documentation"
harvest report
```

## 🔧 Configuración

La configuración se almacena en `~/.harvest/config.json`:

```json
{
    "daily_hours_target": 8.0,
    "daily_standup_hours": 0.25,
    "data_file": "~/.harvest/tasks.json",
    "user_name": "tu_usuario",
    "company": "",
    "timezone": "UTC"
}
```

- **Objetivo diario**: Configurable (por defecto 8 horas)
- **Daily standup**: Configurable (por defecto 0.25h)
- **Archivo de datos**: `~/.harvest/tasks.json`
- **Copia automática**: Requiere `pyperclip` (instalado automáticamente)

## 📱 Integración con Harvest

1. Ejecuta `harvest report` o `week copy`
2. Se copia automáticamente al portapapeles
3. Pega directamente en Harvest

Formato generado:
```
Daily Standup - 0.25h
Feature development - 3.0h
Team sync - 1.0h
Bug fixes - 1.5h
Documentation - 2.5h

Total: 8.25h
```

## 🎉 ¡Listo para usar!

El sistema está diseñado para ser simple, rápido y efectivo. Sin complicaciones, solo comandos directos que funcionan. 