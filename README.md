# Config Server Go

**A lightweight, Spring Cloud Config-compatible configuration server built in Go.**

[![Go Version](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![CI](https://github.com/roniel-rhack/config-server-go/actions/workflows/ci.yml/badge.svg?branch=master)](https://github.com/roniel-rhack/config-server-go/actions/workflows/ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

---

## Overview

Config Server Go is a centralized configuration management microservice that serves versioned YAML configuration files over HTTP. It automatically flattens YAML into dot-notation properties, making it a drop-in backend for applications that consume [Spring Cloud Config](https://spring.io/projects/spring-cloud-config)-style configuration. Manage multiple configuration versions dynamically — list, switch, add, and delete — without restarting the server.

## Features

- **Spring Cloud Config-compatible responses** — returns `propertySources` JSON that Spring Boot clients understand natively
- **YAML-to-properties flattening** — automatically converts nested YAML into dot-notation key-value pairs
- **Multi-version configuration** — maintain multiple config versions side-by-side and switch between them at runtime
- **Raw file serving** — download any config file directly from the active version
- **Hot-reload** — configuration changes are detected via `fsnotify` and applied without restart
- **Graceful shutdown** — handles `SIGINT`/`SIGTERM` cleanly, draining in-flight requests
- **Path traversal protection** — sanitizes file paths to prevent directory escape attacks
- **Docker-ready** — multi-stage build produces a minimal Alpine image with health checks

---

## Quick Start

### Option A: Run with Go

```bash
cd src
go run config_server.go
```

The server starts on [http://localhost:8888](http://localhost:8888) by default.

### Option B: Run with Docker

```bash
docker compose up --build
```

Verify it is running:

```bash
curl http://localhost:8888/health
# ok
```

---

## API Reference

### Health Check

```
GET /health
```

Returns `ok` with a `200` status. Used by Docker `HEALTHCHECK` and load balancers.

---

<details>
<summary><strong>GET</strong> <code>/v1/versions</code> — List all configuration versions</summary>

Returns the current active version and all available versions, including the files present in each version directory.

**Response** `200 OK`

```json
{
  "current": {
    "version": "v1.0.0",
    "folder": "/opt/packages/config-server/configs/v1.0.0/",
    "files": ["application-dev.yaml", "application-prod.yaml"]
  },
  "available": [
    {
      "version": "v1.0.0",
      "folder": "/opt/packages/config-server/configs/v1.0.0/",
      "files": ["application-dev.yaml", "application-prod.yaml"]
    },
    {
      "version": "v2.0.0",
      "folder": "/opt/packages/config-server/configs/v2.0.0/",
      "files": ["application-dev.yaml"]
    }
  ]
}
```

</details>

<details>
<summary><strong>PUT</strong> <code>/v1/versions/set</code> — Switch the active configuration version</summary>

**Request body**

```json
{
  "version": "v2.0.0"
}
```

**Response** `200 OK`

```json
{
  "success": "Version set"
}
```

**Response** `200 OK` (already active)

```json
{
  "success": "Version already set"
}
```

**Error responses**

| Status | Body | Condition |
|--------|------|-----------|
| `400` | `{"error": "Invalid request"}` | Malformed or unparseable JSON body |
| `400` | `{"error": "Invalid version"}` | `version` field is empty |
| `400` | `{"error": "Version not available"}` | Version does not exist in `AVAILABLE_VERSIONS` |
| `500` | `{"error": "Error saving config"}` | Failed to persist the change to disk |

</details>

<details>
<summary><strong>POST</strong> <code>/v1/versions/add</code> — Create a new configuration version</summary>

Creates a new version directory and registers it in the available versions list. Spaces in the version name are replaced with underscores.

**Request body**

```json
{
  "version": "v3.0.0"
}
```

**Response** `200 OK`

```json
{
  "success": "Version added"
}
```

**Error responses**

| Status | Body | Condition |
|--------|------|-----------|
| `400` | `{"error": "Invalid request"}` | Malformed or unparseable JSON body |
| `400` | `{"error": "Invalid version"}` | `version` field is empty |
| `400` | `{"error": "Version already available"}` | Version already exists |
| `500` | `{"error": "Error creating folder"}` | Failed to create version directory |
| `500` | `{"error": "Error saving config"}` | Failed to persist the change to disk |

</details>

<details>
<summary><strong>DELETE</strong> <code>/v1/versions/:version/delete</code> — Remove a configuration version</summary>

Deletes the version directory and removes it from the available versions list. The currently active version cannot be deleted.

**Example**

```
DELETE /v1/versions/v2.0.0/delete
```

**Response** `200 OK`

```json
{
  "success": "Version deleted"
}
```

**Error responses**

| Status | Body | Condition |
|--------|------|-----------|
| `400` | `{"error": "Invalid version"}` | Version parameter is empty |
| `400` | `{"error": "Cannot delete current version"}` | Attempted to delete the active version |
| `400` | `{"error": "Version not available"}` | Version does not exist in `AVAILABLE_VERSIONS` |
| `500` | `{"error": "Error deleting folder"}` | Failed to remove version directory |
| `500` | `{"error": "Error saving config"}` | Failed to persist the change to disk |

</details>

<details>
<summary><strong>GET</strong> <code>/:filename</code> — Download a raw configuration file</summary>

Serves the raw contents of a file from the currently active version directory. Useful for downloading YAML files directly.

**Example**

```
GET /application-dev.yaml
```

**Response** `200 OK` — raw file contents with appropriate content type.

**Error responses**

| Status | Body | Condition |
|--------|------|-----------|
| `400` | `{"error": "Invalid filename"}` | Path traversal attempt detected |
| `404` | `{"error": "File not found"}` | File does not exist in the active version |
| `500` | `{"error": "Error reading file"}` | I/O error while reading the file |

</details>

<details>
<summary><strong>GET</strong> <code>/:appName/:profile</code> — Get config as flattened properties (Spring Cloud Config format)</summary>

Returns configuration from a YAML file named `{appName}-{profile}.yaml` in the active version, flattened into dot-notation key-value pairs. The response format is compatible with Spring Cloud Config clients.

**Example**

```
GET /application/dev
```

Reads from: `<CONFIG_FOLDER>/<CURRENT_VERSION>/application-dev.yaml`

**Response** `200 OK`

```json
{
  "name": "application",
  "propertySources": [
    {
      "name": "application",
      "source": {
        "server.port": "8080",
        "spring.datasource.url": "jdbc:mysql://localhost:3306/mydb",
        "spring.datasource.username": "root",
        "logging.level.root": "INFO"
      }
    }
  ]
}
```

</details>

---

## Configuration

The server reads its configuration from `config.yaml`, searched in the following order:

1. `/opt/packages/config-server/`
2. Current working directory

If no file is found, defaults are used and a new `config.yaml` is written automatically.

### Configuration Keys

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| `SERVER.PORT` | `int` | `8888` | Port the HTTP server listens on |
| `CONFIG_FOLDER` | `string` | `/opt/packages/config-server/configs/` | Root directory containing version subdirectories |
| `CONFIG_FILE` | `string` | `/opt/packages/config-server/config.yaml` | Path where the config file is written/persisted |
| `CURRENT_VERSION` | `string` | _(empty)_ | The currently active configuration version |
| `AVAILABLE_VERSIONS` | `[]string` | `[]` | List of registered version names |

### Sample `config.yaml`

```yaml
SERVER:
  PORT: 8888

CONFIG_FOLDER: /opt/packages/config-server/configs/
CONFIG_FILE: /opt/packages/config-server/config.yaml
CURRENT_VERSION: v1.0.0
AVAILABLE_VERSIONS:
  - v1.0.0
  - v2.0.0
```

---

## Project Structure

```
src/
├── config_server.go        # Entry point, route definitions, graceful shutdown
├── config/
│   └── config.go           # Viper setup, defaults, file watching (fsnotify)
├── custom_logguer/         # Styled terminal logging with charmbracelet/log
├── models/
│   ├── available_versions.go  # AvailableVersions and Version structs
│   ├── property_source.go     # Config and PropertySources response models
│   ├── set_version.go         # SetVersion request model
│   ├── web_200.go             # WebSuccess response model
│   └── web_error.go           # WebError response model
├── pkg/                    # YAML-to-properties parsing and transformation
├── services/
│   ├── properties.go       # GetConfigFile, GetConfig handlers
│   └── versions.go         # GetVersions, SetVersion, AddVersion, DeleteVersion
├── utils/                  # Path resolution and version helpers
└── versions/               # Version loader and file enumeration
```

---

## Docker

### Multi-Stage Build

The `Dockerfile` uses a two-stage build:

1. **Builder stage** (`golang:1.25`) — downloads dependencies, compiles a statically-linked binary with `CGO_ENABLED=0`
2. **Production stage** (`alpine:3`) — copies only the binary, runs as non-root `appuser` (UID 1000)

The final image is minimal and contains no Go toolchain.

### Health Check

The container includes a built-in health check that polls `/health` every 30 seconds:

```dockerfile
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
  CMD wget -qO- http://localhost:8888/health || exit 1
```

### Running with Docker Compose

```bash
docker compose up --build
```

The default `docker-compose.yml` maps port `8888:8888`.

### Mounting a Configuration Volume

To provide your own config and version directories:

```bash
docker run -d \
  -p 8888:8888 \
  -v ./config.yaml:/opt/packages/config-server/config.yaml \
  -v ./configs:/opt/packages/config-server/configs \
  config-server-go
```

---

## Development

All commands are run from the `src/` directory.

```bash
# Install dependencies
go mod download

# Run the server
go run config_server.go

# Build the binary
go build -v ./...

# Run static analysis
go vet ./...

# Run tests with race detection
go test -race -coverprofile=coverage.out ./...

# View coverage report
go tool cover -func=coverage.out
```

---

## CI/CD

The GitHub Actions workflow (`.github/workflows/ci.yml`) runs on pushes and pull requests to `master`:

**Build job** — runs across Go 1.25 and 1.26:
- Verifies module dependencies (`go mod verify`)
- Runs static analysis (`go vet ./...`)
- Builds all packages (`go build -v ./...`)
- Runs tests with race detection and coverage profiling
- Enforces a minimum **85% code coverage** threshold

**Docker job** — runs after the build job passes:
- Builds the Docker image to verify the `Dockerfile` is valid

---

## Tech Stack

| Component | Library |
|-----------|---------|
| HTTP framework | [Fiber v2](https://gofiber.io/) |
| Configuration | [Viper](https://github.com/spf13/viper) |
| YAML parsing | [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) |
| YAML highlighting | [goccy/go-yaml](https://github.com/goccy/go-yaml) |
| File watching | [fsnotify](https://github.com/fsnotify/fsnotify) |
| Logging | [charmbracelet/log](https://github.com/charmbracelet/log) |
| Containerization | [Docker](https://www.docker.com/) with multi-stage build |

---

## License

This project is licensed under the [MIT License](LICENSE).
