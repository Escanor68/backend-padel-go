# Makefile para Backend Padel Go

.PHONY: help build run test clean deps swag docker-build docker-run

# Variables
BINARY_NAME=padel-backend
DOCKER_IMAGE=padel-backend
DOCKER_TAG=latest

# Ayuda
help: ## Mostrar esta ayuda
	@echo "Comandos disponibles:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# Desarrollo
deps: ## Instalar dependencias
	@echo "Instalando dependencias..."
	go mod tidy
	go mod download

build: ## Compilar la aplicación
	@echo "Compilando aplicación..."
	go build -o $(BINARY_NAME) main.go

run: ## Ejecutar la aplicación en modo desarrollo
	@echo "Ejecutando aplicación..."
	go run main.go

test: ## Ejecutar tests
	@echo "Ejecutando tests..."
	go test -v ./...

test-coverage: ## Ejecutar tests con cobertura
	@echo "Ejecutando tests con cobertura..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

# Documentación
swag: ## Generar documentación Swagger
	@echo "Generando documentación Swagger..."
	swag init

# Limpieza
clean: ## Limpiar archivos generados
	@echo "Limpiando archivos..."
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	go clean

# Docker
docker-build: ## Construir imagen Docker
	@echo "Construyendo imagen Docker..."
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

docker-run: ## Ejecutar contenedor Docker
	@echo "Ejecutando contenedor Docker..."
	docker run -p 8080:8080 --env-file .env $(DOCKER_IMAGE):$(DOCKER_TAG)

# Base de datos
db-migrate: ## Ejecutar migraciones de base de datos
	@echo "Ejecutando migraciones..."
	go run main.go migrate

# Linting
lint: ## Ejecutar linter
	@echo "Ejecutando linter..."
	golangci-lint run

fmt: ## Formatear código
	@echo "Formateando código..."
	go fmt ./...

# Desarrollo completo
dev: deps swag run ## Configurar entorno de desarrollo completo

# Producción
prod: clean build ## Compilar para producción
	@echo "Aplicación compilada para producción: $(BINARY_NAME)"
