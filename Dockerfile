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

# Probar el método de instalación enterprise
RUN echo "🚀 Probando instalación enterprise de Harvest CLI..." && \
    echo "📦 Descargando e instalando..." && \
    curl -fsSL https://raw.githubusercontent.com/lucasvidela94/harvest-tracker/main/install-latest.sh | bash

# Verificar que la instalación funcionó
RUN echo "✅ Verificando instalación..." && \
    harvest version && \
    harvest --help

# Probar algunos comandos básicos
RUN echo "🧪 Probando comandos básicos..." && \
    harvest add "Tarea de prueba" 2.0 && \
    harvest status && \
    harvest list

# Mostrar información del sistema
RUN echo "📊 Información del sistema:" && \
    echo "OS: $(uname -s)" && \
    echo "Arch: $(uname -m)" && \
    echo "Harvest CLI: $(harvest version)"

# Comando por defecto
CMD ["harvest", "--help"] 