# CLAUDE.md — Sound of Gothic Go Backend

This file gives Claude Code the context and conventions it needs to work effectively in this codebase.

---

## Project at a Glance

**Language:** Go 1.24.3
**Module:** `soundofgothic.pl/backend`
**Purpose:** Dual-server backend for the Sound of Gothic Polish voice-mod platform:
- `sog-backend` — read-only REST API (port 3001) consumed by the web frontend
- `sog-mod-server` — MCP server (stdio or HTTP/SSE) consumed by AI agents to manage mod production

The project stores the complete Gothic I & II dialogue database (recordings, NPCs, guilds, voice actors) and a mod-management layer for TTS-generated voice replacements.

---

## Essential Commands

### Development

```bash
# Start the pre-seeded PostgreSQL database (required for the servers and integration tests)
make dev

# Run the REST API server
go run ./cmd/sog-backend/ -c config.yml

# Run the MCP server (stdio mode — for claude-code / Claude Desktop)
go run ./cmd/sog-mod-server/ -c config.yml

# Run all tests
go test ./...

# Run only unit tests (no DB required)
go test ./internal/mcp/...

# Run integration tests (requires `make dev` first)
go test ./internal/postgres/...

# Lint (must pass before any PR)
golangci-lint run
```

### Docker / Build

```bash
make docker        # Build production image tagged with git SHA
make docker-push   # Build + push (requires clean working tree)
make version       # Print current image tag
```

---

## Architecture

```
domain/          ← Pure Go structs + service interfaces. No DB, no HTTP.
postgres/        ← Implements domain interfaces using Bun ORM + PostgreSQL.
rest/            ← HTTP handlers (chi router). Reads only; uses domain interfaces.
mcp/             ← MCP tools and resources. Reads and writes via domain interfaces.
config/          ← Viper YAML config loader.
cmd/             ← Two binary entry points (sog-backend, sog-mod-server).
```

The dependency flow is strictly: `cmd → rest/mcp → domain ← postgres`.
Neither `rest/` nor `mcp/` imports `postgres/` directly — they receive `domain.Repositories` via constructor injection.

---

## Code Conventions

### Adding a New REST Endpoint

1. Define an input struct with `httpin` and `validate` tags in the relevant handler file:
   ```go
   type MyInput struct {
       Page int64 `in:"query=page;default=1" validate:"min=1"`
       Name string `in:"query=name"`
   }
   ```
2. Register the route in `internal/rest/rest.go` using `middlewares.ValidatedInput`:
   ```go
   r.With(middlewares.ValidatedInput(MyInput{})).Get("/my-endpoint", env.myHandler)
   ```
3. Retrieve the parsed input from context:
   ```go
   input := r.Context().Value(httpin.Input).(*MyInput)
   ```
4. Use `rjson.RespondWithJSON` / `rjson.NewPagedResponse` / `rjson.InternalError` for responses.

### Adding a New MCP Tool

1. Create a new file in `internal/mcp/tools/`.
2. Implement the `mcp.Tool` interface (defined in `internal/mcp/server.go`):
   ```go
   type Tool interface {
       Handler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
       Tool() mcp.Tool
   }
   ```
3. Use `tools.MapTo[T](request.Params.Arguments)` to parse tool call arguments into a typed struct.
4. Register the tool in `NewGothicToolsPack()` (`internal/mcp/tools/tool.go`).

### Adding a New Domain Entity

