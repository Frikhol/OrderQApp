# Order System

A microservices-based order management system with Beego frontend and gRPC backend.

## Architecture

The system consists of the following services:

- **API Gateway**: HTTP proxy service that handles client requests and forwards them to appropriate gRPC services
- **Order Service**: Handles order-related operations
- **Executor Service**: Manages order execution
- **Notification Service**: Handles system notifications
- **Auth Service**: Manages authentication and authorization

## Prerequisites

- Go 1.16 or later
- Protocol Buffers compiler (protoc)
- Docker and Docker Compose (for deployment)

## Setup

1. Install dependencies:
```bash
go mod tidy
```

2. Generate gRPC code:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/order.proto
```

3. Run services:
```bash
# Run API Gateway
go run cmd/api-gateway/main.go

# Run Order Service
go run cmd/order-service/main.go

# Run other services similarly
```

## Development

- Frontend: `web/` directory contains Beego web application
- Backend: `internal/` directory contains gRPC service implementations
- Common packages: `pkg/` directory contains shared utilities

## Deployment

See `deploy/` directory for Docker and Kubernetes configurations.

## License

MIT 