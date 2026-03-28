# Config Server Go

A lightweight configuration management microservice built with Go and [Fiber](https://gofiber.io/). It serves as a centralized hub for managing versioned application configuration files (YAML/properties), allowing clients to retrieve, switch, add, and delete configuration versions dynamically.

## Features

- Serve configuration files (YAML) over HTTP
- Automatic YAML-to-properties conversion (dot-notation flattening)
- Multi-version configuration management (list, switch, add, delete)
- File watching with hot-reload via fsnotify
- Colorized YAML output in terminal
- Docker-ready with multi-stage build

## API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/v1/versions` | List all available versions and current active version |
| `PUT` | `/v1/versions/set` | Switch to a different configuration version |
| `POST` | `/v1/versions/add` | Create a new configuration version |
| `DELETE` | `/v1/versions/:version/delete` | Remove a configuration version |
| `GET` | `/:filename` | Download a raw config file from the current version |
| `GET` | `/:appName/:profile` | Get config as flattened properties JSON |

### Response Format

**Config response:**
```json
{
  "name": "application",
  "propertySources": [
    {
      "name": "application",
      "source": {
        "server.port": "8888",
        "database.host": "localhost"
      }
    }
  ]
}
```

## Getting Started

### Prerequisites

- Go 1.22+
- Docker (optional)

### Run locally

```bash
cd src
go run config_server.go
```

The server starts on port `8888` by default.

### Run with Docker

```bash
docker compose up --build
```

### Configuration

The server reads from `config.yaml` (looked up in `/opt/packages/config-server/` or current directory):

| Key | Default | Description |
|-----|---------|-------------|
| `SERVER.PORT` | `8888` | Server listen port |
| `CONFIG_FOLDER` | `/opt/packages/config-server/configs/` | Path to versioned config directories |
| `CURRENT_VERSION` | — | Active configuration version |
| `AVAILABLE_VERSIONS` | `[]` | List of available versions |

## Project Structure

```
src/
├── config_server.go          # Entry point, routes, middleware
├── config/config.go          # Viper configuration setup
├── custom_logguer/           # Styled terminal logging
├── models/                   # Request/response structs
├── pkg/                      # YAML/properties parsing library
├── services/                 # HTTP handlers (business logic)
├── utils/                    # Path and version utilities
└── versions/                 # Version loader and file enumeration
```

## Tech Stack

- [Fiber v2](https://gofiber.io/) — HTTP framework
- [Viper](https://github.com/spf13/viper) — Configuration management
- [gopkg.in/yaml.v3](https://pkg.go.dev/gopkg.in/yaml.v3) — YAML parsing
- [charmbracelet/log](https://github.com/charmbracelet/log) — Styled logging
- [goccy/go-yaml](https://github.com/goccy/go-yaml) — YAML syntax highlighting

## License

This project is licensed under the [MIT License](LICENSE).
