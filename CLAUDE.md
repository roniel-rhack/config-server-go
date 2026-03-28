# Config Server Go

## Project Overview

Centralized configuration management microservice in Go (Fiber). Serves versioned YAML/properties configs over HTTP.

## Tech Stack

- Go 1.22+, Fiber v2, Viper, yaml.v3
- Docker multi-stage build (golang:1.22 → alpine:3)
- Module name: `configTest` (in go.mod)

## Project Structure

- `src/config_server.go` — entry point, routes, middleware
- `src/config/` — Viper setup, defaults, file watching
- `src/services/` — HTTP handlers (properties, versions)
- `src/models/` — request/response structs
- `src/pkg/` — YAML/properties parsing and transformation
- `src/utils/` — path and version helpers
- `src/versions/` — version loader, file enumeration
- `src/custom_logguer/` — styled terminal logging

## Build & Run

```bash
cd src && go run config_server.go     # local
docker compose up --build             # docker
```

Default port: 8888

## Conventions

- Use `go vet` and `go fmt` before committing
- Run `go build ./...` to verify compilation
- No tests exist yet — add tests in `*_test.go` files alongside source
- Keep handlers in `services/`, models in `models/`, utilities in `utils/`

## Workflow Orchestration

### 1. Plan Mode Default
- Enter plan mode for ANY non-trivial task (3+ steps or architectural decisions)
- If something goes sideways, STOP and re-plan immediately — don't keep pushing
- Use plan mode for verification steps, not just building
- Write detailed specs upfront to reduce ambiguity

### 2. Subagent Strategy
- Use subagents liberally to keep main context window clean
- Offload research, exploration, and parallel analysis to subagents
- For complex problems, throw more compute at it via subagents
- One task per subagent for focused execution

### 3. Self-Improvement Loop
- After ANY correction from the user: update `tasks/lessons.md` with the pattern
- Write rules for yourself that prevent the same mistake
- Ruthlessly iterate on these lessons until mistake rate drops
- Review lessons at session start for relevant project

### 4. Verification Before Done
- Never mark a task complete without proving it works
- Diff behavior between main and your changes when relevant
- Ask yourself: "Would a staff engineer approve this?"
- Run tests, check logs, demonstrate correctness

### 5. Demand Elegance (Balanced)
- For non-trivial changes: pause and ask "is there a more elegant way?"
- If a fix feels hacky: "Knowing everything I know now, implement the elegant solution"
- Skip this for simple, obvious fixes — don't over-engineer
- Challenge your own work before presenting it

### 6. Autonomous Bug Fixing
- When given a bug report: just fix it. Don't ask for hand-holding
- Point at logs, errors, failing tests — then resolve them
- Zero context switching required from the user
- Go fix failing CI tests without being told how

## Task Management

1. **Plan First**: Write plan to `tasks/todo.md` with checkable items
2. **Verify Plan**: Check in before starting implementation
3. **Track Progress**: Mark items complete as you go
4. **Explain Changes**: High-level summary at each step
5. **Document Results**: Add review section to `tasks/todo.md`
6. **Capture Lessons**: Update `tasks/lessons.md` after corrections

## Core Principles

- **Simplicity First**: Make every change as simple as possible. Impact minimal code.
- **No Laziness**: Find root causes. No temporary fixes. Senior developer standards.
- **Minimal Impact**: Changes should only touch what's necessary. Avoid introducing bugs.
