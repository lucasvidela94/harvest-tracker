# Dockerfile para probar la instalación enterprise de workflow CLI
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
RUN useradd -m -s /bin/bash workflow-user
USER workflow-user
WORKDIR /home/workflow-user

# Probar el método de instalación enterprise
RUN echo "🚀 Probando instalación enterprise de workflow CLI..." && \
    echo "📦 Descargando e instalando..." && \
    curl -fsSL https://raw.githubusercontent.com/lucasvidela94/workflow-cli/main/install-latest.sh | bash

# Verificar que la instalación funcionó
RUN echo "✅ Verificando instalación..." && \
    workflow version && \
    workflow --help

# Probar algunos comandos básicos
RUN echo "🧪 Probando comandos básicos..." && \
    workflow add "Tarea de prueba" 2.0 && \
    workflow status && \
    workflow list

# Mostrar información del sistema
RUN echo "📊 Información del sistema:" && \
    echo "OS: $(uname -s)" && \
    echo "Arch: $(uname -m)" && \
    echo "workflow CLI: $(workflow version)"

# Comando por defecto
CMD ["workflow", "--help"] 