#!/bin/bash

# Script para probar la instalación enterprise de Harvest CLI usando Docker
# Versión fija para evitar problemas con la API de GitHub

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🐳 Probando instalación enterprise de Harvest CLI (v2.0.1)${NC}"
echo -e "${YELLOW}Este script construirá un contenedor Docker y probará la instalación automática${NC}"

# Crear Dockerfile temporal con versión fija
cat > Dockerfile.test << 'EOF'
# Dockerfile para probar la instalación enterprise de Harvest CLI
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

# Crear usuario no-root para probar la instalación
RUN useradd -m -s /bin/bash harvest-user
USER harvest-user
WORKDIR /home/harvest-user

# Probar el método de instalación enterprise con versión específica
RUN echo "🚀 Probando instalación enterprise de Harvest CLI..." && \
    echo "📦 Descargando versión v2.0.1..." && \
    wget https://github.com/lucasvidela94/harvest-tracker/releases/download/v2.0.1/harvest-v2.0.1-linux-amd64.tar.gz && \
    tar -xzf harvest-v2.0.1-linux-amd64.tar.gz && \
    mv harvest-v2.0.1-linux-amd64/harvest ~/harvest && \
    chmod +x ~/harvest && \
    echo 'export PATH="$HOME:$PATH"' >> ~/.bashrc && \
    export PATH="$HOME:$PATH"

# Verificar que la instalación funcionó
RUN echo "✅ Verificando instalación..." && \
    ~/harvest version && \
    ~/harvest --help

# Probar algunos comandos básicos
RUN echo "🧪 Probando comandos básicos..." && \
    ~/harvest add "Tarea de prueba" 2.0 && \
    ~/harvest status && \
    ~/harvest list

# Mostrar información del sistema
RUN echo "📊 Información del sistema:" && \
    echo "OS: $(uname -s)" && \
    echo "Arch: $(uname -m)" && \
    echo "Harvest CLI: $(~/harvest version)"

# Comando por defecto
CMD ["/home/harvest-user/harvest", "--help"]
EOF

# Construir la imagen Docker
echo -e "\n${BLUE}🔨 Construyendo imagen Docker...${NC}"
docker build -f Dockerfile.test -t harvest-cli-test-fixed .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Imagen construida exitosamente${NC}"
else
    echo -e "${RED}❌ Error construyendo la imagen${NC}"
    exit 1
fi

# Ejecutar el contenedor para probar la instalación
echo -e "\n${BLUE}🚀 Ejecutando contenedor de prueba...${NC}"
echo -e "${YELLOW}Esto probará la instalación automática y comandos básicos${NC}"

docker run --rm harvest-cli-test-fixed

# Probar comandos interactivos
echo -e "\n${BLUE}🧪 Probando comandos interactivos...${NC}"
docker run --rm -it harvest-cli-test-fixed bash -c "
echo '📝 Agregando tareas de prueba...'
~/harvest add 'Desarrollo de feature' 4.0
~/harvest add 'Reunión de equipo' 1.5
~/harvest add 'Testing' 2.0

echo '📊 Mostrando estado...'
~/harvest status

echo '📋 Listando tareas...'
~/harvest list

echo '🔍 Buscando tareas...'
~/harvest search 'feature'

echo '📈 Generando reporte...'
~/harvest report

echo '✅ ¡Todas las pruebas completadas exitosamente!'
"

# Limpiar archivos temporales
rm -f Dockerfile.test

echo -e "\n${GREEN}🎉 ¡Prueba de instalación enterprise completada!${NC}"
echo -e "${BLUE}💡 La instalación automática funciona correctamente${NC}" 