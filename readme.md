# tahrir-go

A Go rewrite of [Tahrir](https://github.com/fedora-infra/tahrir), the open badge assertion server originally written in Python. This project exposes a JSON REST API for managing persons, badges, and badge assertions, backed by a PostgreSQL database.

## Features

- **REST API** — Manage persons, badges, and badge assertions over HTTP.
- **Event-driven rule engine** — Automatically award milestone badges when a person reaches a configurable assertion threshold.
- **GORM + PostgreSQL** — Type-safe ORM-backed persistence with a well-defined model layer.
- **Live reload** — Hot-rebuild during development via [Air](https://github.com/air-verse/air).
- **Docker support** — Multi-stage Dockerfile for lean production images.

## Project Structure

```
tahrir-go/
├── cmd/
│   └── server/
│       └── main.go          # Application entrypoint
├── internal/
│   ├── db/                  # Database connection and query helpers
│   ├── handlers/            # HTTP handler functions
│   ├── models/              # GORM model definitions (Person, Badge, Assertion)
│   └── rules/               # Rule engine interface and implementations
│       ├── rules.go         # Event types and Rule interface
│       ├── milestone_rule.go # MilestoneRule: awards a badge after N assertions
│       └── workers/
│           └── worker.go    # Background worker that processes rule events
├── Dockerfile
├── Makefile
└── .env                     # Local environment configuration
```

## Prerequisites

- [Go](https://go.dev/dl/) 1.24+
- [PostgreSQL](https://www.postgresql.org/) (or a compatible container)
- [Air](https://github.com/air-verse/air) — for live reload during development
- [golangci-lint](https://golangci-lint.run/) — for linting
- [Docker](https://www.docker.com/) or [Podman](https://podman.io/) — optional, for containers

## Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/<your-org>/tahrir-go.git
cd tahrir-go
```

### 2. Configure the environment

Copy the example environment file and edit it to match your local PostgreSQL setup:

```bash
cp .env .env.local
```

```ini
# .env
DB_HOST=localhost
DB_PORT=5432
DB_USER=tahrir
DB_PASSWORD=tahrir
DB_NAME=tahrir
```

### 3. Start the database

If you have an existing PostgreSQL container named `tahrir-pg`:

```bash
make db-start
```

Or spin up a fresh one with Docker/Podman:

```bash
docker run -d \
  --name tahrir-pg \
  -e POSTGRES_USER=tahrir \
  -e POSTGRES_PASSWORD=tahrir \
  -e POSTGRES_DB=tahrir \
  -p 5432:5432 \
  postgres:16-alpine
```

### 4. Install dependencies

```bash
make deps
```

### 5. Run the server (with live reload)

```bash
make server
```

The API will be available at `http://localhost:8080`.

To run without live reload:

```bash
go run ./cmd/server
```

## API Reference

All responses are JSON. The server listens on `:8080`.

### Persons

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/persons` | List persons (max 5) |
| `GET` | `/persons/{nickname}` | Find person by nickname |
| `GET` | `/persons/id/{id}` | Find person by numeric ID |
| `GET` | `/persons/nickname/{person_nickname}/badges` | List all badge assertions for a person |

### Badges

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/badges` | List badges (max 5) |
| `GET` | `/badges/{id}` | Get a badge by ID |
| `POST` | `/badges` | Create a new badge |

**Create badge — request body:**
```json
{
  "ID": "badge-unique-id",
  "Name": "First Contribution",
  "Image": "https://example.com/badge.png",
  "Description": "Awarded for making a first contribution.",
  "Criteria": "Submit one pull request.",
  "IssuerID": 1,
  "Tags": "contribution,open-source"
}
```

### Assertions

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/assertions/{id}` | Get an assertion by ID |
| `POST` | `/assertions` | Award a badge to a person |

**Post assertion — request body:**
```json
{
  "ID": "assertion-unique-id",
  "BadgeID": "badge-unique-id",
  "PersonID": 42,
  "Salt": "random-salt",
  "IssuedFor": "Merged PR #123"
}
```

> [!NOTE]
> Creating an assertion fires an internal event that is evaluated against configured rules (e.g. `MilestoneRule`). If a rule is satisfied, an additional badge is automatically awarded.

## Rule Engine

The rule engine runs as a background goroutine and processes events published by the API handlers via an in-process channel.

### MilestoneRule

Configured in `main.go`, the `MilestoneRule` watches for `AssertionCreated` events and automatically grants a badge when a person's total assertion count reaches a configurable threshold.

```go
&rules.MilestoneRule{Threshold: 3, DB: conn}
```

To add your own rule, implement the `rules.Rule` interface:

```go
type Rule interface {
    Evaluate(event Event) bool
}
```

Then register it in the worker's `Rules` slice inside `cmd/server/main.go`.

## Development

| Command | Description |
|---------|-------------|
| `make server` | Start with live reload (Air) |
| `make build` | Build the binary |
| `make format` | Format all Go files with `gofmt` |
| `make lint` | Run `golangci-lint` |
| `make checks` | Run format + lint |
| `make deps` | Tidy and download Go modules |

## Docker

**Build the image:**
```bash
make docker-build
```

**Run the container:**
```bash
make docker-run
```

This uses `--network host` and reads configuration from `.env`. The server is exposed on port `8080`.

## Contributing

1. Fork the repository and create a feature branch.
2. Run `make checks` before submitting a pull request to ensure your code passes formatting and linting.
3. Keep pull requests focused; one feature or fix per PR.

## License

This project is a Fedora Infrastructure initiative. See the original [Tahrir](https://github.com/fedora-infra/tahrir) repository for license and attribution details.
