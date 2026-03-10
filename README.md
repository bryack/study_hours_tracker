# ⏱️ Study Hours Tracker - Mastering Hexagonal Architecture in Go

[![Go Version](https://img.shields.io/badge/Go-1.25.6-blue.svg)](https://golang.org)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15+-blue.svg)](https://www.postgresql.org)
[![Hexagonal Architecture](https://img.shields.io/badge/Architecture-Hexagonal-green.svg)](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software))
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A robust study-tracking application built as a deep-dive into **Clean Code** and **Hexagonal Architecture (Ports & Adapters)** in Go. This project demonstrates how to decouple core business logic from external dependencies (databases, CLIs, and web servers) for maximum testability and maintainability.

---

## 🚀 Key Technical Highlights

This project was built to master advanced Go development practices:

*   **Hexagonal Architecture**: Strict separation of concerns between the core domain and external adapters (CLI, Web, Database).
*   **Dependency Injection**: Using constructors to inject concrete adapter implementations into domain ports.
*   **Real-time Communication**: Utilizing **Gorilla WebSockets** for live Pomodoro session alerts in the browser.
*   **Automated Integration Testing**: Leverages **Testcontainers** to spin up real PostgreSQL instances for database tests.
*   **Concurrency**: Thread-safe Pomodoro timer implementation with phased alerts.
*   **Persistent Storage**: Efficient relational data handling with the **pgx** driver.

---

## 🏗️ Architecture: Ports & Adapters

The project structure is a textbook implementation of Hexagonal Architecture, ensuring the business logic remains "pure" and unaware of the infrastructure:

```text
├── domain/       # The Core: Business logic, entities, and Port interfaces
│   ├── pomodoro/ # Timer logic (Domain Service)
│   └── ...       # No external dependencies allowed here!
├── adapters/     # The Infrastructure: Concrete implementations of Ports
│   ├── cli/      # Interactive terminal interface
│   ├── database/ # PostgreSQL implementation (pgx)
│   ├── server/   # Web & WebSocket server (Gorilla)
│   └── pomodoro/ # Alerter implementation
└── cmd/          # Entry Points: Wiring adapters to the domain via DI
    ├── cli/      # CLI App launcher
    └── webserver/# Web App launcher
```

### Technical Decision: Domain Simplification
For this prototype, a **25-minute Pomodoro session** is automatically recorded as **1 study hour** in the database. This design choice simplifies the initial tracking logic while focusing on the architectural integrity of the system.

---

## 🛠️ Tech Stack & Dependencies

*   **Language**: Go 1.25.6
*   **Database**: PostgreSQL 15+
*   **Web Framework**: Standard `net/http` + Gorilla WebSockets
*   **Testing**: `testify` for assertions, `testcontainers-go` for DB integration.
*   **Frontend**: Vanilla HTML/JS (Embedded in the binary via `go:embed`).

---

## 🚦 Getting Started

### Prerequisites
*   Go 1.25.6+
*   Docker or Podman (for PostgreSQL and tests).

### 1. Start the Database
```bash
docker run --name study-postgres \
  -e POSTGRES_PASSWORD=password \
  -e POSTGRES_DB=study_tracker \
  -p 5432:5432 \
  -d postgres:15
```

### 2. Configure Environment
```bash
export DATABASE_URL="postgres://postgres:password@localhost:5432/study_tracker?sslmode=disable"
# (Optional) For Fedora/Podman users:
export TESTCONTAINERS_RYUK_DISABLED=true
```

### 3. Run the Applications
*   **CLI Mode**: `go run ./cmd/cli` (Try typing `math 2` or `pomodoro tdd`)
*   **Web Mode**: `go run ./cmd/webserver` (Visit `http://localhost:5000/study`)

---

## 🧪 Testing

The project maintains a high standard of quality through multi-layered testing:

```bash
# Run all tests (including integration tests with Testcontainers)
go test ./...

# Check test coverage
go test -cover ./...
```

---

## 🛤️ Roadmap & Future Improvements
- [ ] Add a **gRPC Adapter** for high-performance microservices interaction.
- [ ] Implement **OAuth2** for secure user authentication.
- [ ] Integrate **Prometheus metrics** for tracking session stats.

---

## 📫 Contact

*   **LinkedIn**: [www.linkedin.com/in/anna-nurgaleeva-ba9a6338](https://www.linkedin.com/in/anna-nurgaleeva-ba9a6338)
*   **Telegram**: [@bryacka](https://t.me/bryacka)

---

*This project was built with a commitment to clean code and architectural excellence.*
