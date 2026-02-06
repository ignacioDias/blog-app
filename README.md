# PostAPI

A RESTful API for a social media platform built with Go, featuring user authentication, posts, profiles, and follow relationships.

## Features

- 🔐 JWT-based authentication
- 👤 User registration and login
- 📝 Create, read, update, and delete posts
- 👥 User profiles with customizable information
- 🔗 Follow/unfollow users
- 📊 View followers and following lists

## Tech Stack

- **Language**: Go 1.21+
- **Router**: Gorilla Mux
- **Database**: PostgreSQL
- **Authentication**: JWT (golang-jwt)
- **Database Driver**: sqlx, pq
- **Password Hashing**: bcrypt

## Prerequisites

- Go 1.21 or higher
- PostgreSQL 12 or higher
- Git

## Installation

1. Clone the repository:
```bash
git clone https://github.com/yourusername/postapi.git
cd postapi
```

2. Install dependencies:
```bash
go mod download
```

3. Set up PostgreSQL database:
```bash
createdb postgres
```

4. Configure database connection in `internal/infrastructure/persistence/database.go`:
```go
dbUsername = "postgres"
dbPassword = "postgres"
dbHost     = "localhost"
dbPort     = "5432"
dbTable    = "postgres"
```

## Running the Application

```bash
go run cmd/main.go
```

The server will start on `http://localhost:8080`

## API Endpoints

### Authentication

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/register` | Register a new user | No |
| POST | `/api/login` | Login and get JWT token | No |

### Users

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/users/{username}` | Get user by username | No |

### Posts

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/posts` | Create a new post | Yes |
| GET | `/api/posts/{post_id}` | Get a specific post | No |
| PATCH | `/api/posts/{post_id}` | Update a post | Yes |
| DELETE | `/api/posts/{post_id}` | Delete a post | Yes |
| GET | `/api/{username}/posts` | Get all posts by a user | No |

### Profiles

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| GET | `/api/{username}` | Get user profile | No |
| POST | `/api/{username}` | Create user profile | Yes |
| PATCH | `/api/{username}` | Update user profile | Yes |

### Follow/Unfollow

| Method | Endpoint | Description | Auth Required |
|--------|----------|-------------|---------------|
| POST | `/api/follow/{username}` | Follow a user | Yes |
| DELETE | `/api/unfollow/{username}` | Unfollow a user | Yes |
| GET | `/api/{username}/followers` | Get user's followers | No |
| GET | `/api/{username}/following` | Get users being followed | No |

## Authentication

Protected endpoints require a JWT token in the Authorization header:

```bash
Authorization: Bearer <your_jwt_token>
```

## Example Requests

### Register a user
```bash
curl -X POST http://localhost:8080/api/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "john@example.com",
    "password": "securepassword123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "password": "securepassword123"
  }'
```

### Create a post
```bash
curl -X POST http://localhost:8080/api/posts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <your_token>" \
  -d '{
    "title": "My First Post",
    "content": "This is the content of my post"
  }'
```

## Testing

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test ./... -cover
```

Run tests in verbose mode:
```bash
go test ./... -v
```

**Current Test Coverage:**
- Application Layer: 82.6%
- Domain Layer: 100%
- Middleware: 96.0%

For detailed testing documentation, see [README_TESTS.md](README_TESTS.md)

## Project Structure

```
postapi/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── application/            # Business logic and use cases
│   │   ├── jwt_service.go
│   │   ├── post_usecase.go
│   │   ├── profile_usecase.go
│   │   └── user_usecase.go
│   ├── domain/                 # Domain models and interfaces
│   │   ├── post.go
│   │   ├── profile.go
│   │   ├── repositories.go
│   │   ├── user.go
│   │   └── user_follows.go
│   ├── infrastructure/         # External implementations
│   │   ├── handlers/          # HTTP handlers
│   │   ├── httpserver/        # Server and router setup
│   │   └── persistence/       # Database repositories
│   └── middleware/            # HTTP middleware
│       ├── auth.go
│       └── response.go
├── go.mod
└── go.sum
```

## Architecture

This project follows **Clean Architecture** principles:

- **Domain Layer**: Core business entities and repository interfaces
- **Application Layer**: Business logic, use cases, and services
- **Infrastructure Layer**: External implementations (database, HTTP)
- **Middleware**: Cross-cutting concerns (authentication, logging)

## Security Notes

⚠️ **Important for Production:**

1. Change the JWT secret key from `"secret-key"` to a strong, random secret
2. Use environment variables for sensitive configuration
3. Enable HTTPS/TLS
4. Implement rate limiting
5. Add input validation and sanitization
6. Use prepared statements (already implemented via sqlx)
