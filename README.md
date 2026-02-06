# Study Hours Tracker

Track study time across subjects. Built with Go using hexagonal architecture.

## Quick Start

```bash
# CLI - Manual recording
go build -o study-cli ./cmd/cli
./study-cli
# Interactive session starts:
# Type: math 2
# Type: physics 1
# Type: quit

# CLI - Pomodoro timer (25 minutes hardcoded)
# In interactive session, type:
# pomodoro tdd

# Web Server (http://localhost:5000)
go build -o study-server ./cmd/webserver
./study-server
```

## CLI Features

### Interactive Session
```bash
./study-cli
# Starts interactive session with greeting
# Type commands and press Enter
# Type 'quit' to exit gracefully
```

### Manual Recording
```bash
# In interactive session:
math 2        # Record 2 hours of math study
physics 1     # Record 1 hour of physics study
quit          # Exit the program
```

### Pomodoro Timer
```bash
# In interactive session:
pomodoro tdd  # Start 25-minute focused session for TDD
# Alerts during session:
# - 0 min: "Session started. Stay focused!"
# - 12 min: "Halfway there! Keep it up."
# - 25 min: "Time's up! Recording your hour..."
# Automatically records 1 hour to database
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
domain/       → Business logic & port interfaces
adapters/     → Implementations (CLI, Server, Database, Pomodoro)
testhelpers/  → Test utilities
```

**Stack:** Go 1.25.6 • PostgreSQL • Testify • Testcontainers
