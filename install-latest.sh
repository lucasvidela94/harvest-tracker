#!/bin/bash

# Harvest CLI - Instalador Automático
# Descarga e instala la última versión de Harvest CLI

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuración
REPO="lucasvidela94/harvest-tracker"
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

echo -e "${BLUE}🚀 Harvest CLI - Instalador Automático${NC}"
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
    EXT="tar.gz"
elif [ "$OS" = "linux" ]; then
    EXT="tar.gz"
else
    echo -e "${RED}❌ Sistema operativo no soportado: $OS${NC}"
    exit 1
fi

# URL del release
RELEASE_URL="https://github.com/$REPO/releases/download/$LATEST_VERSION/harvest-$LATEST_VERSION-$OS-$ARCH.$EXT"

echo -e "${BLUE}📦 Descargando desde: $RELEASE_URL${NC}"

# Crear directorio temporal
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Descargar release
echo -e "${YELLOW}⬇️  Descargando...${NC}"
curl -L -o "harvest.$EXT" "$RELEASE_URL"

# Verificar checksum (opcional)
echo -e "${YELLOW}🔐 Verificando integridad...${NC}"
curl -L -o "checksums.txt" "https://github.com/$REPO/releases/download/$LATEST_VERSION/harvest-$LATEST_VERSION-checksums.txt"
if command -v sha256sum >/dev/null 2>&1; then
    sha256sum -c checksums.txt --ignore-missing || echo -e "${YELLOW}⚠️  Advertencia: No se pudo verificar checksum${NC}"
fi

# Extraer
echo -e "${YELLOW}📁 Extrayendo...${NC}"
tar -xzf "harvest.$EXT"

# Instalar
echo -e "${YELLOW}🔧 Instalando...${NC}"
sudo mv "harvest-$LATEST_VERSION-$OS-$ARCH/harvest" /usr/local/bin/harvest
sudo chmod +x /usr/local/bin/harvest

# Limpiar
cd /
rm -rf "$TEMP_DIR"

# Verificar instalación
if command -v harvest >/dev/null 2>&1; then
    echo -e "${GREEN}✅ ¡Harvest CLI instalado exitosamente!${NC}"
    echo -e "${BLUE}📋 Versión instalada:${NC}"
    harvest version
    echo -e "${BLUE}💡 Para ver ayuda: harvest --help${NC}"
else
    echo -e "${RED}❌ Error en la instalación${NC}"
    exit 1
fi 