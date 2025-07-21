#!/bin/bash

# workflow CLI - Instalador Automático
# Descarga e instala la última versión de workflow CLI

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuración
REPO="lucasvidela94/workflow-cli"
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

echo -e "${BLUE}🚀 workflow CLI - Instalador Automático${NC}"
echo -e "${YELLOW}Versión a instalar: $LATEST_VERSION${NC}"

# Detectar plataforma
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}❌ Arquitectura no soportada: $ARCH${NC}"
        exit 1
        ;;
esac

# Determinar extensión del archivo
if [ "$OS" = "darwin" ]; then
    EXT=""
elif [ "$OS" = "linux" ]; then
    EXT=""
elif [ "$OS" = "windows" ]; then
    EXT=".exe"
else
    echo -e "${RED}❌ Sistema operativo no soportado: $OS${NC}"
    exit 1
fi

# URL del release (binario individual)
RELEASE_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/workflow-$OS-$ARCH$EXT"

echo -e "${BLUE}📦 Descargando desde: $RELEASE_URL${NC}"

# Crear directorio temporal
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Descargar release
echo -e "${YELLOW}⬇️  Descargando...${NC}"
curl -L -o "workflow" "$RELEASE_URL"

# Verificar que el archivo se descargó correctamente
if [ ! -f "workflow" ] || [ ! -s "workflow" ]; then
    echo -e "${RED}❌ Error: No se pudo descargar el archivo${NC}"
    exit 1
fi

# Hacer ejecutable
chmod +x workflow

# Instalar
echo -e "${YELLOW}🔧 Instalando...${NC}"
sudo mv "workflow" /usr/local/bin/workflow

# Limpiar
cd /
rm -rf "$TEMP_DIR"

# Verificar instalación
if command -v workflow >/dev/null 2>&1; then
    echo -e "${GREEN}✅ ¡workflow CLI instalado exitosamente!${NC}"
    echo -e "${BLUE}📋 Versión instalada:${NC}"
    workflow version
    echo -e "${BLUE}💡 Para ver ayuda: workflow --help${NC}"
else
    echo -e "${RED}❌ Error en la instalación${NC}"
    exit 1
fi 