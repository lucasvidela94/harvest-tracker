#!/bin/bash

# Script para probar la instalación enterprise de workflow CLI usando Docker
# Usa el binario local para evitar problemas con releases inexistentes

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🐳 Probando workflow CLI con binario local${NC}"
echo -e "${YELLOW}Este script construirá un contenedor Docker y probará la funcionalidad${NC}"

# Verificar que existe el binario
if [ ! -f "./workflow" ]; then
    echo -e "${RED}❌ Error: No se encontró el binario './workflow'${NC}"
    echo -e "${YELLOW}Ejecuta 'make build' primero${NC}"
    exit 1
fi

# Crear Dockerfile temporal que use el binario local
cat > Dockerfile.test << 'EOF'
# Dockerfile para probar workflow CLI
FROM ubuntu:22.04

# Evitar prompts interactivos durante la instalación
ENV DEBIAN_FRONTEND=noninteractive

# Instalar dependencias básicas
RUN apt-get update && apt-get install -y \
    curl \
    wget \
    tar \
    gzip \
    ca-certificates \
    && rm -rf /var/lib/apt/lists/*

# Crear usuario no-root para probar
RUN useradd -m -s /bin/bash workflow-user

# Copiar el binario local y cambiar permisos como root
COPY workflow /home/workflow-user/workflow
RUN chmod +x /home/workflow-user/workflow && \
    chown workflow-user:workflow-user /home/workflow-user/workflow

USER workflow-user
WORKDIR /home/workflow-user

# Verificar que la instalación funcionó
RUN echo "✅ Verificando instalación..." && \
    ~/workflow version && \
    ~/workflow --help

# Probar algunos comandos básicos
RUN echo "🧪 Probando comandos básicos..." && \
    ~/workflow add "Tarea de prueba" 2.0 && \
    ~/workflow status && \
    ~/workflow list

# Mostrar información del sistema
RUN echo "📊 Información del sistema:" && \
    echo "OS: $(uname -s)" && \
    echo "Arch: $(uname -m)" && \
    echo "workflow CLI: $(~/workflow version)"

# Comando por defecto
CMD ["/home/workflow-user/workflow", "--help"]
EOF

# Construir la imagen Docker
echo -e "\n${BLUE}🔨 Construyendo imagen Docker...${NC}"
docker build -f Dockerfile.test -t workflow-cli-test-fixed .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Imagen construida exitosamente${NC}"
else
    echo -e "${RED}❌ Error construyendo la imagen${NC}"
    exit 1
fi

# Ejecutar el contenedor para probar la instalación
echo -e "\n${BLUE}🚀 Ejecutando contenedor de prueba...${NC}"
echo -e "${YELLOW}Esto probará la funcionalidad básica${NC}"

docker run --rm workflow-cli-test-fixed

# Probar comandos interactivos
echo -e "\n${BLUE}🧪 Probando comandos interactivos...${NC}"
docker run --rm -it workflow-cli-test-fixed bash -c "
echo '📝 Agregando tareas de prueba...'
~/workflow add 'Desarrollo de feature' 4.0
~/workflow add 'Reunión de equipo' 1.5
~/workflow add 'Testing' 2.0

echo '📊 Mostrando estado...'
~/workflow status

echo '📋 Listando tareas...'
~/workflow list

echo '🔍 Buscando tareas...'
~/workflow search 'feature'

echo '📈 Generando reporte...'
~/workflow report
"

# Limpiar
echo -e "\n${BLUE}🧹 Limpiando archivos temporales...${NC}"
rm -f Dockerfile.test
docker rmi workflow-cli-test-fixed 2>/dev/null || true

echo -e "${GREEN}✅ Test completado exitosamente${NC}" 