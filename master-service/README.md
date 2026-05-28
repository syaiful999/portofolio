# Moyo Master Service

Backend service for Moyo Master Service, built using the Go Micro framework.

## Tech Stack

- **Language**: Go
- **Framework**: Go Micro v4
- **Registry**: Kubernetes
- **Monitoring**: Prometheus, VictoriaMetrics
- **Configuration**: Godotenv

## Features

- Microservice architecture
- Kubernetes service registry integration
- Built-in monitoring wrappers (Prometheus, VictoriaMetrics)
- Environment variable configuration management

## Getting Started

### Prerequisites

- Go (latest stable version)
- Docker (optional)

### Installation

1. **Clone the repository**
   ```bash
   # git clone <repository-url>
   cd moyo-master-service
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Setup environment variables**
   ```bash
   cp .env.example .env
   ```

4. **Run the application**
   ```bash
   go run main.go
   ```

## License

[Proprietary]
