# 🌾 Harvest CLI (Go Version)

Sistema de tracking de tareas para Harvest, reescrito en Go para mejor performance y distribución.

## 🚀 Características

- ✅ **Binario standalone** - Sin dependencias externas
- ✅ **Cross-platform** - Linux, macOS, Windows
- ✅ **Performance** - Ejecución rápida
- ✅ **Compatibilidad** - Usa los mismos archivos de datos que la versión Python
- ✅ **Migración gradual** - Puede coexistir con la versión Python

## 📦 Instalación

### Desarrollo Local

```bash
# Clonar el repositorio
git clone https://github.com/lucasvidela94/harvest-tracker.git
cd harvest-tracker

# Construir el binario
make build

# Ejecutar
./harvest
```

### Build para Múltiples Plataformas

```bash
# Construir para todas las plataformas
make build-all

# Los binarios se crean en build/
# - harvest-linux-amd64
# - harvest-linux-arm64
# - harvest-darwin-amd64
# - harvest-darwin-arm64
# - harvest-windows-amd64.exe
```

## 🏗️ Estructura del Proyecto

```
harvest-go/
├── cmd/
│   └── harvest/
│       └── main.go          # Punto de entrada
├── internal/
│   ├── core/
│   │   └── config.go        # Gestión de configuración
│   ├── cli/
│   │   └── commands.go      # Comandos CLI con Cobra
│   └── upgrade/             # Sistema de actualización (pendiente)
├── pkg/
│   └── harvest/
│       └── types.go         # Tipos de datos
├── Makefile                 # Comandos de build
├── go.mod                   # Dependencias
└── README.md
```

## 🎯 Estado Actual

### ✅ Implementado
- [x] Estructura básica del proyecto
- [x] Framework CLI con Cobra
- [x] Gestión de configuración
- [x] Tipos de datos
- [x] Build multi-plataforma
- [x] Comando de versión

### 🚧 En Desarrollo
- [ ] Gestión de tareas (TaskManager)
- [ ] Comandos add/tech/meeting/qa
- [ ] Comando status
- [ ] Comando report
- [ ] Sistema de upgrade

### 📋 Pendiente
- [ ] Tests unitarios
- [ ] Documentación completa
- [ ] Script de instalación
- [ ] CI/CD pipeline

## 🔧 Desarrollo

### Comandos Make

```bash
make build      # Construir binario
make clean      # Limpiar archivos
make test       # Ejecutar tests
make deps       # Instalar dependencias
make build-all  # Build multi-plataforma
make run        # Construir y ejecutar
make check      # Verificar código
make help       # Mostrar ayuda
```

### Dependencias

```go
require (
    github.com/spf13/cobra v1.9.1  // CLI framework
)
```

## 🔄 Migración desde Python

### Estrategia de Migración Gradual

1. **Fase 1**: Estructura básica ✅
2. **Fase 2**: Comandos core (add, status, report)
3. **Fase 3**: Sistema de upgrade
4. **Fase 4**: Tests y documentación
5. **Fase 5**: Distribución y CI/CD

### Compatibilidad de Datos

- Usa los mismos archivos de configuración (`~/.harvest/config.json`)
- Usa los mismos archivos de datos (`~/.harvest/tasks.json`)
- Puede coexistir con la versión Python durante la migración

## 🎯 Ventajas sobre Python

| Aspecto | Python | Go |
|---------|--------|----|
| **Dependencias** | Python 3.x + librerías | Solo binario |
| **Distribución** | Script + archivos | Un archivo |
| **Performance** | ⭐⭐ | ⭐⭐⭐⭐⭐ |
| **Instalación** | Compleja | Simple |
| **Cross-platform** | Depende de Python | Nativo |

## 🤝 Contribuir

1. Fork el repositorio
2. Crea una rama para tu feature (`git checkout -b feature/nueva-funcionalidad`)
3. Commit tus cambios (`git commit -am 'Agregar nueva funcionalidad'`)
4. Push a la rama (`git push origin feature/nueva-funcionalidad`)
5. Crea un Pull Request

## 📄 Licencia

MIT License - ver [LICENSE](LICENSE) para detalles.

---

**Nota**: Esta es la versión Go de Harvest CLI. La versión Python original sigue funcionando y se puede encontrar en el directorio `scripts/harvest/`. 