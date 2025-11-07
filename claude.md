# Blog Go - Backend API

A lightweight, well-tested Go backend API service for managing mailing list subscriptions with comprehensive test coverage and automated quality checks.

## Overview

This project provides a REST API for handling mailing list operations, built with Go and designed to work with a blog platform hosted at zhisme.com. The codebase follows clean architecture principles with 77-100% test coverage and strict linting standards.

## Features

- **RESTful API**: Simple mailing list subscription management
- **CSV Storage**: Lightweight data persistence with duplicate detection
- **Input Validation**: Email format and username validation
- **CORS Support**: Configured for specific origins
- **High Test Coverage**: 77-100% across all modules
- **CI/CD Integration**: Automated testing and linting on pull requests

## Tech Stack

- **Go**: 1.23.6
- **Chi Router**: v5.2.2 - Lightweight HTTP router
- **CORS**: v1.2.1 - Cross-Origin Resource Sharing support
- **Testing**: Standard library testing with black-box approach
- **CI/CD**: GitHub Actions with golangci-lint

## Project Structure

```
blog_go/
├── .github/
│   └── workflows/
│       └── pr-validation.yml # CI/CD pipeline (tests, linting, formatting)
├── cmd/
│   └── api/
│       └── main.go           # Application entry point
├── internal/
│   ├── api/
│   │   ├── server.go         # HTTP server configuration
│   │   ├── create_mailing_list.go
│   │   └── handlers/
│   │       └── mailing_list.go
│   ├── dto/
│   │   └── mailing_list.go   # Data transfer objects
│   ├── interfaces/
│   │   ├── repositories.go   # Repository interfaces
│   │   └── validators.go     # Validator interfaces
│   ├── repositories/
│   │   └── csv_mailing_list_repository.go
│   └── validators/
│       └── mailing_list.go
├── test/                     # Black-box tests (mirrors src structure)
│   ├── cmd/api/
│   ├── internal/api/
│   ├── internal/api/handlers/
│   ├── internal/repositories/
│   └── internal/validators/
├── .golangci.yml            # Linter configuration
├── go.mod
└── go.sum
```

## Getting Started

### Prerequisites

- Go 1.23.6 or higher

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd blog_go
```

2. Install dependencies:
```bash
go mod download
```

### Running the Server

```bash
go run cmd/api/main.go
```

The server will start on `localhost:3000`.

## API Endpoints

### Mailing List

- **POST** `/mailing_list` - Create/subscribe to mailing list

## Configuration

The API is configured to accept requests from:
- `http://localhost:1313` (local development)
- `https://zhisme.com/` (production)

Allowed HTTP methods: `POST`, `DELETE`

## Development

The project follows a clean architecture pattern with:
- **Handlers**: HTTP request handlers
- **DTOs**: Data transfer objects for API contracts
- **Repositories**: Data persistence layer
- **Validators**: Input validation logic
- **Interfaces**: Dependency injection contracts

### Code Quality

The codebase maintains high quality standards through:
- **Automated Linting**: golangci-lint checks on every PR
- **Code Formatting**: Enforced gofmt formatting
- **Error Handling**: Explicit error checking (no ignored errors)
- **Memory Optimization**: Optimized struct field alignment
- **No Shadowing**: Clean variable scoping

### CI/CD Pipeline

Pull requests are automatically validated with:
- Go formatting checks (`gofmt`)
- Comprehensive linting (`golangci-lint`)
- Full test suite execution
- Build verification

## Testing

### Running Tests

Run all tests:
```bash
go test ./test/...
```

Run tests with verbose output:
```bash
go test ./test/... -v
```

Run tests with coverage:
```bash
go test ./test/... -cover
```

### Test Structure

The project uses **black-box testing** following Go best practices:
- All tests are in the `test/` directory, separate from source code
- Test packages use `_test` suffix (e.g., `package api_test`)
- Tests only access exported (public) APIs
- Test directory structure mirrors source code structure

### Test Coverage

| Package | Coverage | Description |
|---------|----------|-------------|
| `internal/api` | 77.1% | HTTP server and routing tests |
| `internal/api/handlers` | 87.5% | Business logic and handler tests |
| `internal/repositories` | 80.9% | CSV repository and file operations |
| `internal/validators` | 100.0% | Input validation tests |

### Test Files

```
test/
├── cmd/api/
│   └── main_test.go                    # Integration tests
├── internal/api/
│   ├── server_test.go                  # Server initialization & routing
│   ├── create_mailing_list_test.go     # HTTP handler tests
│   └── handlers/
│       └── mailing_list_test.go        # Business logic tests
├── internal/repositories/
│   └── csv_mailing_list_repository_test.go  # Data persistence tests
└── internal/validators/
    └── mailing_list_test.go            # Validation logic tests
```

### Testing Philosophy

- **Integration Tests**: Test the full request/response cycle
- **Unit Tests**: Test individual components through their public interfaces
- **No Mocks**: Tests use real implementations for reliability
- **Isolation**: Each test runs independently with its own test data files

## License

See LICENSE file for details.
