# Makefile para workflow CLI

# Variables
BINARY_NAME=workflow
BUILD_DIR=build
VERSION=1.1.0

# Colores para output
GREEN=\033[0;32m
YELLOW=\033[1;33m
BLUE=\033[0;34m
RED=\033[0;31m
NC=\033[0m # No Color

# Detectar sistema operativo
OS=$(shell uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(shell uname -m)

# Rutas de instalación
ifeq ($(OS),darwin)
    INSTALL_DIR=$(HOME)/.local/bin
else ifeq ($(OS),linux)
    INSTALL_DIR=$(HOME)/.local/bin
else
    INSTALL_DIR=$(HOME)/bin
endif

.PHONY: all build clean test help install install-script uninstall uninstall-script dist check dev

# Target por defecto
all: build

# Construir el binario
build:
	@echo "$(GREEN)🔨 Building workflow CLI...$(NC)"
	go build -o $(BINARY_NAME) ./cmd/workflow
	@echo "$(GREEN)✅ Build completed!$(NC)"

# Limpiar archivos generados
clean:
	@echo "$(YELLOW)🧹 Cleaning build files...$(NC)"
	rm -f $(BINARY_NAME)
	rm -rf $(BUILD_DIR)
	@echo "$(GREEN)✅ Clean completed!$(NC)"

# Ejecutar tests
test:
	@echo "$(GREEN)🧪 Running tests...$(NC)"
	go test ./...
	@echo "$(GREEN)✅ Tests completed!$(NC)"

# Instalar dependencias
deps:
	@echo "$(GREEN)📦 Installing dependencies...$(NC)"
	go mod tidy
	go mod download
	@echo "$(GREEN)✅ Dependencies installed!$(NC)"

# Build para múltiples plataformas
build-all: clean
	@echo "$(GREEN)🔨 Building for multiple platforms...$(NC)"
	mkdir -p $(BUILD_DIR)
	
	# Linux
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/workflow
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/workflow
	
	# macOS
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/workflow
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/workflow
	
	# Windows
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/workflow
	
	@echo "$(GREEN)✅ Multi-platform build completed!$(NC)"
	@echo "$(YELLOW)📁 Binaries created in $(BUILD_DIR)/$(NC)"

# Instalación manual
install: build
	@echo "$(GREEN)📦 Installing $(BINARY_NAME)...$(NC)"
	@mkdir -p $(INSTALL_DIR)
	@cp $(BINARY_NAME) $(INSTALL_DIR)/
	@chmod +x $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(GREEN)✅ Installed in $(INSTALL_DIR)/$(BINARY_NAME)$(NC)"
	@echo ""
	@echo "$(BLUE)💡 To use 'workflow' from anywhere:$(NC)"
	@echo "1. Add $(INSTALL_DIR) to your PATH"
	@echo "2. Or run: export PATH=\"$(INSTALL_DIR):\$$PATH\""
	@echo ""
	@echo "$(BLUE)📝 Usage examples:$(NC)"
	@echo "  workflow --help"
	@echo "  workflow add 'Example task' 2.0"

# Instalación usando script (recomendado)
install-script:
	@echo "$(GREEN)📦 Installing using script...$(NC)"
	@chmod +x install.sh
	@./install.sh

# Desinstalación manual
uninstall:
	@echo "$(YELLOW)🗑️ Uninstalling $(BINARY_NAME)...$(NC)"
	@rm -f $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "$(GREEN)✅ Uninstalled from $(INSTALL_DIR)$(NC)"

# Desinstalación usando script
uninstall-script:
	@echo "$(YELLOW)🗑️ Uninstalling using script...$(NC)"
	@chmod +x uninstall.sh
	@./uninstall.sh

# Crear distribución
dist: build-all
	@echo "$(GREEN)📦 Creating distribution...$(NC)"
	@mkdir -p $(BUILD_DIR)/dist
	@for binary in $(BUILD_DIR)/$(BINARY_NAME)-*; do \
		platform=$$(basename $$binary | sed 's/$(BINARY_NAME)-//'); \
		dist_name="$(BINARY_NAME)-$(VERSION)-$$platform"; \
		mkdir -p "$(BUILD_DIR)/dist/$$dist_name"; \
		cp "$$binary" "$(BUILD_DIR)/dist/$$dist_name/$(BINARY_NAME)"; \
		cp README.md "$(BUILD_DIR)/dist/$$dist_name/" 2>/dev/null || true; \
		cp install.sh "$(BUILD_DIR)/dist/$$dist_name/" 2>/dev/null || true; \
		cp uninstall.sh "$(BUILD_DIR)/dist/$$dist_name/" 2>/dev/null || true; \
		cd "$(BUILD_DIR)/dist" && tar -czf "$$dist_name.tar.gz" "$$dist_name"; \
		rm -rf "$$dist_name"; \
		echo "$(GREEN)✅ Created $$dist_name.tar.gz$(NC)"; \
	done

# Verificar instalación
check:
	@echo "$(BLUE)🔍 Checking installation...$(NC)"
	@if command -v $(BINARY_NAME) >/dev/null 2>&1; then \
		echo "$(GREEN)✅ $(BINARY_NAME) is installed and available in PATH$(NC)"; \
		$(BINARY_NAME) --version; \
	else \
		echo "$(RED)❌ $(BINARY_NAME) is not in PATH$(NC)"; \
		if [ -f "$(INSTALL_DIR)/$(BINARY_NAME)" ]; then \
			echo "$(GREEN)✅ $(BINARY_NAME) is installed in $(INSTALL_DIR) but not in PATH$(NC)"; \
			echo "$(YELLOW)Run: export PATH=\"$(INSTALL_DIR):\$$PATH\"$(NC)"; \
		else \
			echo "$(RED)❌ $(BINARY_NAME) is not installed$(NC)"; \
		fi; \
	fi

# Ejecutar el binario
run: build
	@echo "$(GREEN)🚀 Running workflow CLI...$(NC)"
	./$(BINARY_NAME)

# Modo desarrollo
dev: build
	@echo "$(GREEN)🚀 Running in development mode...$(NC)"
	@./$(BINARY_NAME) --help

# Verificar el código
code-check:
	@echo "$(GREEN)🔍 Checking code...$(NC)"
	go vet ./...
	gofmt -d .
	@echo "$(GREEN)✅ Code check completed!$(NC)"

# Ayuda
help:
	@echo "$(GREEN)🌾 workflow CLI - Makefile$(NC)"
	@echo ""
	@echo "$(YELLOW)🔨 Build commands:$(NC)"
	@echo "  make build         - Build the binary"
	@echo "  make build-all     - Build for all platforms"
	@echo "  make clean         - Clean build files"
	@echo "  make deps          - Install dependencies"
	@echo ""
	@echo "$(YELLOW)📦 Installation:$(NC)"
	@echo "  make install       - Install manually"
	@echo "  make install-script- Install using script (recommended)"
	@echo "  make uninstall     - Uninstall manually"
	@echo "  make uninstall-script- Uninstall using script"
	@echo ""
	@echo "$(YELLOW)🧪 Testing & Development:$(NC)"
	@echo "  make test          - Run tests"
	@echo "  make check         - Check installation"
	@echo "  make dev           - Run in development mode"
	@echo "  make code-check    - Check code quality"
	@echo ""
	@echo "$(YELLOW)📦 Distribution:$(NC)"
	@echo "  make dist          - Create distribution files"
	@echo ""
	@echo "$(YELLOW)💡 Recommendation:$(NC)"
	@echo "  Use 'make install-script' for automatic installation" 