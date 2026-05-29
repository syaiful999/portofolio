# Moyo - Microservices Portfolio

A production-grade microservices architecture built with Go, featuring gRPC communication, API gateway pattern, and PostgreSQL.

## Architecture

```
┌──────────────────┐         ┌──────────────────┐
│  Gateway Service │  gRPC   │  Master Service  │
│  (Fiber + gRPC-  │────────▶│  (Go Micro v4)   │
│   Gateway)       │         │                  │
│  :8080 / :8081   │         │  :8095 / :8096   │
└──────────────────┘         └────────┬─────────┘
                                      │
                             ┌────────▼─────────┐
                             │   PostgreSQL      │
                             │   + Redis Cache   │
                             └──────────────────┘
```

## Services

| Service | Description | Port |
|---------|-------------|------|
| `gateway-service` | API Gateway — REST to gRPC proxy, JWT auth, RBAC, caching | 8080 (gRPC-Gateway), 8081 (Fiber) |
| `master-service` | Core business logic — User management, Enum/lookup data | 8095 (gRPC), 8096 (HTTP) |

## Tech Stack

- **Language**: Go 1.24
- **Gateway**: GoFiber v2 + gRPC-Gateway v2
- **Microservice**: Go Micro v4 (gRPC)
- **Database**: PostgreSQL with golang-migrate
- **Cache**: Redis + In-memory (go-cache)
- **Auth**: JWT + bcrypt password hashing
- **Object Storage**: MinIO
- **Monitoring**: Prometheus + VictoriaMetrics
- **Broker**: NATS
- **Registry**: Kubernetes / mDNS (local)
- **Container**: Docker (Alpine multi-stage builds)

## Getting Started

### Prerequisites

- Go 1.24+
- PostgreSQL 14+
- Redis
- Docker (optional)
- [golang-migrate](https://github.com/golang-migrate/migrate) CLI

### Quick Setup

```bash
# 1. Clone the repository
git clone <repository-url>
cd Portfolio

# 2. Setup master-service
cd master-service
cp .env.example .env       # Edit with your DB credentials
go mod tidy

# 3. Create database & run migrations
make db-create
make migrate-up

# 4. Run the service
go run main.go

# 5. In another terminal, setup gateway-service
cd ../gateway-service
cp .env.example .env       # Edit with your service config
go mod tidy
go run main.go
```

### Database Migrations

All migration commands are available via Make in `master-service/`:

```bash
make migrate-up         # Run all pending migrations
make migrate-down       # Rollback last migration
make migrate-down-all   # Rollback all migrations
make migrate-reset      # Drop & recreate everything
make migrate-version    # Show current migration version
make migrate-create NAME=add_new_table  # Create new migration files
make db-setup           # Create DB + run all migrations (fresh start)
```

**Environment variables** for database connection:

```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=your_username
export DB_PASSWORD=your_password
export DB_NAME=master_db
```

### Generate Proto

```bash
# Install protoc tools (run once)
make init

# Generate gRPC code
make proto
```

## Project Structure

```
├── gateway-service/
│   ├── cache/              # Redis + in-memory cache abstraction
│   ├── config/             # Environment configuration
│   ├── middleware/         # Auth, RBAC, logging middleware
│   ├── pkg/entities/       # Request/response entities
│   ├── proto/              # gRPC proto definitions & generated code
│   └── utils/              # JWT, logging, response helpers
│
└── master-service/
    ├── config/             # YAML + env configuration
    ├── database/
    │   ├── connection.go   # DB connection with pool config
    │   └── migrations/     # SQL migration files (golang-migrate)
    ├── pkg/
    │   ├── enum/           # Enum/lookup domain (handler → usecase → repo)
    │   └── user/           # User domain (handler → usecase → repo)
    ├── plugin/             # Logging (logrus) + file upload (MinIO)
    ├── server/             # Service initialization & handler registration
    └── utils/              # Shared utilities (validation, response, etc.)
```

## Database Schema

| Schema | Table | Description |
|--------|-------|-------------|
| `master` | `master_enum` | Lookup/enum values (roles, departments, locations, etc.) |
| `master` | `master_user` | User accounts with auth & profile data |
| `master` | `v_master_user` | View: users joined with role/dept/location names |
| `master` | `v_master_user_groupby_role` | View: user count grouped by role |
| `transaction` | `transact_redirect_page` | Password reset/activation links |
| `transaction` | `transact_outsource` | Outsource-department relationships |

## Default Credentials (Seed Data)

After running migrations, a default admin account is available:

| Field | Value |
|-------|-------|
| Username | `admin` |
| Email | `admin@moyo.local` |
| Password | `admin123` |
| Role | Super Admin |

> Change the default password immediately in non-development environments.

## License

Proprietary
