# 🚀 Fase 2: Comandos Core - Plan Detallado

## 🎯 **Objetivo de la Fase 2**

Implementar los comandos core de Harvest CLI en Go, manteniendo compatibilidad con la versión Python y verificando cada punto antes de continuar.

## 📋 **Puntos de la Fase 2**

### **Punto 1: TaskManager (Gestión de Tareas)**
**Archivo**: `harvest-go/internal/core/task.go`

**Funcionalidades a implementar:**
- [ ] Cargar tareas desde JSON
- [ ] Guardar tareas en JSON
- [ ] Agregar nueva tarea
- [ ] Obtener tareas del día
- [ ] Calcular total de horas

**Verificación:**
```bash
# Test: Cargar configuración existente
cd harvest-go && go run cmd/harvest/main.go
# Debe mostrar ayuda sin errores

# Test: Verificar compatibilidad con datos Python
# Los datos de ~/.harvest/tasks.json deben ser legibles
```

### **Punto 2: Comando `add` (Agregar Tareas)**
**Archivo**: `harvest-go/internal/cli/commands.go`

**Funcionalidades a implementar:**
- [ ] Comando `add <description> <hours> [category]`
- [ ] Validación de horas (números y texto)
- [ ] Categoría por defecto "general"
- [ ] Mostrar estado después de agregar

**Verificación:**
```bash
# Test: Agregar tarea básica
./harvest add "Test task" 2.0
# Debe agregar tarea y mostrar estado

# Test: Agregar con categoría
./harvest add "Test task" 1.5 tech
# Debe agregar con categoría tech

# Test: Validación de horas
./harvest add "Test" invalid
# Debe mostrar error de validación
```

### **Punto 3: Comando `status` (Mostrar Estado)**
**Archivo**: `harvest-go/internal/cli/commands.go`

**Funcionalidades a implementar:**
- [ ] Mostrar fecha actual
- [ ] Mostrar horas trabajadas vs objetivo
- [ ] Mostrar horas restantes
- [ ] Listar tareas del día con iconos
- [ ] Barra de progreso visual

**Verificación:**
```bash
# Test: Mostrar estado sin tareas
./harvest status
# Debe mostrar 0h / 8h

# Test: Mostrar estado con tareas
./harvest add "Task 1" 2.0
./harvest add "Task 2" 1.5
./harvest status
# Debe mostrar 3.5h / 8h y listar tareas
```

### **Punto 4: Comandos Específicos (`tech`, `meeting`, `qa`)**
**Archivo**: `harvest-go/internal/cli/commands.go`

**Funcionalidades a implementar:**
- [ ] Comando `tech <description> <hours>`
- [ ] Comando `meeting <description> <hours>`
- [ ] Comando `qa <description> <hours>`
- [ ] Comando `daily` (0.25h automático)

**Verificación:**
```bash
# Test: Comandos específicos
./harvest tech "Development" 3.0
./harvest meeting "Team sync" 1.0
./harvest qa "Testing" 1.5
./harvest daily

# Verificar que se agregaron con categorías correctas
./harvest status
```

### **Punto 5: Comando `report` (Generar Reporte)**
**Archivo**: `harvest-go/internal/cli/commands.go`

**Funcionalidades a implementar:**
- [ ] Generar reporte para Harvest
- [ ] Formato: "Description - X.Xh"
- [ ] Mostrar total de horas
- [ ] Copiar al portapapeles (opcional)

**Verificación:**
```bash
# Test: Generar reporte
./harvest add "Task 1" 2.0
./harvest add "Task 2" 1.5
./harvest report
# Debe mostrar reporte formateado
```

## 🔄 **Estrategia de Desarrollo**

### **Metodología: Test-Driven Development**
1. **Implementar** funcionalidad básica
2. **Probar** con datos reales
3. **Verificar** compatibilidad con Python
4. **Commit** solo si todo funciona
5. **Continuar** al siguiente punto

### **Orden de Implementación**
1. **TaskManager** (base para todo)
2. **Comando `add`** (funcionalidad core)
3. **Comando `status`** (visualización)
4. **Comandos específicos** (conveniencia)
5. **Comando `report`** (salida)

## 🧪 **Tests de Compatibilidad**

### **Test 1: Datos Existentes**
```bash
# Verificar que puede leer datos de Python
./harvest status
# Debe mostrar tareas existentes de Python
```

### **Test 2: Escritura Compatible**
```bash
# Agregar tarea con Go
./harvest add "Go task" 2.0

# Verificar que Python puede leerla
cd .. && python3 harvest status
# Debe mostrar la tarea agregada por Go
```

### **Test 3: Formato JSON**
```bash
# Verificar que el formato JSON es compatible
cat ~/.harvest/tasks.json
# Debe tener el mismo formato que Python
```

## 📊 **Criterios de Éxito**

### **Funcionalidad**
- [ ] Todos los comandos funcionan
- [ ] Validaciones correctas
- [ ] Mensajes de error claros
- [ ] Compatibilidad con Python

### **Performance**
- [ ] Inicio rápido (< 100ms)
- [ ] Respuesta inmediata
- [ ] Sin errores de memoria

### **UX**
- [ ] Mensajes claros
- [ ] Iconos correctos
- [ ] Formato consistente
- [ ] Ayuda útil

## 🎯 **Próximo Punto: TaskManager**

**Archivo a crear**: `harvest-go/internal/core/task.go`

**Funciones a implementar**:
- `LoadTasks() ([]Task, error)`
- `SaveTasks(tasks []Task) error`
- `AddTask(description string, hours float64, category string) error`
- `GetTodayTasks() ([]Task, error)`
- `GetTotalHours(tasks []Task) float64`

**Verificación**:
```bash
cd harvest-go
# Crear archivo task.go
# Implementar funciones básicas
# Test: go run cmd/harvest/main.go
# Debe compilar sin errores
```

---

**Estado**: 🟡 **Listo para comenzar**
**Próximo paso**: Implementar TaskManager
**Tiempo estimado**: 30-60 minutos por punto 