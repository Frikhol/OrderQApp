FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Install goose for migrations
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Copy the source code
COPY . .
COPY .env .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o auth_service ./cmd/auth/main.go

# Start a new stage from scratch
FROM alpine:latest

RUN apk --no-cache add ca-certificates postgresql-client go

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/auth_service .
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy the migrations directory
COPY --from=builder /app/migrations ./migrations

# Expose port
EXPOSE 9000

# Command to run
CMD ["./auth_service"] 