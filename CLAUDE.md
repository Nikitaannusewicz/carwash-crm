# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Carwash CRM is a Go-based API for managing carwash business operations. The project follows standard Go project layout conventions with a modular domain-driven architecture.

## Development Commands

### Running the Application
```bash
go run cmd/api/main.go
```

### Building
```bash
go build -o bin/api cmd/api/main.go
```

### Testing
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests for a specific package
go test ./internal/server

# Run a specific test
go test -run TestName ./internal/package
```

### Dependencies
```bash
# Add a dependency
go get github.com/some/package

# Tidy dependencies
go mod tidy

# Vendor dependencies
go mod vendor
```

## Architecture

### Project Structure
- `cmd/api/main.go` - Application entrypoint; initializes config, logger, and server
- `internal/config/` - Environment-based configuration loader (PORT, DATABASE_URL, JWT_SECRET)
- `internal/server/` - HTTP server setup using standard library `http.ServeMux` with route registration
- `internal/database/` - Database connection and migration logic (currently empty)
- `internal/modules/` - Domain modules implementing business logic:
  - `identity/` - User authentication and authorization
  - `operations/` - Core carwash operations management
  - `scheduling/` - Appointment and scheduling features

### Key Architectural Patterns

**Server Initialization Flow:**
1. `main.go` loads config via `config.LoadConfig()`
2. Creates new server with `server.NewServer(cfg)`
3. Server constructor calls `registerRoutes()` to set up HTTP handlers
4. Server runs with configured timeouts (idle: 1min, read: 10s, write: 30s)

**Configuration:**
Environment variables are loaded through `internal/config/config.go` with fallback defaults. Key variables:
- `PORT` (default: 8080)
- `DATABASE_URL` (default: empty)
- `JWT_SECRET` (default: "change-me-in-prod")

**HTTP Routing:**
Uses Go 1.22+ method-specific routing pattern: `s.router.HandleFunc("GET /path", handler)`. Routes are registered in `server.registerRoutes()`.

### Module Organization

Each module in `internal/modules/` should follow a consistent structure:
- `handler.go` - HTTP handlers
- `service.go` - Business logic
- `repository.go` - Data access layer
- `model.go` - Domain models

When adding new routes, register them in `server.registerRoutes()` (internal/server/server.go:29).
