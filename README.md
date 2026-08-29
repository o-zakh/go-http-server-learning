# Chirpy

A lightweight, scalable HTTP API server for a Twitter-like social media platform built with Go. Chirpy provides user authentication, post management, and premium features through webhook integrations.

## Features

- **User Management** - Registration and profile updates with secure password hashing using Argon2id
- **Authentication** - JWT-based access tokens with refresh token support and revocation
- **Chirps (Posts)** - Create, read, and delete posts with automatic profanity filtering
- **Chirpy Red Premium** - Webhook-based premium tier system for user upgrades
- **Database** - PostgreSQL with type-safe queries via SQLC and automated migrations via Goose
- **Admin Dashboard** - Metrics tracking and system monitoring endpoints
- **Security** - Bearer token authentication, password hashing, and secure token management

## Tech Stack

- **Language**: Go 1.26.5
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: Argon2id (alexedwards/argon2id)
- **Database Tools**: 
  - [Goose](https://github.com/pressly/goose) - Migration tool
  - [SQLC](https://sqlc.dev/) - SQL code generation
- **HTTP Server**: Go's standard `net/http` with custom routing
- **Utilities**: UUID generation, environment configuration, webhook handling

## Quick Start

### Prerequisites

- Go 1.26.5 or later
- PostgreSQL 12 or later
- Goose CLI (for database migrations)
- SQLC (for code generation)

### Installation

1. **Clone the repository**

```bash
git clone https://github.com/o-zakh/go-http-server-learning
cd go-http-server-learning
```

2. **Install dependencies**

```bash
go mod download
```

3. **Set up environment variables**

Create a `.env` file in the project root:

```env
DB_URL=postgres://username:password@localhost:5432/chirpy
TOKEN_SECRET=your-secret-key-for-jwt-tokens
POLKA_KEY=your-polka-webhook-secret-key
```

4. **Set up the database**

Create a PostgreSQL database:

```bash
createdb chirpy
```

Run migrations:

```bash
goose postgres "postgres://username:password@localhost:5432/chirpy" up
```

5. **Run the server**

```bash
go run .
```

The server will start on `http://localhost:8080`

## API Usage

### Health Check

```bash
GET /api/healthz
```

**Response**: `200 OK`

### User Management

#### Create User

```bash
POST /api/users
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response**: 
```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "created_at": "2024-01-15T10:30:00Z",
  "updated_at": "2024-01-15T10:30:00Z",
  "email": "user@example.com",
  "is_chirpy_red": false
}
```

#### Login

```bash
POST /api/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

**Response**: Returns access token and refresh token (1 hour expiration)

#### Refresh Access Token

```bash
POST /api/refresh
Authorization: Bearer <refresh_token>
```

#### Revoke Refresh Token

```bash
POST /api/revoke
Authorization: Bearer <refresh_token>
```

#### Update User Email/Password

```bash
PUT /api/users
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "email": "newemail@example.com",
  "password": "newpassword123"
}
```

### Chirps (Posts)

#### Create Chirp (Max 140 characters)

```bash
POST /api/chirps
Authorization: Bearer <access_token>
Content-Type: application/json

{
  "body": "Hello, Chirpy!"
}
```

**Note**: Profanity in the following words is automatically filtered: "kerfuffle", "sharbert", "fornax"

#### Get All Chirps

```bash
GET /api/chirps
```

**Query Parameters**:
- `author_id` - Filter by author UUID
- `sort` - Sort order: `asc` (oldest first) or `desc` (newest first, default)

#### Get Specific Chirp

```bash
GET /api/chirps/{chirpID}
```

#### Delete Chirp

```bash
DELETE /api/chirps/{chirpID}
Authorization: Bearer <access_token>
```

Only the chirp author can delete their own chirps.

### Admin Endpoints

#### View Metrics

```bash
GET /admin/metrics
```

Returns an HTML page displaying the number of times the `/app` endpoint has been accessed.

#### Reset Metrics

```bash
POST /admin/reset
```

Resets the metrics counter to zero.

### Premium Features

#### Polka Webhook (User Upgrade)

When users upgrade to Chirpy Red through Polka, the webhook endpoint receives:

```bash
POST /api/polka/webhooks
Authorization: ApiKey <polka_key>
Content-Type: application/json

{
  "event": "user.upgraded",
  "data": {
    "user_id": "550e8400-e29b-41d4-a716-446655440000"
  }
}
```

Upgrade events set the user's `is_chirpy_red` flag to true.

## Database Schema

The project uses PostgreSQL with Goose for migrations. Key tables:

- **users** - User accounts with email and hashed passwords
- **chirps** - User-generated posts with content and timestamps
- **refresh_tokens** - Session management with expiration and revocation tracking

Migration files are located in [sql/schema/](sql/schema/).

## Project Structure

```
.
├── main.go                 # Server entry point and route setup
├── chirps.go              # Chirp creation and retrieval handlers
├── users_edit.go          # User management handlers
├── auth_tokens.go         # Token refresh and revocation logic
├── polka.go               # Webhook handling for premium features
├── prof_filter.go         # Profanity filtering utility
├── json.go                # JSON response utilities
├── internal/
│   ├── auth/              # Authentication utilities (JWT, password hashing)
│   └── database/          # SQLC-generated database queries
├── sql/
│   ├── schema/            # Database migration files
│   └── queries/           # SQL query definitions for SQLC
├── assets/                # Static files
└── .github/
    └── prompts/           # Prompt templates for development
```

## Development

### Code Generation

When you modify SQL queries in [sql/queries/](sql/queries/), regenerate the database code:

```bash
sqlc generate
```

### Running Tests

Database unit tests are available in [internal/auth/unit_test.go](internal/auth/unit_test.go):

```bash
go test ./internal/...
```

### Debugging

The project includes debugging utilities for profiling and filtering logic located in [prof_filter.go](prof_filter.go).

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

For issues, questions, or suggestions, please open an [Issue](https://github.com/o-zakh/go-http-server-learning/issues) on GitHub.

## License

This project is provided as-is for educational and learning purposes.

## Maintainer

- [o-zakh](https://github.com/o-zakh)
