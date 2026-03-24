# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **Go room service** built with the [Kratos framework](https://github.com/go-kratos/kratos) — a Go microservices framework. The project follows clean architecture with dependency injection via Wire.

## Common Commands

```bash
# Initialize development environment (install tools)
make init

# Generate code from protobuf definitions (API + internal)
make api          # API proto files
make config       # Internal proto files
make all          # Both

# Generate Wire dependency injection
make generate     # or: cd cmd/room && wire

# Build the service
make build

# Run the service
./bin/room -conf ./configs
```

## Architecture

The codebase follows Kratos's clean architecture pattern with four layers:

### 1. **API Layer** (`api/`)
- Protocol Buffer definitions for service contracts
- `make api` generates: Go code, HTTP/gRPC bindings, OpenAPI spec
- Example: `api/helloworld/v1/greeter.proto`

### 2. **Service Layer** (`internal/service/`)
- Transport-agnostic business logic handlers
- Implements generated gRPC/HTTP interfaces
- Receives proto requests, calls biz layer, returns proto responses
- Wire providers: `service.ProviderSet`

### 3. **Business Logic Layer** (`internal/biz/`)
- Core business logic and use cases
- Defines repository interfaces (data access contracts)
- Wire providers: `biz.ProviderSet`

### 4. **Data Layer** (`internal/data/`)
- Repository implementations (database, external APIs)
- Data access and persistence
- Wire providers: `data.ProviderSet`

### Dependency Injection
- **Wire** is used for compile-time dependency injection
- Entry point: `cmd/room/wire.go` (injective) → `wire_gen.go` (generated)
- Each layer exports a `ProviderSet` with its constructors
- `cmd/room/main.go` loads config and starts the application

### Configuration
- Config file: `configs/config.yaml`
- Proto-based config in `internal/conf/conf.proto`
- Loaded at startup via `config.New()` with file source

### Servers
- HTTP server: `0.0.0.0:8000`
- gRPC server: `0.0.0.0:9000`
- Both registered in `internal/server/`

## Adding a New Feature

1. Define API in `api/{service}/v1/*.proto`
2. Run `make api` to generate bindings
3. Create biz layer: usecase + repo interface in `internal/biz/`
4. Implement repo in `internal/data/`
5. Create service handler in `internal/service/`
6. Wire everything: add to respective `ProviderSet`
7. Run `make generate` to regenerate wire