1. Create `internal/domain/<entity>.go` — struct + service interface (no DB tags beyond Bun's `BaseModel`).
2. Create `internal/postgres/<entity>.go` — repository implementation using `commonRepository[T]`.
3. Add the service accessor to `domain.Repositories` in `internal/domain/domain.go`.
4. Implement the accessor in `internal/postgres/postgres.go`.
5. If the entity needs a new DB table, add a migration in `internal/postgres/migrations/`.

### Query Modifiers

Complex queries are built using composable `QueryModifier` functions in `internal/postgres/mods/base.go`:

```go
// Good — compose modifiers
modifiers := []mods.QueryModifier{
    mods.WithPaging(page, pageSize),
    mods.WithIn("r.game_id", gameIDs),
    mods.WithFullTextSearch("r.transcript_tsv", query),
    mods.WithRelations("Game", "NPC", "Voice"),
}
```

Never build raw SQL strings in repository files; use the provided modifier functions or add new ones to `mods/base.go`.

---

## Key Files Quick Reference

| File | What It Does |
|------|-------------|
| `internal/domain/domain.go` | `Repositories` interface — the central seam |
| `internal/domain/alternative.go` | `AlternativeState` enum and `Alternative` struct |
| `internal/rest/rest.go` | Route registration, middleware wiring |
| `internal/mcp/server.go` | `MCPServer` constructor, `Tool`/`Resource` interfaces |
| `internal/mcp/tools/tool.go` | `MapTo[T]()` helper, `NewGothicToolsPack()` |
| `internal/postgres/postgres.go` | DB connection, `NewPostgresRepositories` |
| `internal/postgres/mods/base.go` | All composable `QueryModifier` functions |
| `internal/postgres/migrations/` | Embedded SQL migrations |
| `config.yml` | Local dev config (address, db, mcp) |
| `docker/config-template.yaml` | Production config template (envsubst) |

---

## Configuration

Config is loaded from YAML by `internal/config/config.go` via Viper. Search order:
1. `-config <path>` CLI flag
2. `./config.yml`
3. `/app/config.yml` (Docker)
4. `../../config.yml` (test working directories)

**Local dev defaults (`config.yml`):**
```yaml
address: localhost:3001
db:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  database: postgres
mcp:
  address: localhost:8080   # or "stdio"
```

**Docker env vars:** `ADDRESS`, `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`

---

## Database

- **Engine:** PostgreSQL (pre-seeded image `soundofgothic/db:v2.1.2`)
- **ORM:** [Bun](https://bun.uptrace.dev/) with `pgdriver` (pure Go, no `lib/pq`)
- **Migrations:** Embedded SQL files in `internal/postgres/migrations/`, run automatically by `sog-mod-server` on startup
- **Full-text search:** Polish Ispell dictionary (`pl_ispell`), used for `filter` on recordings

### Development Database

```bash
make dev
# Starts postgres on localhost:5432
# Credentials: postgres/postgres, database: postgres
# Pre-loaded with Gothic I + II dialogue data
```

---

## Testing

- **Unit tests** — `internal/mcp/tools/*_test.go`, `internal/mcp/resources/*_test.go`
  - No DB required. Mock `domain.Repositories` with table-driven tests.
- **Integration tests** — `internal/postgres/*_test.go`
  - Require a live PostgreSQL connection. See `internal/test/db_connection.go`.
  - Run after `make dev`.

**Pattern for mocking repositories in unit tests:**
```go
// Implement only the methods you need; embed a nil for the rest.
type mockRepositories struct {
    domain.Repositories
    mod mockModService
}
func (m *mockRepositories) Mod() domain.ModService { return &m.mod }
```

---

## Linting

CI runs `golangci-lint v2.1`. Run locally before pushing:

```bash
golangci-lint run
```

Common issues to avoid:
- Unused imports / variables
- `err` not checked
- Exported identifiers without doc comments (if adding to `domain/`)

---

## MCP Transport Modes

| Mode | Config | Use Case |
|------|--------|----------|
| `stdio` | `mcp.address: stdio` | `claude-code`, Claude Desktop, any MCP CLI client |
| HTTP/SSE | `mcp.address: localhost:8080` | Remote or browser-accessible agents; endpoint at `/mcp` |

The MCP server logs every tool call via `slog.Info` with the tool name and arguments.

---

## Pitfalls & Gotchas

- **REST server does NOT run migrations.** Only `sog-mod-server` calls `migrations.Migrate()` on startup. If you add a migration and only start `sog-backend`, the new tables will not exist.
- **`Alternative` has a composite PK.** `(mod_id, recording_id)` — both fields are `pk` in Bun tags. Upsert uses `ON CONFLICT (mod_id, recording_id) DO UPDATE`.
- **`gameID` in `RecordingGetInput` is validated as `max=2`** — only Gothic I (1) and Gothic II (2) exist.
- **`filter` on recordings uses full-text search.** Very short queries (under ~3 chars) may return no results from `plainto_tsquery`. Prefer LIKE-based `filter` for short strings (available via `WithTextSearch`).
- **Array query parameters** accept both `voiceID=1&voiceID=2` and `voiceID[]=1&voiceID[]=2` syntaxes (handled by `httpin`).
- **CORS is fully open.** `middlewares.CORS` allows all origins — intentional for the public read-only API.
