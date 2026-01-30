# Study Hours Tracker

A simple, efficient application for tracking study time across different subjects. Built with Go using clean architecture principles, offering both CLI and web interfaces.

## Features

- **Study Session Recording** - Log hours spent studying specific subjects
- **Historical Data Retrieval** - Access past study records by subject  
- **Progress Reporting** - Generate comprehensive reports showing study patterns
- **Multiple Interfaces** - CLI for quick logging, web server for visual access
- **PostgreSQL Storage** - Reliable data persistence with full ACID compliance

## Quick Start

### Prerequisites

- Go 1.25.6 or later
- PostgreSQL database (local or containerized)

### Installation

```bash
# Clone the repository
git clone https://github.com/bryack/study_hours_tracker
cd study_hours_tracker

# Install dependencies
go mod download

# Run tests to verify setup
go test ./...
```

### Usage

#### Command Line Interface

```bash
# Build and run CLI
go build -o study-cli ./cmd/cli
./study-cli math 2    # Record 2 hours of math study
./study-cli physics 1 # Record 1 hour of physics study
```

#### Web Server

```bash
# Build and run web server
go build -o study-server ./cmd/webserver
./study-server

# Server runs on http://localhost:8080
# POST /subjects/{subject} - Record study hours
# GET /subjects/{subject} - Get total hours for subject
# GET /report - Get complete study report
```

## Architecture

This project demonstrates **Hexagonal Architecture** (Ports & Adapters):

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│   CLI Adapter   │    │   HTTP Adapter   │    │  Database       │
│                 │    │                  │    │  Adapter        │
└─────────┬───────┘    └─────────┬────────┘    └─────────┬───────┘
          │                      │                       │
          └──────────────────────┼───────────────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │     Domain Layer        │
                    │  (Business Logic)       │
                    │                         │
                    │  • StudyActivity        │
                    │  • Report               │
                    │  • SubjectStore (port)  │
                    └─────────────────────────┘
```

### Key Components

- **Domain Layer** (`domain/`) - Pure business logic, no external dependencies
- **Ports** (`store/`) - Interface definitions for external services
- **Adapters** (`adapters/`) - Concrete implementations of external interfaces
- **Entry Points** (`cmd/`) - Application launchers for different interfaces

## Development

### Running Tests

```bash
# Unit tests
go test ./...

# Integration tests with race detection
go test -race ./...

# Test coverage
go test -cover ./...
```

### Database Setup

The application uses PostgreSQL. For development, you can use Docker:

```bash
# Start PostgreSQL container
docker run --name study-postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=study_tracker \
  -p 5432:5432 \
  -d postgres:15

# Set environment variables
export DATABASE_URL="postgres://postgres:password@localhost:5432/study_tracker?sslmode=disable"
```

### Code Quality

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Run linter (if golangci-lint is installed)
golangci-lint run
```

## API Reference

### REST Endpoints

#### Record Study Hours
```http
POST /subjects/{subject}
Content-Type: application/json

{
  "hours": 2
}
```

#### Get Subject Hours
```http
GET /subjects/{subject}

Response:
{
  "subject": "math",
  "hours": 5
}
```

#### Get Study Report
```http
GET /report

Response:
[
  {
    "subject": "math",
    "hours": 5
  },
  {
    "subject": "physics", 
    "hours": 3
  }
]
```

## Project Structure

```
study_hours_tracker/
├── cmd/                    # Application entry points
│   ├── cli/               # CLI application
│   └── webserver/         # HTTP server application
├── domain/                # Core business logic
├── store/                 # Interface definitions (ports)
├── adapters/              # External interface implementations
│   ├── cli/              # CLI implementation
│   ├── database/         # PostgreSQL implementation
│   └── server/           # HTTP server implementation
├── testhelpers/          # Shared testing utilities
├── go.mod               # Go module definition
└── README.md            # This file
```

## Contributing

1. Follow Go conventions and project architecture patterns
2. Write tests for new functionality
3. Ensure all tests pass before submitting changes
4. Use meaningful commit messages
5. Update documentation as needed

## License

This project is for educational purposes demonstrating clean architecture in Go.

## Technology Stack

- **Language:** Go 1.25.6
- **Database:** PostgreSQL with pgx driver
- **Testing:** Testify + Testcontainers
- **Architecture:** Hexagonal (Ports & Adapters)
- **HTTP:** Standard library net/http
