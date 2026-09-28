# {{ cookiecutter.project_name }}

{{ cookiecutter.description }}

Go backend + React frontend compiled into one static binary (frontend embedded
with `go:embed`) and shipped as one container.

## Requirements

[mise](https://mise.jdx.dev). In this directory:

```sh
mise trust      # first time only
mise install    # installs pinned Go {{ cookiecutter.go_version }}, Node {{ cookiecutter.node_version }}, golangci-lint
```

## Quickstart

```sh
mise run build
./bin/app       # http://localhost:8080
```

## Tasks

```sh
mise run dev       # backend (:8080) + Vite (:5173) concurrently
mise run test      # go test -race + vitest
mise run lint      # golangci-lint + eslint + prettier
mise run check     # lint + test
mise run fmt       # gofmt + prettier
mise run docker    # docker build
mise run run       # build image and run on :8080
mise run clean     # remove build artifacts
```

`mise tasks` lists everything with descriptions.

## Layout

```
cmd/app/                 entrypoint: config, slog, graceful shutdown
internal/config/         env-var config (PORT, LOG_LEVEL)
internal/httpserver/     routes: /live, /api/v1/*, SPA static fallback
internal/ui/             go:embed of the frontend build (dist/; committed
                         placeholder keeps go build/test green pre-frontend)
web/                     React 19 + TypeScript + Vite (vitest, eslint, prettier)
mise.toml                pinned toolchain + all project tasks
Dockerfile               node -> go -> alpine (non-root, tini)
```

`mise run build` and the Dockerfile both copy `web/dist` into
`internal/ui/dist` before compiling the binary.

## Configuration

| Env var     | Default | Meaning                    |
| ----------- | ------- | -------------------------- |
| `PORT`      | `8080`  | HTTP listen port           |
| `LOG_LEVEL` | `info`  | `debug` \| `info` \| `warn` \| `error` |

## Extending

- Add API handlers in `internal/httpserver` under the `/api/v1` pattern.
- Add a database: create `internal/db`, add `DATABASE_URL` to
  `internal/config`, and a `/ready` probe that pings it.
