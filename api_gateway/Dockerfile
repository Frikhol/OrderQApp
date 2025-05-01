FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api_gateway ./cmd/api-gateway/main.go

# Start a new stage from scratch
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/api_gateway .
# Copy conf directory
COPY --from=builder /app/conf ./conf

# Create directories for volumes
RUN mkdir -p views static

# Expose port
EXPOSE 8080

# Command to run
CMD ["./api_gateway"] 