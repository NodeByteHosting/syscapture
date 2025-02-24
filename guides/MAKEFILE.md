# SysCapture Make Commands

## Overview
This document describes the available make commands for building, running, and maintaining the SysCapture application.

## Basic Commands

### Build and Run
```bash
# Build the application
make build

# Start the application (builds first)
make start

# Run in development mode
make dev
```

### Development Tools
```bash
# Install project dependencies
make install

# Check application version
make version

# Clean build artifacts
make clean
```

## Documentation

### API Documentation
```bash
# Generate Swagger documentation
make docs
```

## Code Quality

### Testing
```bash
# Run all tests
make test
```

### Code Maintenance
```bash
# Format code
make fmt

# Run linter
make lint

# Check for potential errors
make vet

# Run all code quality checks
make check
```

## Command Details

### Build Commands
- `build`: Compiles the application to `dist/syscapture`
- `start`: Builds and runs the compiled application
- `dev`: Runs the application directly using `go run`

### Dependency Management
- `install`: Installs project dependencies
- `deps`: Alternative command for dependency installation
  - Updates Swagger dependencies
  - Runs `go mod tidy`

### Documentation
- `docs`: Generates OpenAPI/Swagger documentation
  - Output directory: `docs/`
  - Processes API endpoints

### Testing and Quality
- `test`: Runs all project tests
- `fmt`: Formats Go code using `gofmt`
- `lint`: Runs `golint` for code style checks
- `vet`: Runs `go vet` to check for common errors
- `check`: Runs all code quality checks (fmt, lint, vet)

### Utility Commands
- `version`: Displays the current version
- `clean`: Removes build artifacts and documentation

## Directory Structure
```
syscapture/
├── dist/           # Build output
├── client/         # Client application code
├── docs/           # Generated Swagger docs
└── api/            # API definitions
```

## Examples

### Development Workflow
```bash
# Install dependencies
make install

# Run in development mode
make dev
```

### Build and Deploy
```bash
# Run all checks
make check

# Run tests
make test

# Build for deployment
make build
```

### Documentation Update
```bash
# Generate fresh documentation
make docs
```

## Notes
- Commands use `@` prefix for cleaner output
- Most commands include status messages
- Build artifacts are placed in `dist/` directory
- Documentation is auto-generated from API code