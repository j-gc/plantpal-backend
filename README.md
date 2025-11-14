# PlantPal Backend

A modular monolith backend for PlantPal application built with Go.

## Project Structure

```
plantpal-backend/
├── cmd/
│   └── api/                    # Application entrypoint
│       └── main.go
├── internal/
│   ├── modules/                # Business modules
│   │   ├── users/             # User management module
│   │   │   ├── domain/
│   │   │   ├── application/
│   │   │   └── infrastructure/
│   │   │       ├── http/
│   │   │       └── persistence/
│   │   └── notifications/     # Notification module
│   │       ├── domain/
│   │       ├── application/
│   │       └── infrastructure/
│   │           ├── http/
│   │           └── persistence/
│   └── shared/                # Shared components
│       ├── config/            # Configuration management
│       ├── database/          # Database connection
│       ├── domain/            # Shared domain models
│       ├── errors/            # Error handling
│       ├── events/            # Event bus for inter-module communication
│       ├── infrastructure/    # Shared infrastructure
│       ├── logger/            # Logging
│       └── middleware/        # HTTP middleware
├── pkg/                       # Public libraries
├── config/                    # Configuration files
├── migrations/                # Database migrations
├── scripts/                   # Build and deployment scripts
├── .env.example              # Example environment variables
├── .gitignore
├── go.mod
└── README.md
```

## Architecture

This project follows a **Modular Monolith** architecture with **Clean Architecture** principles:

- **Domain Layer**: Business logic and entities (no external dependencies)
- **Application Layer**: Use cases and application services
- **Infrastructure Layer**: External concerns (HTTP, database, etc.)

### Module Communication

Modules communicate through:
1. **HTTP APIs**: For external communication
2. **Event Bus**: For asynchronous inter-module communication
3. **Shared Interfaces**: For synchronous communication (use sparingly)

### Key Principles

- Each module is independent and can be extracted into a microservice
- Dependencies point inward (Infrastructure → Application → Domain)
- Shared code is minimal and well-defined
- Database per module (logical separation)

## Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 14 or higher

### Installation

1. Clone the repository:
```bash
git clone https://github.com/V01D0/plantpal-backend.git
cd plantpal-backend
```

2. Copy the environment file:
```bash
cp .env.example .env
```

3. Update the `.env` file with your configuration

4. Install dependencies:
```bash
go mod download
```

5. Run migrations:
```bash
# Add migration tool and run migrations
```

6. Start the server:
```bash
go run cmd/api/main.go
```

The server will start on `http://localhost:8080`

## Development

### Running Tests

```bash
go test ./...
```

### Running with Hot Reload

```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with air
air
```

### Database Migrations

```bash
# Create a new migration
# migrate create -ext sql -dir migrations -seq migration_name

# Run migrations
# migrate -path migrations -database "postgresql://user:password@localhost:5432/plantpal?sslmode=disable" up

# Rollback
# migrate -path migrations -database "postgresql://user:password@localhost:5432/plantpal?sslmode=disable" down
```

## API Documentation

API documentation will be available at `/api/docs` when the server is running.

## Contributing

1. Create a feature branch
2. Make your changes
3. Write tests
4. Submit a pull request

## License

MIT
