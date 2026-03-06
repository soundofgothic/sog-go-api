# Sound of Gothic — Go Backend API

A dual-server Go backend for the **Sound of Gothic** project — a Polish fan community voice-mod platform for the Gothic game series. The project exposes a read-optimised **REST API** for web clients and a **Model Context Protocol (MCP) server** enabling AI agents to manage voice mod production workflows.

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Architecture](#architecture)
3. [Repository Layout](#repository-layout)
4. [Domain Model](#domain-model)
5. [REST API Reference](#rest-api-reference)
6. [MCP Server Reference](#mcp-server-reference)
7. [Configuration](#configuration)
8. [Getting Started (Local Development)](#getting-started-local-development)
9. [Building & Deployment](#building--deployment)
10. [Testing](#testing)
11. [CI/CD](#cicd)
12. [Dependencies](#dependencies)
13. [Contributing](#contributing)

---

## Project Overview

Sound of Gothic is a community project that creates Polish voice-over mods for the Gothic I and Gothic II video games. This backend service:

- Stores and serves the complete **dialogue database** — every NPC line, its source file, the assigned voice actor, and the guild/faction.
- Manages **voice mods**: packaged sets of alternative TTS-generated recordings that replace original dialogue.
- Provides an **MCP interface** so that AI agents can query game data and drive the entire mod-creation pipeline (listing recordings, creating mods, scheduling TTS generation, tracking synthesis state).

---

## Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        Clients                                  │
│        Web Frontend            AI Agent / Claude                │
└──────────┬─────────────────────────────┬───────────────────────┘
           │ HTTP                        │ MCP (stdio / SSE)
           ▼                             ▼
┌──────────────────┐          ┌──────────────────────┐
│  sog-backend     │          │  sog-mod-server       │
│  (REST, :3001)   │          │  (MCP, stdio/:8080)   │
└──────────┬───────┘          └──────────┬────────────┘
           │                             │
           └──────────┬──────────────────┘
                      │ domain.Repositories
                      ▼
           ┌──────────────────┐
           │   PostgreSQL      │
           │  (soundofgothic   │
           │   /db:v2.1.2)    │
           └──────────────────┘
```

### Two Binary Architecture

| Binary | Package | Purpose |
|--------|---------|---------|
| `sog-backend` | `cmd/sog-backend` | Read-only REST API for the web frontend |
| `sog-mod-server` | `cmd/sog-mod-server` | MCP server for AI-driven mod production; also runs DB migrations |

Both binaries share the same `internal/` packages — they are two different process entry-points over a common library.

### Key Design Decisions

- **Repository pattern** — `domain.Repositories` interface abstracts all data access; swapping the database layer requires no change to handlers.
- **Functional query modifiers** — `internal/postgres/mods` provides composable `QueryModifier` functions that build Bun `SelectQuery` objects, avoiding query string concatenation.
- **Generics** — `commonRepository[T]` provides a generic base for all repository implementations.
- **Clean Architecture layers** — `domain/` (entities + interfaces) → `postgres/` (implementations) → `rest/` or `mcp/` (delivery).

---

## Repository Layout

```
sog-go-api/
│
├── cmd/
│   ├── sog-backend/        # REST API entry point
│   │   └── main.go
│   └── sog-mod-server/     # MCP server entry point
│       └── main.go
│
├── internal/
│   ├── config/             # Viper-based YAML config loader
│   │   └── config.go
│   │
│   ├── domain/             # Pure domain: entities + service interfaces
│   │   ├── domain.go       # Repositories aggregate interface
│   │   ├── game.go
│   │   ├── guild.go
│   │   ├── npc.go
│   │   ├── voice.go
│   │   ├── recording.go
│   │   ├── source_file.go
│   │   ├── mod.go
│   │   ├── alternative.go
│   │   └── tts_voice.go
│   │
│   ├── postgres/           # Repository implementations (Bun ORM)
│   │   ├── postgres.go     # Connection setup, NewPostgresRepositories
│   │   ├── common.go       # Generic commonRepository[T]
│   │   ├── game.go
│   │   ├── guild.go
│   │   ├── npc.go
│   │   ├── voice.go
│   │   ├── recording.go
│   │   ├── source_file.go
│   │   ├── mod.go
│   │   ├── alternative.go
│   │   ├── tts_voice.go
│   │   ├── mods/           # Composable query modifiers
│   │   │   ├── base.go
│   │   │   ├── search_struct.go
│   │   │   └── recordings_count.go
│   │   └── migrations/     # SQL migrations (embedded)
│   │       ├── migrations.go
│   │       └── 20250516_add_mods_equipment.up.sql
│   │
│   ├── rest/               # HTTP delivery layer
│   │   ├── rest.go         # Router setup, middleware registration
│   │   ├── games.go
│   │   ├── guilds.go
│   │   ├── npcs.go
│   │   ├── voices.go
│   │   ├── recordings.go
│   │   ├── source_files.go
│   │   ├── middlewares/
│   │   │   ├── cors.go            # Permissive CORS middleware
│   │   │   └── validated_input.go # httpin + validator middleware
│   │   └── rjson/
│   │       ├── json.go            # JSON response helpers
│   │       └── errors.go          # Standardised error responses
│   │
│   ├── mcp/                # MCP delivery layer
│   │   ├── server.go       # MCPServer constructor, Tool/Resource registration
│   │   ├── resources/
│   │   │   ├── actors.go          # actors:/// resource
│   │   │   └── actors_test.go
│   │   └── tools/
│   │       ├── tool.go            # MapTo() helper, shared types
│   │       ├── mod_list.go
│   │       ├── mod_create.go
│   │       ├── mod_delete.go
│   │       ├── alternatives_list.go
│   │       ├── alternatives_persist.go
│   │       ├── guild_list.go
│   │       ├── npc_list.go
│   │       ├── voice_list.go
│   │       ├── tts_voice_list.go
│   │       ├── recordings_list.go
│   │       ├── scripts_list.go
│   │       └── *_test.go
│   │
│   └── test/
│       └── db_connection.go  # Test helpers for DB integration tests
│
├── development/
│   └── docker-compose.yml    # Local dev: pre-seeded PostgreSQL
│
├── docker/
│   ├── Dockerfile            # Multi-stage Alpine production build
│   ├── config-template.yaml  # envsubst-based config template
│   └── entrypoint.sh         # Docker entrypoint
│
├── .github/
│   └── workflows/
│       ├── tests.yml         # golangci-lint on every push
│       └── release.yml       # Docker image publish on release
│
├── config.yml                # Local development config
├── go.mod
├── go.sum
└── Makefile
```

---

## Domain Model

### Entity Overview

```
Game ─────┬──── Guild ──── NPC ──┐
          │                      │
          ├──── SourceFile        ├──── Recording ──── Alternative
          │                      │         │               │
          └─────────────────────-┘      Voice          TTSVoice
                                           │               │
                                           └─── Mod ───────┘
```

### Entities

#### `Game`
Represents a Gothic game title (Gothic I, Gothic II, etc.).

| Field | Type | Description |
|-------|------|-------------|
| `id` | int64 | Primary key |
| `name` | string | Display name of the game |

#### `Guild`
An in-game faction or guild (e.g. Ore Barons, Old Camp).

| Field | Type | Description |
|-------|------|-------------|
| `id` | int64 | Primary key |
| `name` | string | Display name |
| `inGameID` | string | Identifier used inside the game scripts |
| `gameID` | int64 | FK → Game |

#### `NPC`
A Non-Player Character with a dialogue voice.

| Field | Type | Description |
|-------|------|-------------|
| `id` | int64 | Primary key |
| `name` | string | NPC name |
| `inGameID` | string | Identifier in game scripts |
| `voiceID` | int64 | FK → Voice (assigned voice actor) |
| `guildID` | int64 | FK → Guild |

#### `Voice`
A Polish voice actor who recorded dialogue.

| Field | Type | Description |
|-------|------|-------------|
| `id` | int64 | Primary key |
| `name` | string | Actor's name |

#### `SourceFile`
A game script file that is the source of dialogue.

| Field | Type | Description |
|-------|------|-------------|
| `id` | int64 | Primary key |
| `name` | string | File name |
| `type` | string | File type (e.g. `ou`, `ou_addon`) |
| `gameID` | int64 | FK → Game |

#### `Recording`
A single dialogue line — the core entity of the system.

| Field | Type | Description |
|-------|------|-------------|
| `id` | int64 | Primary key |
| `wave` | string | WAV filename (unique per game) |
| `transcript` | string | Polish dialogue text |
| `gameID` | int64 | FK → Game |
| `sourceFileID` | int64 | FK → SourceFile |
| `npcID` | *int64 | FK → NPC (nullable) |
| `guildID` | *int64 | FK → Guild (nullable) |
| `voiceID` | *int64 | FK → Voice (nullable) |
| `title` | *string | Optional contextual title |

#### `Mod`
A voice-over mod package — a named collection of alternative recordings.

| Field | Type | Description |
|-------|------|-------------|
| `id` | int64 | Primary key (auto-increment) |
| `gameID` | int64 | FK → Game |
| `name` | string | Mod display name |
| `description` | string | Mod description |
| `version` | string | Semantic version string |

#### `Alternative`
A TTS-synthesised replacement recording for a specific mod+recording combination. Composite primary key: `(mod_id, recording_id)`.

| Field | Type | Description |
|-------|------|-------------|
| `modID` | int64 | PK / FK → Mod |
| `recordingID` | int64 | PK / FK → Recording |
| `ttsVoiceID` | int64 | FK → TTSVoice |
| `transcript` | string | Custom transcript for TTS |
| `state` | AlternativeState | Synthesis state (see below) |
| `wave` | string | Output WAV filename |

**AlternativeState values:**

| State | Meaning |
|-------|---------|
| `none` | Not yet scheduled |
| `scheduled` | Queued for TTS synthesis |
| `in_progress` | Currently being synthesised |
| `completed` | Synthesis finished successfully |
| `failed` | Synthesis failed |

#### `TTSVoice`
A Text-to-Speech voice provider/model entry.

| Field | Type | Description |
|-------|------|-------------|
| `id` | int64 | Primary key |
| `name` | string | Provider/voice identifier |

---

## REST API Reference

**Base URL (local):** `http://localhost:3001`

All endpoints are `GET`. The API is read-only (no write endpoints are exposed via REST — write operations go through MCP). Middleware applied to all routes:

- `middleware.Logger` — structured request logging
- `middleware.Recoverer` — panic recovery
- `middlewares.CORS` — allows all origins, `GET/POST/PUT/DELETE/OPTIONS`, standard headers

### Pagination

Paginated endpoints return:

```json
{
  "total": 1234,
  "page": 1,
  "pageSize": 50,
  "results": [...]
}
```

Non-paginated endpoints return a plain JSON array `[...]`.

### Error Response

```json
{ "internal_error": "human-readable message" }
```

Validation errors return HTTP 400 with field-level detail.

---

### `GET /games`

Returns all games. No parameters.

**Response:**
```json
[
  { "id": 1, "name": "Gothic" },
  { "id": 2, "name": "Gothic II" }
]
```

---

### `GET /guilds`

Paginated list of guilds, optionally filtered.

| Parameter | Type | Default | Constraints | Description |
|-----------|------|---------|-------------|-------------|
| `page` | int | 1 | min=1 | Page number |
| `pageSize` | int | 50 | min=10, max=100 | Results per page |
| `filter` | string | — | — | Text search on guild name |
| `gameID[]` | int[] | — | — | Filter by game ID(s) |
| `inGameAlias` | string | — | — | Filter by in-game identifier |
| `voiceID[]` | int[] | — | — | Filter by voice actor ID(s) |
| `scriptID[]` | int[] | — | — | Filter by source file ID(s) |
| `npcID[]` | int[] | — | — | Filter by NPC ID(s) |
| `id[]` | int[] | — | — | Filter by guild ID(s) |

Results include a `recordingsCount` field.

---

### `GET /npcs`

Paginated list of NPCs.

| Parameter | Type | Default | Constraints | Description |
|-----------|------|---------|-------------|-------------|
| `page` | int | 1 | min=1 | Page number |
| `pageSize` | int | 50 | min=10, max=100 | Results per page |
| `filter` | string | — | — | Text search on NPC name |
| `gameID[]` | int[] | — | — | Filter by game ID(s) |
| `voiceID[]` | int[] | — | — | Filter by voice actor ID(s) |
| `guildID[]` | int[] | — | — | Filter by guild ID(s) |
| `scriptID[]` | int[] | — | — | Filter by source file ID(s) |
| `id[]` | int[] | — | — | Filter by NPC ID(s) |

---

### `GET /voices`

List of voice actors. Not paginated, returns all matching voices.

| Parameter | Type | Description |
|-----------|------|-------------|
| `filter` | string | Text search on voice actor name |
| `gameID[]` | int[] | Filter by game ID(s) |
| `npcID[]` | int[] | Filter by NPC ID(s) |
| `guildID[]` | int[] | Filter by guild ID(s) |
| `scriptID[]` | int[] | Filter by source file ID(s) |
| `id[]` | int[] | Filter by voice ID(s) |

Results include a `recordingsCount` field.

---

### `GET /source_files`

Paginated list of source script files.

| Parameter | Type | Default | Constraints | Description |
|-----------|------|---------|-------------|-------------|
| `page` | int | 1 | min=1 | Page number |
| `pageSize` | int | 50 | min=10, max=100 | Results per page |
| `filter` | string | — | — | Text search on filename |
| `type` | string | — | — | Exact match on file type |
| `gameID[]` | int[] | — | — | Filter by game ID(s) |
| `guildID[]` | int[] | — | — | Filter by guild ID(s) |
| `npcID[]` | int[] | — | — | Filter by NPC ID(s) |
| `voiceID[]` | int[] | — | — | Filter by voice actor ID(s) |
| `id[]` | int[] | — | — | Filter by source file ID(s) |

---

### `GET /recordings`

Paginated list of dialogue recordings. Results include eagerly-loaded relations (Game, SourceFile, NPC, Guild, Voice, Alternative).

| Parameter | Type | Default | Constraints | Description |
|-----------|------|---------|-------------|-------------|
| `page` | int | 1 | min=1 | Page number |
| `pageSize` | int | 50 | min=10, max=100 | Results per page |
| `filter` | string | — | — | Full-text search on transcript (Polish `pl_ispell`) |
| `sourceFileID[]` | int[] | — | — | Filter by source file ID(s) |
| `voiceID[]` | int[] | — | — | Filter by voice actor ID(s) |
| `gameID[]` | int[] | — | — | Filter by game ID(s) |
| `npcID[]` | int[] | — | — | Filter by NPC ID(s) |
| `guildID[]` | int[] | — | — | Filter by guild ID(s) |

**Note:** The `filter` parameter uses PostgreSQL full-text search with Polish Ispell stemming (`plainto_tsquery('pl_ispell', ?)`). Short queries (< 3 characters) fall back to `LOWER(column) LIKE` matching.

---

### `GET /recordings/{gameID}/{wave}`

Fetch a single recording by its game + WAV filename composite key.

| Parameter | Type | Constraints | Description |
|-----------|------|-------------|-------------|
| `gameID` | int | required, min=1, max=2 | Game ID (path parameter) |
| `wave` | string | required | WAV filename (path parameter) |

**Response:** Single `Recording` object with all relations populated, or `404`/`500` on error.

---

## MCP Server Reference

The MCP server implements the [Model Context Protocol](https://modelcontextprotocol.io/) and is designed to be consumed by AI agents (e.g. Claude via `claude-code` or any MCP-compatible client).

**Server name:** `Sound of Gothic`
**Version:** `0.1.0`
**Transport:** `stdio` (default) or HTTP/SSE (configured via `mcp.address`)

When using SSE transport, the base path is `/mcp`.

### Connecting via stdio (Claude Desktop / claude-code)

```json
{
  "mcpServers": {
    "soundofgothic": {
      "command": "/path/to/sog-mod-server",
      "args": ["-c", "/path/to/config.yml"]
    }
  }
}
```

### Connecting via SSE

```
GET http://localhost:8080/mcp/sse
```

---

### MCP Resources

Resources are read-only data sources that an AI agent can fetch.

#### `actors:///`

Returns a list of all voice actors (voices) with their recording counts.

**Response:** JSON array of `Voice` objects including `recordingsCount`.

---

### MCP Tools

Tools are callable functions that allow an AI agent to perform actions.

#### `list_mods`

List all mods in the database.

**Parameters:** None

**Returns:** JSON array of `Mod` objects.

---

#### `create_or_update_mod`

Create a new mod or update an existing one (upsert by ID).

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `name` | string | Yes | Mod display name |
| `description` | string | Yes | Mod description |
| `version` | string | Yes | Semantic version (e.g. `1.0.0`) |
| `id` | int | No | If provided, updates existing mod |
| `game_id` | int | Yes | ID of the target game |

---

#### `mod_delete`

Delete a mod by ID.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `id` | int | Yes | Mod ID to delete |

---

#### `guild_list`

List guilds, with optional filtering.

| Parameter | Type | Description |
|-----------|------|-------------|
| `game_id` | int | Filter by game |
| `voice_id` | int | Filter by voice actor |

---

#### `npc_list`

List NPCs, with optional filtering.

| Parameter | Type | Description |
|-----------|------|-------------|
| `game_id` | int | Filter by game |
| `guild_id` | int | Filter by guild |
| `voice_id` | int | Filter by voice actor |

---

#### `voice_list`

List voice actors.

| Parameter | Type | Description |
|-----------|------|-------------|
| `game_id` | int | Filter by game |

---

#### `tts_voice_list`

List available Text-to-Speech voices.

**Parameters:** None

---

#### `recordings_list`

List dialogue recordings with filtering.

| Parameter | Type | Description |
|-----------|------|-------------|
| `game_id` | int | Filter by game |
| `npc_id` | int | Filter by NPC |
| `guild_id` | int | Filter by guild |
| `voice_id` | int | Filter by voice actor |
| `source_file_id` | int | Filter by source file |
| `query` | string | Full-text search on transcript |
| `page` | int | Page number (default: 1) |
| `page_size` | int | Results per page (default: 50) |

---

#### `scripts_list`

List source script files.

| Parameter | Type | Description |
|-----------|------|-------------|
| `game_id` | int | Filter by game |
| `type` | string | Filter by file type |

---

#### `alternatives_list`

List alternative (TTS) recordings for a mod.

| Parameter | Type | Description |
|-----------|------|-------------|
| `mod_id` | int | Filter by mod |
| `state` | string | Filter by state (`none`, `scheduled`, `in_progress`, `completed`, `failed`) |
| `tts_voice_id` | int | Filter by TTS voice |

---

#### `alternatives_persist`

Create or update an alternative recording entry.

| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| `mod_id` | int | Yes | Mod ID |
| `recording_id` | int | Yes | Recording ID |
| `tts_voice_id` | int | Yes | TTS voice to use |
| `transcript` | string | Yes | Custom transcript for TTS input |
| `state` | string | No | Initial state (default: `none`) |
| `wave` | string | No | Output WAV filename |

---

## Configuration

Configuration is loaded from a YAML file. The loader searches the following paths in order:

1. Path specified via `-config <path>` / `-c <path>` CLI flag
2. `./config.yml` (current directory)
3. `/app/config.yml` (Docker container path)
4. `../../config.yml` (relative, for tests)

If no file is found, built-in defaults are used (except PostgreSQL credentials, which are required).

### Full `config.yml` Schema

```yaml
# Address the REST API server listens on.
# Default: ":8080"
address: "localhost:3001"

# PostgreSQL connection settings. All fields are required.
db:
  host: "localhost"       # required
  port: 5432              # required
  user: "postgres"        # required
  password: "postgres"    # required
  database: "postgres"    # required

# MCP server transport configuration.
# Use "stdio" for stdin/stdout (Claude Desktop, claude-code).
# Use "host:port" for HTTP/SSE transport.
# Default: "stdio"
mcp:
  address: "stdio"
```

### Docker Environment Variables

In the Docker production image the config is generated from a template (`docker/config-template.yaml`) via `envsubst`. The following environment variables must be set:

| Variable | Config Key | Example |
|----------|-----------|---------|
| `ADDRESS` | `address` | `:3001` |
| `DB_HOST` | `db.host` | `postgres` |
| `DB_PORT` | `db.port` | `5432` |
| `DB_USER` | `db.user` | `sog` |
| `DB_PASSWORD` | `db.password` | `secret` |
| `DB_NAME` | `db.database` | `soundofgothic` |

---

## Getting Started (Local Development)

### Prerequisites

- **Go 1.24.3+**
- **Docker** or **Podman** (for the database)
- **make**

### 1. Start the Database

The development database is a pre-seeded PostgreSQL image containing all Gothic game data:

```bash
make dev
```

This runs `docker compose up -d` inside `development/`, starting:
- PostgreSQL on `localhost:5432` (image: `soundofgothic/db:v2.1.2`)
  - User: `postgres`, Password: `postgres`, Database: `postgres`

### 2. Configure

The default `config.yml` in the root is already set up for local development:

```yaml
address: localhost:3001
db:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  database: postgres
mcp:
  address: localhost:8080
```

### 3. Run the REST API Server

```bash
go run ./cmd/sog-backend/
# or with explicit config:
go run ./cmd/sog-backend/ -c config.yml
```

The API is now available at `http://localhost:3001`.

### 4. Run the MCP Server

```bash
# stdio mode (for use with claude-code / Claude Desktop):
go run ./cmd/sog-mod-server/ -c config.yml

# HTTP/SSE mode:
# Set mcp.address to "localhost:8080" in config.yml, then:
go run ./cmd/sog-mod-server/ -c config.yml
```

**Note:** The MCP server runs database migrations automatically on startup. The REST server does not.

### 5. Verify

```bash
curl http://localhost:3001/games
# → [{"id":1,"name":"Gothic"},{"id":2,"name":"Gothic II"}]

curl "http://localhost:3001/recordings?gameID=1&pageSize=10"
# → {"total":...,"page":1,"pageSize":10,"results":[...]}
```

---

## Building & Deployment

### Makefile Targets

```bash
make list          # Show all available targets
make dev           # Start development database via Docker Compose
make docker        # Build production Docker image (tagged with git short SHA)
make docker-push   # Build + push image (requires clean git working tree)
make version       # Print the current image tag (soundofgothic/api:<sha>)
```

The Makefile auto-detects `podman` and uses it instead of `docker` if available.

### Docker Build

The production image uses a two-stage build:

```
Stage 1 (golang:1.24.3-alpine):
  - Downloads Go module dependencies
  - Builds both binaries: sog-backend, sog-mod-server

Stage 2 (alpine:3.21):
  - Copies binaries
  - Copies config template + entrypoint script
  - Installs gettext (for envsubst)
  - Entrypoint: /app/entrypoint.sh
```

```bash
# Build the image
make docker
# → soundofgothic/api:<git-sha>

# Push to registry (requires no uncommitted changes)
make docker-push
```

### Versioning

The image is tagged with the short git commit SHA (`git rev-parse --short HEAD`). Semantic version tags are managed via the GitHub release workflow (`.github/workflows/release.yml`).

---

## Testing

### Unit Tests

```bash
go test ./...
```

### Test Coverage by Package

| Package | Test File(s) | What's Tested |
|---------|-------------|---------------|
| `internal/mcp/tools` | `tool_test.go` | `MapTo()` struct mapping utility |
| `internal/mcp/tools` | `mod_create_test.go` | Mod creation tool handler |
| `internal/mcp/tools` | `mod_delete_test.go` | Mod deletion tool handler |
| `internal/mcp/tools` | `mod_list_test.go` | Mod listing tool handler |
| `internal/mcp/resources` | `actors_test.go` | Actor resource handler |
| `internal/postgres` | `recording_test.go` | Recording repository (integration) |
| `internal/postgres` | `alternative_test.go` | Alternative repository (integration) |

### Integration Tests

Database integration tests (`internal/postgres/*_test.go`) use `internal/test/db_connection.go` to establish a real PostgreSQL connection. They require the development database to be running:

```bash
make dev
go test ./internal/postgres/...
```

---

## CI/CD

### `.github/workflows/tests.yml`

Runs on every push to any branch:

- **Linter:** `golangci-lint v2.1` (Go 1.21)

### `.github/workflows/release.yml`

Triggered on GitHub releases:

- Builds the Docker image
- Pushes to the container registry as `soundofgothic/api:<version>`

---

## Dependencies

### Direct Dependencies

| Package | Version | Purpose |
|---------|---------|---------|
| `github.com/go-chi/chi/v5` | v5.0.8 | HTTP router |
| `github.com/uptrace/bun` | v1.2.11 | PostgreSQL ORM |
| `github.com/uptrace/bun/dialect/pgdialect` | v1.2.11 | Bun PostgreSQL dialect |
| `github.com/uptrace/bun/driver/pgdriver` | v1.2.11 | Pure-Go PostgreSQL driver |
| `github.com/mark3labs/mcp-go` | v0.28.0 | Model Context Protocol server |
| `github.com/spf13/viper` | v1.15.0 | YAML configuration |
| `github.com/ggicci/httpin` | v0.10.1 | HTTP query/path parameter binding |
| `github.com/go-playground/validator/v10` | v10.12.0 | Struct validation |
| `github.com/samber/lo` | v1.50.0 | Generic utility functions |
| `github.com/enhanced-tools/errors` | v0.0.0-… | Structured error handling with stack traces |
| `github.com/stretchr/testify` | v1.10.0 | Test assertions |

### Go Version

**Required: Go 1.24.3** (uses generics, `log/slog`, range-over-integer).

---

## Contributing

1. Fork the repository and create a feature branch from `main`.
2. Ensure `golangci-lint run` passes locally before opening a PR.
3. Add tests for any new functionality.
4. Keep domain models (`internal/domain/`) free of infrastructure dependencies.
5. New MCP tools should implement the `mcp.Tool` interface (see `internal/mcp/server.go`).
6. New REST endpoints should use `middlewares.ValidatedInput` for input parsing.
