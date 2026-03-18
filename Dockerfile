# Production Dockerfile
FROM golang:1.26.1-alpine AS build

# Install build dependencies
RUN apk add --no-cache curl libstdc++ libgcc

WORKDIR /app

# Copy go mod files and download dependencies
COPY go.mod go.sum ./
RUN go mod download && go mod verify

# Copy source code
COPY cmd ./cmd
COPY internal ./internal
COPY server ./server
COPY storage ./storage
COPY concurrency ./concurrency
COPY Makefile ./
COPY sqlc.yaml ./

# Copy env file for production
COPY .env ./

# Build the application
RUN go build -o ./bin/diva-server ./cmd/server/main.go

# Production stage
FROM alpine:3.20.1 AS prod

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Create non-root user
RUN addgroup -g 1001 -S diva && \
    adduser -u 1001 -S diva -G diva

# Copy binary from build stage
COPY --from=build /app/bin/diva-server /app/bin/diva-server

# Change ownership to non-root user
RUN chown -R diva:diva /app

USER diva

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:${PORT}/health || exit 1
