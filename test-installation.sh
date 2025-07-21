#!/bin/bash

# Script para probar la instalación enterprise de workflow CLI usando Docker

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🐳 Probando instalación enterprise de workflow CLI${NC}"
echo -e "${YELLOW}Este script construirá un contenedor Docker y probará la instalación automática${NC}"

# Construir la imagen Docker
echo -e "\n${BLUE}🔨 Construyendo imagen Docker...${NC}"
docker build -t workflow-cli-test .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✅ Imagen construida exitosamente${NC}"
else
    echo -e "${RED}❌ Error construyendo la imagen${NC}"
    exit 1
fi

# Ejecutar el contenedor para probar la instalación
echo -e "\n${BLUE}🚀 Ejecutando contenedor de prueba...${NC}"
echo -e "${YELLOW}Esto probará la instalación automática y comandos básicos${NC}"

docker run --rm workflow-cli-test

# Probar comandos interactivos
echo -e "\n${BLUE}🧪 Probando comandos interactivos...${NC}"
docker run --rm -it workflow-cli-test bash -c "
echo '📝 Agregando tareas de prueba...'
workflow add 'Desarrollo de feature' 4.0
workflow add 'Reunión de equipo' 1.5
workflow add 'Testing' 2.0

echo '📊 Mostrando estado...'
workflow status

echo '📋 Listando tareas...'
workflow list

echo '🔍 Buscando tareas...'
workflow search 'feature'

echo '📈 Generando reporte...'
workflow report

echo '✅ ¡Todas las pruebas completadas exitosamente!'
"

echo -e "\n${GREEN}🎉 ¡Prueba de instalación enterprise completada!${NC}"
echo -e "${BLUE}💡 La instalación automática funciona correctamente${NC}" 