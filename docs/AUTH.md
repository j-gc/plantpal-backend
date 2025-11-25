# Authentication Endpoints

## Overview

User registration and login endpoints are implemented using JWT tokens.

## Endpoints

### POST /api/v1/auth/register

Register a new user account.

**Request Body:**
```json
{
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@example.com",
  "password": "SecurePass123"
}
```

**Success Response (201 Created):**
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "email": "john.doe@example.com",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request body or missing required fields
- `409 Conflict` - Email already in use

### POST /api/v1/auth/login

Authenticate and receive a JWT access token.

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "SecurePass123"
}
```

**Success Response (200 OK):**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "email": "john.doe@example.com",
    "first_name": "John",
    "last_name": "Doe"
  }
}
```

**Error Responses:**
- `400 Bad Request` - Invalid request body
- `401 Unauthorized` - Invalid email or password

## Testing

### Using curl

**Register:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "first_name": "John",
    "last_name": "Doe",
    "email": "john.doe@example.com",
    "password": "SecurePass123"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "SecurePass123"
  }'
```

### Using the test script

```bash
./scripts/test_auth.sh
```

## Implementation Details

### Architecture Layers

1. **Domain** (`internal/modules/authdomain/`)
   - User entity
   - UserRepository interface

2. **Application** (`internal/modules/authapplication/`)
   - Business logic for registration and login
   - Input/output DTOs
   - Password hashing and token generation

3. **Infrastructure**
   - **HTTP** (`internal/modules/authinfrastructure/http/`)
     - Gin handlers for auth endpoints
   - **Persistence** (`internal/modules/authinfrastructure/persistence/`)
     - PostgreSQL implementation of UserRepository
   - **Security** (`internal/modules/authinfrastructure/security/`)
     - Bcrypt password hashing

4. **Shared** (`internal/shared/`)
   - JWT token issuer (`internal/shared/jwt/`)
   - Configuration
   - Logger

### Security

- Passwords are hashed using bcrypt with default cost (10)
- JWT tokens are signed with HS256 algorithm
- Tokens expire after 24 hours
- JWT secret is configured via `JWT_SECRET` environment variable

### Database Schema

The `users` table is created via migration:
```sql
CREATE TABLE users (
    id uuid PRIMARY KEY DEFAULT gen_random_uuid(),
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);
```

Run migrations with:
```bash
make migrate-up
```
