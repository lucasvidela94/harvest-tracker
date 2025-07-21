# 🧪 Testing de workflow CLI

Este documento describe cómo ejecutar tests automatizados para workflow CLI usando Docker.

## 🚀 Inicio Rápido

### Prerrequisitos
- Docker instalado y funcionando
- Binario `workflow` compilado en el directorio raíz

### Ejecutar todos los tests
```bash
make -f Makefile.test test-run
```

### Ejecutar test rápido
```bash
make -f Makefile.test test-quick
```

## 📋 Comandos Disponibles

### Construir imagen de test
```bash
make -f Makefile.test test-build
```

### Ejecutar tests completos
```bash
make -f Makefile.test test-run
```

### Shell interactivo para testing
```bash
make -f Makefile.test test-shell
```

### Test específico de migración
```bash
make -f Makefile.test test-migration
```

### Limpiar recursos
```bash
make -f Makefile.test test-clean
```

### Ver ayuda
```bash
make -f Makefile.test test-help
```

## 🧪 Qué se Testea

### 1. **Funcionalidad Básica**
- ✅ Verificación de versión y ayuda
- ✅ Creación de directorios de configuración
- ✅ Comandos básicos (add, edit, complete, delete)

### 2. **Gestión de Tareas**
- ✅ Agregar tareas con diferentes categorías
- ✅ Editar tareas (horas, descripción, categoría)
- ✅ Completar tareas
- ✅ Eliminar tareas
- ✅ Duplicar tareas

### 3. **Fechas y Categorías**
- ✅ Tareas para fechas específicas (--date)
- ✅ Tareas para ayer (--yesterday)
- ✅ Tareas para mañana (--tomorrow)
- ✅ Categorías específicas (tech, meeting, qa, daily)

### 4. **Búsqueda y Filtros**
- ✅ Búsqueda por texto
- ✅ Filtros por categoría
- ✅ Filtros por estado

### 5. **Reportes y Exportación**
- ✅ Reportes diarios
- ✅ Reportes semanales
- ✅ Formato legacy para workflow
- ✅ Exportación a CSV
- ✅ Exportación a JSON

### 6. **Migración de Datos**
- ✅ Simulación de datos JSON viejos
- ✅ Migración dry-run
- ✅ Verificación de estructura de archivos

### 7. **Estructura de Archivos**
- ✅ Creación de directorio .workflow
- ✅ Archivo de configuración
- ✅ Base de datos SQLite

## 🔧 Personalización

### Agregar nuevos tests
Edita el archivo `test.sh` y agrega nuevos casos de prueba:

```bash
# Test personalizado
print_info "Test personalizado"
workflow add "Mi test" 1.0
check_exit "workflow add personalizado"
```

### Modificar Dockerfile
Edita `Dockerfile.test` para cambiar la imagen base o agregar dependencias:

```dockerfile
# Cambiar imagen base
FROM debian:bullseye

# Agregar dependencias
RUN apt-get update && apt-get install -y \
    sqlite3 \
    jq
```

### Ejecutar tests específicos
Puedes ejecutar solo partes del script:

```bash
# Solo tests básicos
docker run --rm workflow-test bash -c "
    workflow --version && \
    workflow add 'test' 1.0 && \
    workflow status
"
```

## 🐛 Troubleshooting

### Error: "workflow: command not found"
- Verifica que el binario `workflow` existe en el directorio raíz
- Asegúrate de que tiene permisos de ejecución: `chmod +x workflow`

### Error: "permission denied"
- Ejecuta con sudo si es necesario: `sudo make -f Makefile.test test-run`
- Verifica permisos del directorio: `ls -la`

### Error: "Docker daemon not running"
- Inicia Docker: `sudo systemctl start docker`
- Verifica que Docker esté funcionando: `docker --version`

### Tests fallan en CI/CD
- Asegúrate de que el binario esté compilado para la arquitectura correcta
- Verifica que todas las dependencias estén en el Dockerfile

## 📊 Interpretación de Resultados

### ✅ Éxito
```
✅ workflow --version
✅ workflow add task 1
✅ Todos los tests básicos pasaron
```

### ❌ Fallo
```
❌ workflow add task 1
Error: could not add task
```

### ⚠️ Advertencia
```
⚠️ config.json no existe (se creará automáticamente)
```

## 🔄 Integración con CI/CD

### GitHub Actions
```yaml
name: Test workflow CLI
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - name: Build binary
        run: make build
      - name: Run tests
        run: make -f Makefile.test test-run
```

### GitLab CI
```yaml
test:
  stage: test
  image: docker:latest
  services:
    - docker:dind
  script:
    - make build
    - make -f Makefile.test test-run
```

## 📝 Notas de Desarrollo

- Los tests se ejecutan en un entorno aislado (Ubuntu 22.04)
- Cada test es independiente y puede ejecutarse por separado
- Los datos de test se crean y destruyen en cada ejecución
- El script verifica tanto la funcionalidad como la estructura de archivos

## 🤝 Contribuir

Para agregar nuevos tests:

1. Edita `test.sh` y agrega tu test
2. Verifica que funcione localmente
3. Actualiza esta documentación si es necesario
4. Haz commit de los cambios

¡Los tests ayudan a mantener la calidad y confiabilidad de workflow CLI! 