# Stage 1: Base - Common foundation for all environments
FROM golang:1.25.1-alpine AS base

# Install basic dependencies
RUN apk add --no-cache git ca-certificates

# Set working directory
WORKDIR /app

# Copy dependency files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Stage 2: Development - Hot reload with Air
FROM base AS dev

# Install Air for hot-reload
RUN go install github.com/air-verse/air@latest

# Install Swag for Swagger generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code (will be mounted as volume in docker-compose)
COPY . .

# Expose application port
EXPOSE 8080

# Run Air
CMD ["air", "-c", ".air.toml"]

# Stage 3: Builder - Build optimized binary for stage/prod
FROM base AS builder

# Accept build environment as argument
ARG BUILD_ENV=prod

# Install Swag for Swagger generation
RUN go install github.com/swaggo/swag/cmd/swag@latest

# Copy source code
COPY . .

# Generate Swagger docs
RUN swag init -g cmd/api/main.go -o docs

# Build binary with appropriate build tag
RUN CGO_ENABLED=0 GOOS=linux go build \
  -tags ${BUILD_ENV} \
  -ldflags="-s -w" \
  -o /app/api \
  ./cmd/api

# Stage 4: Staging - Minimal runtime for staging
FROM alpine:latest AS stage

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 appuser && \
  adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/api .

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose application port
EXPOSE 8080

# Run the binary
CMD ["./api"]

# Stage 5: Production - Minimal runtime for production
FROM alpine:latest AS prod

# Install CA certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN addgroup -g 1000 appuser && \
  adduser -D -u 1000 -G appuser appuser

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/api .

# Change ownership
RUN chown -R appuser:appuser /app

# Switch to non-root user
USER appuser

# Expose application port
EXPOSE 8080

# Run the binary
CMD ["./api"]