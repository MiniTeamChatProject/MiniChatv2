# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

**MiniChat** is a real-time team collaboration chat platform built with a microservices architecture. The project consists of:

- **User Service** (`service/user`) - User registration, login, JWT authentication
- **Room Service** (`service/room`) - Room management, real-time messaging, WebSocket
- **Web Frontend** (`web`) - Next.js + React UI

**Team**: Wuhan University of Technology students - NameIsNotName (Guo Zifan), Wang Haokun, Chen Xinqi, Jonas Peng

## Quick Start

```bash
# Start all services (databases + microservices)
docker-compose up -d

# User Service: http://localhost:8000
# Room Service: http://localhost:8001
# User DB: localhost:15432
# Room DB: localhost:5432

# View Swagger UI
# User Service: http://localhost:8000/q/
# Room Service: http://localhost:8001/q/
```

## Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                    MiniChat 微服务架构                        │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  ┌──────────────┐        ┌──────────────┐                   │
│  │  用户服务     │        │  房间服务     │                   │
│  │  :8000       │   ◄──► │  :8001       │                   │
│  │              │  gRPC  │              │                   │
│  │  注册/登录    │        │  房间/聊天     │                   │
│  │  JWT Token   │        │  WebSocket   │                   │
│  └──────┬───────┘        └──────┬───────┘                   │
│         │                        │                           │
│         ▼                        ▼                           │
│  ┌──────────────┐        ┌──────────────┐                   │
│  │  用户数据库   │        │  房间数据库   │                   │
│  │  benutzer_db │        │ room_service │                   │
│  │  :15432      │        │   :5432      │                   │
│  └──────────────┘        └──────────────┘                   │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

## Service Port Allocation

| Service | HTTP | gRPC | Database | DB Port |
|---------|------|------|----------|---------|
| User Service | 8000 | 9000 | benutzer_db | 15432 |
| Room Service | 8001 | 9001 | room_service | 5432 |

## Project Structure

```
MiniChatv2/
├── service/
│   ├── user/           # User registration service
│   │   ├── api/        # Protobuf definitions
│   │   ├── cmd/        # Application entry point
│   │   ├── internal/   # Business logic (biz/data/service)
│   │   └── docker-compose.yaml
│   └── room/           # Room & chat service
│       ├── api/        # Protobuf definitions
│       ├── cmd/        # Application entry point
│       ├── internal/   # Business logic
│       └── docs/       # Integration documentation
├── web/                # Next.js frontend (WIP)
├── docker-compose.yml  # Multi-service orchestration
├── PRD.md             # Product requirements document
└── CLAUDE.md          # This file
```

## User Service (`service/user`)

### Overview
Handles user registration, login, JWT token generation, and user profile management.

### Tech Stack
- Go 1.22+, Kratos v2 framework
- PostgreSQL 16 with GORM
- JWT authentication + bcrypt password hashing
- HTTP/gRPC endpoints

### Key APIs
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/register` | POST | User registration |
| `/login` | POST | User login (returns JWT) |
| `/user/profile` | GET | Get user profile |
| `/user/profile` | PUT | Update user profile |
| `/user/{id}` | DELETE | Delete user |

### Quick Commands
```bash
cd service/user

# Install dependencies
go mod tidy

# Generate proto code
kratos proto client ./api/helloworld/v1/registration.proto

# Run service
kratos run
# OR
cd cmd/registration && go run .

# Build
go build -o ../../bin/user ./cmd/registration

# One-click setup (Linux)
chmod +x presetup.sh && ./presetup.sh
```

### Database
- Database: `benutzer_db`
- Table: `users` (id, username, password, nickname, email)
- Auto-migration on startup

## Room Service (`service/room`)

### Overview
Manages chat rooms, member permissions, real-time messaging via WebSocket, and message history.

### Tech Stack
- Go 1.24+, Kratos v2 framework
- PostgreSQL 16 with GORM v1.25.12
- gorilla/websocket for real-time communication
- Clean architecture (biz/data/service layers)

### Key APIs
| Endpoint | Method | Description |
|----------|--------|-------------|
| `/v1/rooms` | POST | Create room |
| `/v1/rooms/{id}` | GET | Get room info |
| `/v1/rooms/{id}/join` | POST | Join room |
| `/v1/rooms/{id}/leave` | POST | Leave room |
| `/v1/rooms/{id}/messages` | POST | Send message |
| `/v1/rooms/{id}/messages` | GET | Get message history |
| `/ws/room/{room_id}` | WebSocket | Real-time chat |

### Quick Commands
```bash
cd service/room

