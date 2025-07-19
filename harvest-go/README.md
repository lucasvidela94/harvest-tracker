# 🌾 Harvest CLI

Una herramienta de línea de comandos para gestionar tareas y reportes de tiempo, diseñada para integrarse con Harvest.

## 🚀 Instalación Rápida

### Opción 1: Instalación Automática (Recomendada)

```bash
# Clonar el repositorio
git clone https://github.com/lucasvidela94/harvest-tracker.git
cd harvest-tracker/harvest-go

# Instalar usando el script automático
./install.sh
```

### Opción 2: Instalación Manual

```bash
# Compilar e instalar
make install-script

# O manualmente
make build
make install
```

### Opción 3: Desde el código fuente

```bash
# Clonar y compilar
git clone https://github.com/lucasvidela94/harvest-tracker.git
cd harvest-tracker/harvest-go
go build -o harvest ./cmd/harvest

# Mover a PATH
sudo mv harvest /usr/local/bin/
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

### Configuración de Daily Standup

```bash
# Configurar horas del daily (por defecto: 0.25h)
# Se puede modificar en ~/.harvest/config.json
```

## 🔄 Actualizaciones

El sistema incluye un sistema de upgrade automático:

```bash
# Verificar actualizaciones
harvest upgrade

# El sistema:
# 1. Detecta la versión actual
# 2. Crea backup automático
# 3. Descarga nueva versión
# 4. Instala y migra datos
# 5. Proporciona rollback automático
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

### Problemas de permisos

```bash
# Hacer ejecutable
chmod +x ~/.local/bin/harvest
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
cd harvest-tracker/harvest-go

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
harvest-go/
├── cmd/harvest/          # Punto de entrada
├── internal/
│   ├── cli/             # Comandos CLI
│   ├── core/            # Lógica principal
│   └── upgrade/         # Sistema de upgrade
├── pkg/harvest/         # Tipos y utilidades
├── install.sh           # Script de instalación
├── uninstall.sh         # Script de desinstalación
└── Makefile             # Comandos de build
```

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

---

**¡Disfruta usando Harvest CLI! 🌾** 