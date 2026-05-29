# Moyo Gateway Service

API Gateway for the Moyo microservices platform. Handles REST-to-gRPC translation, authentication, authorization, and response caching.

## Ports

| Port | Handler | Purpose |
|------|---------|---------|
| 8080 | gRPC-Gateway (mux) | REST API proxy to backend gRPC services |
| 8081 | GoFiber | Additional HTTP endpoints, health checks |

## Features

- **gRPC-Gateway**: Automatically translates REST/JSON requests to gRPC calls
- **JWT Authentication**: Token validation and user context extraction
- **RBAC Authorization**: Role-based access control middleware
- **Caching**: Redis (distributed) + in-memory (local) with configurable TTL
- **Request Logging**: Structured request/response logging

## Getting Started

### Prerequisites

- Go 1.23+
- Running `master-service` instance
- Redis (optional, for distributed cache)

### Setup

```bash
cp .env.example .env    # Configure service hosts and credentials
go mod tidy
go run main.go
```

### Configuration

See `.env.example` for all available environment variables:

- `SERVICE_*` — Service name, port, version
- `CACHE_*` — Redis connection settings
- `MASTER_DATA_SERVICE_HOST` — gRPC address of master-service

### Generate Proto

```bash
make init   # Install protoc plugins (once)
make proto  # Regenerate gRPC gateway code
```

## Architecture

```
Client (HTTP/REST)
    │
    ▼
┌─────────────────────────────────┐
│         Gateway Service         │
├─────────────────────────────────┤
│  Middleware: Auth → RBAC → Log  │
├─────────────────────────────────┤
│  Cache Layer (Redis/Memory)     │
├─────────────────────────────────┤
│  gRPC-Gateway (proto → REST)    │
└──────────────┬──────────────────┘
               │ gRPC
               ▼
        Backend Services
```

## License

Proprietary
