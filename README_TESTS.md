# Test Documentation

This document describes the test suite for the PostAPI project.

## Running Tests

### Run all tests
```bash
go test ./...
```

### Run tests with coverage
```bash
go test ./... -cover
```

### Run tests with detailed coverage report
```bash
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

### Run tests for a specific package
```bash
go test ./internal/application
go test ./internal/middleware
go test ./internal/domain
```

### Run tests in verbose mode
```bash
go test ./... -v
```

### Run a specific test
```bash
go test ./internal/application -run TestJWTService_GenerateToken
```

## Test Structure

### Application Layer Tests (`internal/application`)

**jwt_service_test.go**
- `TestJWTService_GenerateToken`: Tests token generation with various usernames
- `TestJWTService_ValidateToken`: Tests token validation including invalid tokens
- `TestJWTService_TokenExpiration`: Tests token validity checks
- `TestJWTService_RoundTrip`: Tests complete token generation and validation cycle

**mappers_test.go**
- `TestMapUserToJson`: Tests user to JSON conversion
- `TestMapPostToJson`: Tests post to JSON conversion
- `TestMapFollowToJson`: Tests follow relationship to JSON conversion
- `TestMapProfileToJson`: Tests profile to JSON conversion

### Middleware Tests (`internal/middleware`)

**auth_test.go**
- `TestAuthMiddleware_MissingHeader`: Tests missing Authorization header
- `TestAuthMiddleware_InvalidFormat`: Tests invalid token format
- `TestAuthMiddleware_InvalidToken`: Tests invalid JWT tokens
- `TestAuthMiddleware_ValidToken`: Tests valid authentication flow
- `TestAuthMiddleware_ContextKey`: Tests context value storage and retrieval
- `TestAuthMiddleware_DifferentTokens`: Tests multiple users with different tokens
- `TestNewAuthMiddleware`: Tests middleware initialization

**response_test.go**
- `TestParse`: Tests JSON request body parsing
- `TestSendResponse`: Tests JSON response sending with various status codes
- `TestSendResponse_JSONEncoding`: Tests proper JSON encoding
- `TestSendResponse_Array`: Tests array response serialization
- `TestParse_EmptyBody`: Tests handling of empty request bodies

### Domain Layer Tests (`internal/domain`)

**models_test.go**
- `TestUser_ToResponse`: Tests User to UserResponse conversion
- `TestUserStructTags`: Verifies struct field accessibility
- `TestPostModel`: Tests Post model fields
- `TestProfileModel`: Tests Profile model fields
- `TestUserFollowModel`: Tests UserFollow model fields
- `TestJsonUserModel`: Verifies password exclusion from JSON representation
- `TestPostRequestModel`: Tests PostRequest model
- `TestProfileRequestModel`: Tests ProfileRequest model

## Test Coverage Goals

- **Application Layer**: 80%+ coverage
- **Middleware**: 90%+ coverage
- **Domain Models**: 95%+ coverage

## Testing Best Practices

1. **Isolation**: Each test should be independent and not rely on others
2. **Clear Names**: Test names should describe what is being tested
3. **Table-Driven**: Use table-driven tests for multiple similar scenarios
4. **Mock Dependencies**: Use mocks for external dependencies (database, services)
5. **Edge Cases**: Test both happy paths and error scenarios

## Future Test Additions

To complete the test suite, consider adding:

1. **Handler Tests**: Integration tests for HTTP handlers with mocked repositories
2. **Repository Tests**: Tests for database operations (requires test database)
3. **End-to-End Tests**: Full API integration tests
4. **Performance Tests**: Benchmark tests for critical paths
5. **Security Tests**: Test authentication, authorization, and input validation

## Mock Setup

For handler and integration tests, you'll need to create mocks for:
- `UserRepository`
- `PostRepository`
- `ProfileRepository`
- `UserFollowRepository`
- `JWTService`

Consider using a mocking library like `github.com/stretchr/testify/mock` or `github.com/golang/mock`.

## Example: Running Specific Test Suites

```bash
# Test JWT functionality
go test ./internal/application -run JWT -v

# Test middleware
go test ./internal/middleware -v

# Test with race detection
go test ./... -race

# Generate coverage for CI/CD
go test ./... -coverprofile=coverage.out -covermode=atomic
```

## Continuous Integration

Add to your `.github/workflows/test.yml`:

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - run: go test ./... -v -race -coverprofile=coverage.out
      - run: go tool cover -func=coverage.out
```
