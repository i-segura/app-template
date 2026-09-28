# {{ cookiecutter.project_name }}

{{ cookiecutter.description }}

Go backend + React frontend compiled into **one static binary** (the frontend
build is embedded with `go:embed`) and shipped as one slim container.

## Requirements

- Go {{ cookiecutter.go_version }}+
- Node {{ cookiecutter.node_version }}+ and npm

## Quickstart

```sh
make build     # build frontend -> copy to internal/ui/dist -> build Go binary (bin/app)
./bin/app      # serve everything on http://localhost:8080
```

### Development

```sh
make dev       # prints the two commands to run
```

Run them in two terminals:

```sh
go run ./cmd/app        # backend on :8080
cd web && npm run dev   # Vite on :5173, proxies /api and /live to :8080
```

### Test / lint / container

```sh
make test      # go test ./... + vitest
make lint      # golangci-lint (if installed) + eslint + prettier
make docker    # docker build (multi-stage, non-root final image)
make run       # build image and run it on :8080
```

## How the single-container pattern works

1. Vite builds the React app to `web/dist`.
2. `web/dist` is copied to `internal/ui/dist` (by `make build` and by the
   Dockerfile).
3. `internal/ui/static.go` embeds that directory with `//go:embed all:dist`.
4. The Go server serves embedded assets and falls back to `index.html` for
   unknown non-`/api` routes, so client-side routing works.

A placeholder `internal/ui/dist/index.html` is committed so `go build` and
`go test` work before the frontend has ever been built.

## Layout

```
cmd/app/                 thin entrypoint: config, slog, graceful shutdown
internal/config/         env-var config (PORT, LOG_LEVEL)
internal/httpserver/     routes: /live, /api/v1/*, SPA static fallback
internal/ui/             go:embed of the frontend build (dist/)
web/                     React 19 + TypeScript + Vite (+ vitest, eslint, prettier)
Dockerfile               node -> go -> alpine (non-root, tini)
```

## Configuration

| Env var     | Default | Meaning                    |
| ----------- | ------- | -------------------------- |
| `PORT`      | `8080`  | HTTP listen port           |
| `LOG_LEVEL` | `info`  | `debug` \| `info` \| `warn` \| `error` |

## Extending

- Add API handlers in `internal/httpserver` under the `/api/v1` pattern.
- Add a database: create `internal/db`, add `DATABASE_URL` to
  `internal/config`, and a `/ready` probe that pings it.
- CI (`.github/workflows/ci.yml`) runs Go lint+test, frontend lint+test, and
  a `docker build` validation on every push/PR.
