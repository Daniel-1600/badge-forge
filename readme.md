# badge-forge

This project is inspired [Tahrir](https://github.com/fedora-infra/tahrir) and [Tahrir-api](https://github.com/fedora-infra/tahrir-api) where I have been a contributor the open badge assertion server originally written in Python. This project exposes a JSON REST API for managing persons, badges, and badge assertions, backed by a PostgreSQL database.

## Features

- **REST API** — Manage persons, badges, and badge assertions over HTTP.
- **Event-driven rule engine** — Automatically award milestone badges when a person reaches a configurable assertion threshold, and badges for merged GitHub pull requests.
- **GORM + PostgreSQL** — Type-safe ORM-backed persistence with a well-defined model layer.
- **Live reload** — Hot-rebuild during development via [Air](https://github.com/air-verse/air).
- **Docker support** — Multi-stage Dockerfile for lean production images.

## Project Structure

```
tahrir-go/
├── cmd/
│   ├── server/
│   │   └── main.go          # Application entrypoint
│   └── award-prs/
│       └── main.go          # Awards merged_prs badges from stored GitHub PRs
├── internal/
│   ├── db/                  # Database connection and query helpers
│   ├── handlers/            # HTTP handler functions
│   ├── models/              # GORM model definitions (Person, Badge, Assertion, GitHub*)
│   └── rules/               # Rule engine interface and implementations
│       ├── rules.go         # Event types and Rule interface
│       ├── milestone_rule.go # MilestoneRule: awards a badge after N assertions
│       ├── merged_prs_rule.go # MergedPRsRule: awards a badge after N merged PRs
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
| `POST` | `/persons` | Find or create a person by nickname (`201` if created, `200` if they already exist) |

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

Configured in `config/rules.yaml`, the `MilestoneRule` watches for `AssertionCreated` events and awards its `badge_id` once a person holds `threshold` assertions of other badges. Each person earns a milestone badge at most once, and the badge must already exist.

```yaml
rules:
  - id: first-milestone
    type: milestone
    enabled: true
    threshold: 3
    badge_id: badge-first-steps
```

### MergedPRsRule

The `merged_prs` rule awards its `badge_id` once a person has authored `threshold` merged pull requests in the `github_*` tables. `repository` limits it to one repository, `labels` only counts PRs carrying at least one of the listed labels, and `exclude_labels` skips PRs carrying any of them. Each person earns the badge at most once.

```yaml
rules:
  - id: compiler-contributor
    name: "Compiler Contributor"
    description: "Had three compiler pull requests merged into rust-lang/rust."
    type: merged_prs
    enabled: true
    repository: rust-lang/rust
    labels: [T-compiler]
    exclude_labels: [rollup]
    threshold: 3
    badge_id: rust-compiler-contributor
```

These rules don't run in the server. `make award-prs` runs them: it creates any missing `merged_prs` badges from the rule's name and description, gives each human PR author a person (matching by nickname, or a new one with their GitHub noreply email), and awards the badges. Re-running it is safe.

```bash
make db-seed-rust   # load rust-lang/team and recent rust-lang/rust PRs
make award-prs      # award the merged_prs badges
```

To add your own rule, implement the `rules.Rule` interface:

```go
type Rule interface {
    Evaluate(event Event) bool
    AwardBadgeID() string
}
```

Then add a case for it in `rules.LoadRules` (`internal/rules/config.go`).

## Trying it on real data

`scripts/rust_replay.sh` replays recently merged `rust-lang/rust` pull requests into a running server, so you can watch real contributors earn badges and the milestone rule fire. It needs `gh` (logged in), `jq`, `curl` and `openssl`.

```bash
make server                      # in one terminal
./scripts/rust_replay.sh         # in another
API=http://localhost:9090 PAGES=5 REPO=rust-lang/cargo ./scripts/rust_replay.sh
```

Bots and rollup PRs are skipped, and assertion IDs come from PR numbers, so re-running is safe (repeats show up as `500`). Keep the results local: the contributors haven't opted in.

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

