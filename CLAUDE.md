# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build Commands

```bash
# Build server application
make build

# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make format

# Lint code (requires golangci-lint)
make lint

# Build Docker image
make docker-build
```

## Running a Single Test

```bash
go test -v -run TestFunctionName ./pkg/...
```

## Architecture

### Entry Point

- **`cmd/server/main.go`** - HTTP API server using go-arg

### Core Packages

- **`pkg/event/`** - Event metadata handling
  - `event.go` - Event struct and ParseEvent function

- **`pkg/clip/`** - Clip file discovery and organization
  - `discovery.go` - FilesByType, GetFileType, camera constants
  - `sorting.go` - SortByName for directory entries

- **`pkg/server/`** - HTTP API server (Gin framework)
  - `server.go` - Server setup, middleware, health endpoints
  - `clips.go` - Clip retrieval endpoints

- **`pkg/config/`** - Configuration management

- **`pkg/errors/`** - Standardized error types with codes

- **`pkg/sentry.go`** - Backward compatibility facade (deprecated)

### Configuration Priority

1. CLI flags
2. Environment variables (prefix: `TESLA_SENTRY_`)
3. Config file (`~/.tesla-sentry-viewer.yaml`)
4. Default values

### Key Dependencies

- **gin-gonic/gin** - HTTP framework
- **sirupsen/logrus** - Logging

### Tesla Sentry Clip Structure

Tesla Sentry clips are organized by timestamp directories matching `\d{4}-\d{2}-\d{2}_\d{2}-\d{2}-\d{2}`. Each directory contains:
- `event.json` - Event metadata
- Multiple `.mp4` files grouped by camera type (front, back, left_repeater, right_repeater)
