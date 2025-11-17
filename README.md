# Diva Server

A high-performance Go backend server that provides RESTful APIs and real-time services for the Diva ecosystem.

## Description

Diva Server is the core backend infrastructure for the Diva platform, handling authentication, data persistence, real-time communication, and business logic for all client applications.

## Requirements

- **Go**: 1.21 or later
- **Make**: GNU Make 3.8 or later
- **PostgreSQL**: 13.0 or later (for production)
- **Redis**: 6.0 or later (for caching)

## Installation

### From Source

```bash
git clone <repository-url>
cd diva-server
make install
```

### Docker

```bash
# Build Docker image
make docker-build

# Run with Docker Compose
make docker-up
```

## Building

### Using Make

```bash
# Build the server binary
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run integration tests
make test-integration

# Clean build artifacts
make clean

# Install dependencies
make deps

# Generate code
make generate

# Run linter
make lint

# Format code
make fmt

# Run server in development mode
make run
```

### Available Make Targets

- `build`: Build for the current platform
- `build-all`: Build for all supported platforms (Linux, macOS, Windows)
- `run`: Run the server in development mode
- `test`: Run unit tests
- `test-integration`: Run integration tests
- `test-coverage`: Generate test coverage report
- `clean`: Remove build artifacts
- `deps`: Download and install dependencies
- `generate`: Run code generation tools
- `lint`: Run Go linter
- `fmt`: Format Go source code
- `migrate-up`: Run database migrations
- `migrate-down`: Rollback database migrations
- `docker-build`: Build Docker image
- `docker-up`: Start services with Docker Compose

## Configuration

Create a `.env` file in the root directory:

```env
# Server
PORT=8080
HOST=localhost

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=diva
DB_USER=diva_user
DB_PASSWORD=diva_password

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=

# JWT
JWT_SECRET=your-secret-key
JWT_EXPIRES_IN=24h
```

## Usage

### Development

```bash
# Run in development mode
make run

# Or run directly
go run cmd/server/main.go
```

### Production

```bash
# Build for production
make build

# Run the binary
./bin/diva-server
```

## API Documentation

When the server is running, API documentation is available at:
- Swagger UI: `http://localhost:8080/swagger/index.html`
- OpenAPI JSON: `http://localhost:8080/swagger/doc.json`

## Development

### Project Structure

```
server/
```

### Adding New Endpoints

1. Define the handler in `internal/api/`
2. Add the route in the router setup
3. Add models in `internal/models/`
4. Write tests in the appropriate test files
5. Update API documentation

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass: `make test`
6. Run linting: `make lint`
7. Submit a pull request

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.