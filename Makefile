# Variables
APP_NAME=syscapture
CLIENT_DIR=client
BUILD_DIR=build
SWAGGER_DIR=docs
PORT=42000

# Targets
.PHONY: all build run clean swagger

all: build

# Build the application
build:
	@echo "Building $(APP_NAME)..."
	@go build -o $(APP_NAME) ./$(CLIENT_DIR)/main.go

# Run the application
run: build
	@echo "Starting $(APP_NAME) on port $(PORT)..."
	@./$(APP_NAME)

run-client:
	@echo "Starting $(APP_NAME) on port $(PORT)..."
	@go run $(CLIENT_DIR)/main.go

# Clean up build artifacts
clean:
	@echo "Cleaning up..."
	@rm -f $(APP_NAME)
	@rm -rf $(SWAGGER_DIR)

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger docs..."
	@cd $(CLIENT_DIR) && swag init -g main.go
	@echo "Swagger docs generated in $(SWAGGER_DIR)"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod tidy
	@go get -u github.com/swaggo/gin-swagger
	@go get -u github.com/swaggo/files

# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "Linting code..."
	@golint ./...

# Check for errors
vet:
	@echo "Checking for errors..."
	@go vet ./...

# Run all checks (fmt, lint, vet)
check: fmt lint vet 