.PHONY: dev build run test clean

# Desarrollo con hot reload
dev:
	air

# Construir binario
build:
	go build -o bin/api ./cmd/api

# Ejecutar modo producción
run:
	go run ./cmd/api

# Ejecutar tests
test:
	go test ./...

# Limpiar archivos temporales
clean:
	rm -rf tmp/ bin/

# Instalar dependencias de desarrollo
install-dev:
	go install github.com/air-verse/air@latest
