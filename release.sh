#!/bin/bash

# Script para crear releases de Harvest
# Uso: ./release.sh [major|minor|patch]

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para mostrar mensajes
log() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verificar que estamos en el directorio correcto
if [ ! -f "VERSION" ] || [ ! -f "harvest" ]; then
    error "Debes ejecutar este script desde el directorio del proyecto harvest"
    exit 1
fi

# Leer versión actual
CURRENT_VERSION=$(cat VERSION)
log "Versión actual: $CURRENT_VERSION"

# Determinar tipo de release
RELEASE_TYPE=${1:-patch}

if [[ ! "$RELEASE_TYPE" =~ ^(major|minor|patch)$ ]]; then
    error "Tipo de release debe ser: major, minor, o patch"
    echo "Uso: $0 [major|minor|patch]"
    exit 1
fi

# Calcular nueva versión
IFS='.' read -ra VERSION_PARTS <<< "$CURRENT_VERSION"
MAJOR=${VERSION_PARTS[0]}
MINOR=${VERSION_PARTS[1]}
PATCH=${VERSION_PARTS[2]}

case $RELEASE_TYPE in
    major)
        NEW_MAJOR=$((MAJOR + 1))
        NEW_MINOR=0
        NEW_PATCH=0
        ;;
    minor)
        NEW_MAJOR=$MAJOR
        NEW_MINOR=$((MINOR + 1))
        NEW_PATCH=0
        ;;
    patch)
        NEW_MAJOR=$MAJOR
        NEW_MINOR=$MINOR
        NEW_PATCH=$((PATCH + 1))
        ;;
esac

NEW_VERSION="$NEW_MAJOR.$NEW_MINOR.$NEW_PATCH"
log "Nueva versión: $NEW_VERSION"

# Verificar que no hay cambios sin commitear
if [ -n "$(git status --porcelain)" ]; then
    error "Hay cambios sin commitear. Por favor, haz commit de todos los cambios antes de crear un release."
    git status --short
    exit 1
fi

# Actualizar archivo VERSION
echo "$NEW_VERSION" > VERSION
log "Archivo VERSION actualizado"

# Actualizar CHANGELOG
TODAY=$(date +%Y-%m-%d)
CHANGELOG_ENTRY="## [$NEW_VERSION] - $TODAY

### Added
- Nuevas características en esta versión

### Changed
- Cambios en funcionalidades existentes

### Fixed
- Correcciones de bugs

"

# Insertar nueva entrada al inicio del CHANGELOG (después del header)
sed -i "/^# Changelog$/a\\$CHANGELOG_ENTRY" CHANGELOG.md
log "CHANGELOG actualizado"

# Commit de los cambios
git add VERSION CHANGELOG.md
git commit -m "Bump version to $NEW_VERSION"
log "Commit de versión creado"

# Crear tag
git tag -a "v$NEW_VERSION" -m "Release v$NEW_VERSION"
log "Tag v$NEW_VERSION creado"

# Mostrar resumen
echo
log "Release $NEW_VERSION creado exitosamente!"
echo
echo "Para publicar el release:"
echo "  git push origin main"
echo "  git push origin v$NEW_VERSION"
echo
echo "Para crear un release en GitHub:"
echo "  1. Ve a https://github.com/[usuario]/[repo]/releases"
echo "  2. Haz clic en 'Create a new release'"
echo "  3. Selecciona el tag v$NEW_VERSION"
echo "  4. Agrega descripción del release"
echo "  5. Publica" 