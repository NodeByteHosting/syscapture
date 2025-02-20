# Variables
APP_NAME=syscapture
BUILD_DIR=dist
CLIENT_DIR=client
SWAGGER_DIR=docs

# Targets
.PHONY: all build run clean swagger

all: build

# Install dependencies
install:
	@echo "Installing dependencies..."
	@go mod tidy
	@go get -u github.com/swaggo/gin-swagger
	@go get -u github.com/swaggo/files

# Build the application
build:
	@echo "Building $(APP_NAME), please wait..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) ./$(CLIENT_DIR)/main.go


# Run the application
run: build
	@echo "Starting $(APP_NAME), please wait..."
	@./$(BUILD_DIR)/$(APP_NAME)

run-client:
	@echo "Starting $(APP_NAME), please wait... "
	@go run $(CLIENT_DIR)/main.go

# Clean up build artifacts
clean:
	@echo "Cleaning up build artifacts"
	@rm -f $(APP_NAME)
	@rm -rf $(SWAGGER_DIR)

# Build Documentation for the V2 API Endpoints
docs:
	@echo "Generating OpenAPI documentation"
	@swag i -g main.go --dir api
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

version:
	@echo "Checking for the Version Information"
	@go run $(CLIENT_DIR)/main.go --version
