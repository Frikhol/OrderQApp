# Build stage
FROM golang:1.24.2-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Final stage
FROM alpine:latest

WORKDIR /app

# Install necessary packages
RUN apk add --no-cache ca-certificates

# Copy the binary and necessary files
COPY --from=builder /app/main .
COPY --from=builder /app/views ./views
COPY --from=builder /app/conf ./conf
COPY --from=builder /app/static ./static

# Create necessary directories and set permissions
RUN mkdir -p /app/views/auth && \
    chmod -R 755 /app/views && \
    chmod -R 755 /app/static

# Expose port
EXPOSE 8080

# Run the application
CMD ["./main"] 