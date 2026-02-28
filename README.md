# keygen-service

Microservice for generating and validating activation keys.

![Go](https://img.shields.io/badge/Go-1.22+-00ADD8)
![Gin](https://img.shields.io/badge/Gin-1.10-blue)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-blue)
![License](https://img.shields.io/badge/License-MIT-lightgrey)

## Overview

A lean microservice for managing activation keys. Suitable for SaaS products, plugin stores, or any system requiring key-based access control. Keys support expiration, usage limits, per-IP logging, and revocation.

## Stack

- Go 1.22+
- Gin
- GORM (PostgreSQL)
- JWT (golang-jwt/jwt)
- bcrypt
- Custom in-memory rate limiter (no external dependency)

## Structure

```
cmd/
└── main.go             application entry point
config/
└── config.go           environment loading
internal/
├── handlers/           key_handler, auth_handler
├── middleware/         auth, ratelimit
├── models/             User, Key, KeyUsageLog
├── repository/         key_repository, user_repository
└── services/           key_service, auth_service
```

## Setup

```bash
git clone https://github.com/m4trixdev/keygen-service.git
cd keygen-service
cp .env.example .env
# Fill in DATABASE_URL and JWT_SECRET
```

Install dependencies and run:

```bash
go mod tidy
go run ./cmd/main.go
```

### With Docker

```bash
docker build -t keygen-service .
docker run -p 8080:8080 --env-file .env keygen-service
```

## API

### Auth

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/auth/login` | Get JWT token | No |
| POST | `/api/v1/auth/register` | Create user | Admin |

### Keys

| Method | Endpoint | Description | Auth |
|--------|----------|-------------|------|
| POST | `/api/v1/keys` | Generate a key | Admin |
| GET | `/api/v1/keys` | List all keys | Admin |
| POST | `/api/v1/keys/validate` | Validate a key | User |
| DELETE | `/api/v1/keys/:id` | Revoke a key | Admin |

### Example: Generate a key

```json
POST /api/v1/keys
Authorization: Bearer <token>

{
  "label": "Plugin License",
  "max_uses": 1,
  "expires_at": "2025-12-31T00:00:00Z"
}
```

Response:

```json
{
  "id": "...",
  "value": "a1b2c3d4-e5f6a7b8-c9d0e1f2-a3b4c5d6",
  "label": "Plugin License",
  "max_uses": 1,
  "uses": 0,
  "revoked": false,
  "expires_at": "2025-12-31T00:00:00Z",
  "created_at": "..."
}
```

### Example: Validate a key

```json
POST /api/v1/keys/validate
Authorization: Bearer <token>

{ "key": "a1b2c3d4-e5f6a7b8-c9d0e1f2-a3b4c5d6" }
```

Response:

```json
{ "valid": true, "key": { ... } }
```

## Key format

Keys are generated using `crypto/rand` (cryptographically secure) and formatted as:

```
xxxxxxxx-xxxxxxxx-xxxxxxxx-xxxxxxxx
```

## Environment variables

| Variable | Description | Default |
|----------|-------------|---------|
| `DATABASE_URL` | PostgreSQL DSN | required |
| `JWT_SECRET` | JWT signing secret | required |
| `PORT` | HTTP port | `8080` |
| `RATE_LIMIT_PER_MIN` | Max requests per IP per minute | `60` |

## Author

**M4trixDev** — [github.com/m4trixdev](https://github.com/m4trixdev)

## License

MIT
