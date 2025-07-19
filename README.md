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

## 🛠️ Comandos Disponibles

### Gestión de Tareas
- `harvest add <descripción> <horas>` - Agregar nueva tarea
- `harvest tech <descripción> <horas>` - Agregar tarea técnica
- `harvest meeting <descripción> <horas>` - Agregar reunión
- `harvest qa <descripción> <horas>` - Agregar tarea de QA
- `harvest daily` - Agregar daily standup (automático)

### Información y Reportes
- `harvest status` - Ver estado actual de tareas
- `harvest report` - Generar reporte para Harvest

### Sistema
- `harvest upgrade` - Actualizar a la última versión
- `harvest rollback` - Gestionar rollbacks

## ⚙️ Configuración

El CLI se configura automáticamente en `~/.harvest/`:

- `config.json` - Configuración general
- `tasks.json` - Datos de tareas

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

- **⚡ Alto Rendimiento**: Escrito en Go para máxima velocidad
- **🔧 Fácil Instalación**: Scripts automáticos de instalación
- **🔄 Actualizaciones Automáticas**: Sistema de upgrade integrado
- **🛡️ Seguridad**: Backup y rollback automáticos
- **📱 Multiplataforma**: Soporte para Linux, macOS y Windows
- **📊 Reportes Inteligentes**: Generación automática de reportes para Harvest

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

- [ ] Integración directa con API de Harvest
- [ ] Interfaz web para gestión de tareas
- [ ] Sincronización en tiempo real
- [ ] Reportes avanzados y analytics
- [ ] Integración con otros sistemas de gestión de tiempo

---

**¡Disfruta usando Harvest CLI! 🌾** 