# 🚀 Fase 3: Sistema de Upgrade - Plan Detallado

## 🎯 **Objetivo de la Fase 3**

Implementar un sistema de upgrade que permita a los usuarios migrar desde la versión Python a la versión Go de Harvest CLI, manteniendo todos sus datos y configuración.

## 📋 **Puntos de la Fase 3**

### **Punto 1: Detección de Versión Actual**
**Archivo**: `harvest-go/internal/upgrade/version.go`

**Funcionalidades a implementar:**
- [ ] Detectar si existe instalación Python
- [ ] Leer versión actual de Python
- [ ] Comparar con versión Go disponible
- [ ] Determinar si hay actualización disponible

**Verificación:**
```bash
# Test: Detectar versión Python
./harvest --upgrade
# Debe detectar versión Python 1.1.0

# Test: Sin instalación previa
# En sistema limpio, debe crear instalación nueva
```

### **Punto 2: Backup de Datos**
**Archivo**: `harvest-go/internal/upgrade/backup.go`

**Funcionalidades a implementar:**
- [ ] Crear backup de configuración Python
- [ ] Crear backup de datos de tareas
- [ ] Crear backup de archivos de instalación
- [ ] Verificar integridad del backup

**Verificación:**
```bash
# Test: Backup automático
./harvest --upgrade
# Debe crear ~/.harvest/backup/ con todos los datos

# Test: Verificar backup
ls ~/.harvest/backup/
# Debe contener config.json, tasks.json, etc.
```

### **Punto 3: Descarga de Nueva Versión**
**Archivo**: `harvest-go/internal/upgrade/download.go`

**Funcionalidades a implementar:**
- [ ] Conectar a GitHub API
- [ ] Obtener última versión disponible
- [ ] Descargar binario para la plataforma
- [ ] Verificar checksum del archivo

**Verificación:**
```bash
# Test: Descarga de versión
./harvest --upgrade
# Debe descargar binario desde GitHub

# Test: Verificar archivo
ls -la ~/.harvest/
# Debe tener nuevo binario
```

### **Punto 4: Instalación y Migración**
**Archivo**: `harvest-go/internal/upgrade/install.go`

**Funcionalidades a implementar:**
- [ ] Reemplazar binario Python con Go
- [ ] Restaurar configuración y datos
- [ ] Actualizar PATH si es necesario
- [ ] Verificar instalación exitosa

**Verificación:**
```bash
# Test: Instalación completa
./harvest --upgrade
# Debe reemplazar Python con Go

# Test: Verificar funcionalidad
harvest status
# Debe funcionar con datos restaurados
```

### **Punto 5: Rollback y Recuperación**
**Archivo**: `harvest-go/internal/upgrade/rollback.go`

**Funcionalidades a implementar:**
- [ ] Detectar fallos en instalación
- [ ] Restaurar versión anterior
- [ ] Recuperar datos del backup
- [ ] Notificar al usuario

**Verificación:**
```bash
# Test: Rollback automático
# Simular fallo y verificar recuperación

# Test: Recuperación manual
./harvest --rollback
# Debe restaurar versión anterior
```

## 🔄 **Estrategia de Desarrollo**

### **Metodología: Migración Segura**
1. **Detectar** instalación existente
2. **Backup** completo de datos
3. **Descargar** nueva versión
4. **Instalar** y migrar datos
5. **Verificar** funcionamiento
6. **Rollback** si hay problemas

### **Orden de Implementación**
1. **Detección de versión** (base para todo)
2. **Sistema de backup** (seguridad)
3. **Descarga de archivos** (obtener nueva versión)
4. **Instalación y migración** (proceso principal)
5. **Rollback** (seguridad adicional)

## 🧪 **Tests de Compatibilidad**

### **Test 1: Migración desde Python**
```bash
# Instalar versión Python
python3 harvest --upgrade  # Instalar Python

# Migrar a Go
./harvest --upgrade
# Debe migrar exitosamente
```

### **Test 2: Datos Preservados**
```bash
# Agregar tareas con Python
python3 harvest add "Python task" 2.0

# Migrar a Go
./harvest --upgrade

# Verificar datos
./harvest status
# Debe mostrar tarea de Python
```

### **Test 3: Configuración Migrada**
```bash
# Configurar Python
python3 harvest config set daily_hours_target 6.0

# Migrar a Go
./harvest --upgrade

# Verificar configuración
./harvest status
# Debe usar 6.0h como objetivo
```

## 📊 **Criterios de Éxito**

### **Funcionalidad**
- [ ] Migración automática desde Python
- [ ] Preservación completa de datos
- [ ] Configuración migrada correctamente
- [ ] Rollback automático en caso de fallo

### **Seguridad**
- [ ] Backup completo antes de migración
- [ ] Verificación de integridad
- [ ] Rollback automático
- [ ] Notificaciones claras al usuario

### **UX**
- [ ] Proceso automático y transparente
- [ ] Mensajes claros de progreso
- [ ] Confirmación antes de cambios
- [ ] Instrucciones de recuperación

## 🎯 **Próximo Punto: Detección de Versión**

**Archivo a crear**: `harvest-go/internal/upgrade/version.go`

**Funciones a implementar**:
- `DetectPythonInstallation() (bool, string, error)`
- `GetCurrentVersion() (string, error)`
- `GetLatestVersion() (string, error)`
- `CompareVersions(current, latest string) (bool, error)`

**Verificación**:
```bash
cd harvest-go
# Crear archivo version.go
# Implementar funciones básicas
# Test: go run cmd/harvest/main.go --upgrade
# Debe detectar versión Python
```

---

**Estado**: 🟡 **Listo para comenzar**
**Próximo paso**: Implementar detección de versión
**Tiempo estimado**: 30-60 minutos por punto 