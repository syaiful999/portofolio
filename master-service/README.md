# Moyo Master Service

Core backend service for the Moyo platform. Manages users, roles, enums/lookups, and business logic via gRPC.

## Tech Stack

- **Language**: Go 1.24
- **Framework**: Go Micro v4 (gRPC)
- **Database**: PostgreSQL + golang-migrate
- **Monitoring**: Prometheus, VictoriaMetrics
- **Storage**: MinIO (file uploads)
- **Logging**: Logrus (structured)

## Features

- Clean architecture: Handler → UseCase → Repository
- Kubernetes service registry (mDNS for local dev)
- Database migrations with rollback support
- Password policy enforcement (expiry, lockout, failed attempts)
- Excel import/export for enum data

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL 14+
- [golang-migrate CLI](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate)

### Setup

```bash
# Install dependencies
go mod tidy

# Configure environment
cp .env.example .env

# Install migration tool
make init

# Create database and run migrations
make db-setup

# Run the service
go run main.go
```

### Database Commands

```bash
make migrate-up                      # Apply all pending migrations
make migrate-down                    # Rollback last migration
make migrate-down-all                # Rollback all migrations
make migrate-reset                   # Drop everything & re-migrate
make migrate-version                 # Show current version
make migrate-create NAME=add_table   # Create new migration pair
make db-create                       # Create the database
make db-drop                         # Drop the database
make db-setup                        # Create DB + migrate (fresh start)
```

### Proto Generation

```bash
make init   # Install protoc plugins (once)
make proto  # Generate gRPC service code
```

## Project Structure

```
├── config/             # YAML + env configuration structs
├── database/
│   ├── connection.go   # DB connection pool management
│   └── migrations/     # SQL migration files
├── pkg/
│   ├── enum/           # Enum/lookup domain
│   │   ├── handler.go
│   │   ├── usecase.go
│   │   ├── repository.go
│   │   └── models.go
│   └── user/           # User domain
│       ├── handler.go
│       ├── usecase.go
│       ├── repository.go
│       └── models.go
├── plugin/             # Logrus logger + MinIO uploader
├── server/             # Service bootstrap & handler registration
└── utils/              # Shared helpers (validation, response, etc.)
```

## License

Proprietary
