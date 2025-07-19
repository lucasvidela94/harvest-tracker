# Makefile para Harvest CLI

# Variables
BINARY_NAME=harvest
BUILD_DIR=build
VERSION=1.1.0

# Colores para output
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

.PHONY: all build clean test help

# Target por defecto
all: build

# Construir el binario
build:
	@echo "$(GREEN)🔨 Building Harvest CLI...$(NC)"
	go build -o $(BINARY_NAME) ./cmd/harvest
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
	GOOS=linux GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 ./cmd/harvest
	GOOS=linux GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 ./cmd/harvest
	
	# macOS
	GOOS=darwin GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 ./cmd/harvest
	GOOS=darwin GOARCH=arm64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 ./cmd/harvest
	
	# Windows
	GOOS=windows GOARCH=amd64 go build -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe ./cmd/harvest
	
	@echo "$(GREEN)✅ Multi-platform build completed!$(NC)"
	@echo "$(YELLOW)📁 Binaries created in $(BUILD_DIR)/$(NC)"

# Ejecutar el binario
run: build
	@echo "$(GREEN)🚀 Running Harvest CLI...$(NC)"
	./$(BINARY_NAME)

# Verificar el código
check:
	@echo "$(GREEN)🔍 Checking code...$(NC)"
	go vet ./...
	gofmt -d .
	@echo "$(GREEN)✅ Code check completed!$(NC)"

# Ayuda
help:
	@echo "$(GREEN)🌾 Harvest CLI - Makefile$(NC)"
	@echo ""
	@echo "$(YELLOW)Available commands:$(NC)"
	@echo "  make build      - Build the binary"
	@echo "  make clean      - Clean build files"
	@echo "  make test       - Run tests"
	@echo "  make deps       - Install dependencies"
	@echo "  make build-all  - Build for all platforms"
	@echo "  make run        - Build and run"
	@echo "  make check      - Check code quality"
	@echo "  make help       - Show this help" 