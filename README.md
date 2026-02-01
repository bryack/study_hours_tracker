# Study Hours Tracker

Track study time across subjects. Built with Go using hexagonal architecture.

## Quick Start

```bash
# CLI
go build -o study-cli ./cmd/cli
./study-cli math 2

# Web Server (http://localhost:5000)
go build -o study-server ./cmd/webserver
./study-server
```

## API

```bash
# Record hours
POST /tracker/math?hours=2    # 202 Accepted

# Get hours
GET /tracker/math             # Returns: 5

# Get report
GET /report                   # Returns: [{"subject":"math","hours":5}]
```

### Validation

- Subject cannot be empty → `400 Bad Request`
- Hours must be positive integer → `400 Bad Request`

## Development

```bash
# Setup PostgreSQL
docker run -d -p 5432:5432 \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=study_tracker \
  postgres:15

export DATABASE_URL="postgres://postgres:password@localhost:5432/study_tracker?sslmode=disable"

# Test (Fedora/Podman users: disable Ryuk)
export TESTCONTAINERS_RYUK_DISABLED=true
go test ./...

# Format
go fmt ./...
```

## Architecture

**Hexagonal (Ports & Adapters)**

```
cmd/          → Entry points (CLI, Web)
domain/       → Business logic
store/        → Port interfaces
adapters/     → Implementations (CLI, Server, Database)
testhelpers/  → Test utilities
```

**Stack:** Go 1.25.6 • PostgreSQL • Testify • Testcontainers
