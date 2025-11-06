# Blog Go - Backend API

A lightweight Go backend API service for managing mailing list subscriptions.

## Overview

This project provides a REST API for handling mailing list operations, built with Go and designed to work with a blog platform hosted at zhisme.com.

## Tech Stack

- **Go**: 1.23.6
- **Chi Router**: v5.2.2 - Lightweight HTTP router
- **CORS**: v1.2.1 - Cross-Origin Resource Sharing support

## Project Structure

```
blog_go/
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

## License

See LICENSE file for details.
