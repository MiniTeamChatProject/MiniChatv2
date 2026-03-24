# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a **Go room service** built with the [Kratos framework](https://github.com/go-kratos/kratos) — a Go microservices framework. The project follows clean architecture with dependency injection via Wire, using **GORM v2** for database access and **PostgreSQL** for persistence.

## Common Commands

```bash
# Initialize development environment (install tools)
make init

# Generate code from protobuf definitions
make api          # API proto files (generates Go code, HTTP/gRPC bindings, OpenAPI spec)
make config       # Internal proto files (config structs)
make all          # Both

# Generate Wire dependency injection
make generate     # or: cd cmd/room && wire

# Build the service
make build

# Run tests
go test ./...

# Run the service
./bin/room -conf ./configs

# Start/stop PostgreSQL (Docker)
docker-compose up -d
docker-compose down
docker-compose down -v    # Reset database (deletes all data)
```

## Architecture

The codebase follows Kratos's clean architecture pattern with four layers:

### 1. **API Layer** (`api/`)
- Protocol Buffer definitions for service contracts
- `make api` generates: Go code, HTTP/gRPC bindings, OpenAPI spec (`openapi.yaml`)
- Current service: `api/room/v1/room.proto`

### 2. **Service Layer** (`internal/service/`)
- Transport-agnostic business logic handlers
- Implements generated gRPC/HTTP interfaces
- Converts proto requests to biz domain models and back
- Wire providers: `service.ProviderSet`
- **TODO**: user_id is currently hardcoded (needs JWT/context integration)

### 3. **Business Logic Layer** (`internal/biz/`)
- Core business logic and use cases (`RoomUsecase`, `RoomMemberUsecase`)
- Domain models (`Room`, `RoomMember`) with enum types
- Repository interfaces (data access contracts)
- Permission logic (role-based access control)
- Wire providers: `biz.ProviderSet`

### 4. **Data Layer** (`internal/data/`)
- **GORM v2** for PostgreSQL access
- Data models in `internal/data/model/` with GORM tags and enums
- Repository implementations (`RoomRepository`, `RoomMemberRepository`)
- Auto-migration enabled in development (`enable_auto_migrate: true`)
- Wire providers: `data.ProviderSet`

### Dependency Injection
- **Wire** for compile-time dependency injection
- Entry point: `cmd/room/wire.go` → `wire_gen.go` (generated)
- Each layer exports a `ProviderSet` with its constructors

### Configuration
- Config file: `configs/config.yaml`
- Proto-based config schema: `internal/conf/conf.proto`
- Loaded at startup via `config.New()` with file source

### Servers
- HTTP server: `0.0.0.0:8000`
- gRPC server: `0.0.0.0:9000`
- Registered in `internal/server/`

### Database
- **PostgreSQL 16** in Docker (`docker-compose.yml`)
- Connection: `postgres://minichat:minichat123@localhost:5432/room_service`
- **GORM v1.25.12** (compatible with Go 1.24)
- Auto-migration runs on startup when `enable_auto_migrate: true`

## Data Model Structure

### Core Models (`internal/data/model/`)

**Room:**
- Types: `RoomType` (Group=1, Voice=2, Video=3, Live=4)
- Status: `RoomStatus` (Normal=1, Muted=2, Banned=3)
- Fields: ID, Name, OwnerID, Type, MaxCount, Avatar, Description, Tags, IsPublic, Status, CreatedAt, UpdatedAt

**RoomMember:**
- Roles: `MemberRole` (Member=1, Admin=2, Owner=3)
- Status: `MemberStatus` (Normal=1, Muted=2, Left=3)
- Fields: ID, RoomID, UserID, Role, Status, MuteUntil, JoinedAt, UpdatedAt

### Permission Model

| Operation | Owner | Admin | Member |
|-----------|-------|-------|--------|
| Update room | ✓ | ✓ | ✗ |
| Delete room | ✓ | ✗ | ✗ |
| Kick member | ✓ | ✓ (not admin/owner) | ✗ |
| Update role | ✓ | ✗ | ✗ |
| Mute member | ✓ | ✓ | ✗ |

## Adding a New Feature

1. Define API in `api/room/v1/*.proto` (include `google/api/annotations.proto` for HTTP)
2. Run `make api` to generate Go code and update `openapi.yaml`
3. Add domain model to `internal/biz/` (if new entity)
4. Add GORM model to `internal/data/model/` (if new table)
5. Create repo interface in `internal/biz/` and implementation in `internal/data/`
6. Create usecase in `internal/biz/` (if business logic needed)
7. Create service handler in `internal/service/`
8. Add to respective `ProviderSet`
9. Run `make generate` to regenerate wire
10. Update `internal/server/grpc.go` and `http.go` to register service

## Important Notes

### Database Schema Changes
- When modifying GORM models, run `docker-compose down -v` to reset the database
- Auto-migration will recreate tables on next startup

### GORM Version
- Using **GORM v1.25.12** for Go 1.24 compatibility
- Latest GORM requires Go 1.25+

### OpenAPI/Swagger
- `openapi.yaml` is generated in project root by `make api`
- Use https://editor.swagger.io/ to view API documentation
- Import `openapi.yaml` into Postman/Insomnia for testing

### User Authentication
- Currently **hardcoded** as `user_id = 1` in service layer
- TODO: Integrate JWT/context to extract real user_id