# Generate proto code
make api

# Generate Wire dependency injection
make generate

# Build
make build

# Run
./bin/room -conf ./configs

# Start database
docker-compose up -d
docker-compose down -v    # Reset database
```

### Data Models
- **Room**: id, name, owner_id, type, max_count, status, created_at
- **RoomMember**: id, room_id, user_id, role, status, mute_until
- **Message**: id, room_id, user_id, content, type, created_at

### WebSocket
Connect to: `ws://localhost:8001/ws/room/{room_id}`

Message format:
```json
{
  "type": "message",
  "data": {
    "content": "Hello!"
  }
}
```

## User-Service Integration

The room service integrates with user service for authentication:

### Option 1: API Gateway (Recommended)
```
Client → API Gateway (JWT validation) → Room Service
         ↓ adds user headers
    X-User-ID: 123
    X-Username: alice
```

### Option 2: Direct gRPC Call
Room service calls user service via gRPC to validate tokens:

```go
// In room service, call user service
client := pb.NewUserServiceClient(conn)
user, err := client.VerifyToken(ctx, &pb.VerifyTokenRequest{Token: token})
```

### Implementation Status
- ✅ Middleware created (`internal/middleware/`)
- ✅ User context extraction from headers
- ⚠️ Currently uses hardcoded `user_id = 1` for development
- TODO: Integrate with user service gRPC client

## Development Workflow

### Adding New Features to Room Service

1. **Define API** in `api/room/v1/room.proto`
2. **Generate code**: `make api`
3. **Add domain model** to `internal/biz/`
4. **Add GORM model** to `internal/data/model/`
5. **Implement repository** in `internal/data/`
6. **Implement usecase** in `internal/biz/`
7. **Implement service** in `internal/service/`
8. **Update Wire**: `make generate`
9. **Register in server**: update `internal/server/`

### Database Schema Changes

```bash
# Reset database (WARNING: deletes all data)
cd service/room
docker-compose down -v
docker-compose up -d

# Auto-migration will recreate tables on service startup
```

### Testing

```bash
# Test user service
curl -X POST http://localhost:8000/register \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "password": "pass123", "nickname": "Test User"}'

curl -X POST http://localhost:8000/login \
  -H "Content-Type: application/json" \
  -d '{"username": "test", "password": "pass123"}'

# Test room service (with simulated user header)
curl -X POST http://localhost:8001/v1/rooms/1/messages \
  -H "Content-Type: application/json" \
  -H "X-User-ID: 123" \
  -d '{"content": "Hello!"}'

# Test WebSocket
wscat -c "ws://localhost:8001/ws/room/1"
```

## Important Notes

### GORM Version Compatibility
- Room service uses **GORM v1.25.12** for Go 1.24 compatibility
- Latest GORM requires Go 1.25+

### WebSocket Authentication
Currently accepts user_id via:
- URL parameter: `ws://localhost:8001/ws/room/1?user_id=123`
- Header: `X-User-ID: 123` (during WebSocket handshake)

### Message Storage
Messages are stored in the **room service database** (`room_service.messages` table), not in the user database. This is the correct architectural pattern.

### Swagger UI
- User Service: http://localhost:8000/q/
- Room Service: http://localhost:8001/q/

## Documentation

- `PRD.md` - Product Requirements Document (detailed 5-day plan)
- `service/user/README.md` - User service documentation
- `service/room/CLAUDE.md` - Room service detailed guide
- `service/room/docs/user-integration.md` - User integration guide
- `service/room/docs/api-gateway-example.conf` - API Gateway examples

## Technology Stack Summary

| Component | Technology |
|-----------|-----------|
| Backend Framework | Go + Kratos v2 |
| Database | PostgreSQL 16 |
| ORM | GORM v2 |
| Authentication | JWT + bcrypt |
| Real-time | WebSocket (gorilla/websocket) |
| API Protocol | HTTP + gRPC + Protobuf |
| Dependency Injection | Wire |
| Frontend | Next.js + React (planned) |
| Containerization | Docker + docker-compose |

## Branch Strategy

- `dev-room` - Room service development + user service integration
- `feat/Microservice_for_Registration_with_Kratos_Dockerfile` - Original user service
- `main` - Production releases
