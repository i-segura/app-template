# {{ cookiecutter.project_name }}

{{ cookiecutter.description }}

Go backend + React frontend compiled into **one static binary** (the frontend
build is embedded with `go:embed`) and shipped as one slim container.

## Requirements

Just [mise](https://mise.jdx.dev) — everything else is pinned in `mise.toml`:

```sh
curl https://mise.run | sh   # or: brew install mise / apt etc.
mise install                 # installs Go {{ cookiecutter.go_version }}, Node {{ cookiecutter.node_version }}, golangci-lint for this directory
```

## Quickstart

```sh
mise run build     # frontend -> internal/ui/dist -> Go binary (bin/app)
./bin/app          # serve everything on http://localhost:8080
```

### Development

```sh
mise run dev       # backend (:8080) + Vite (:5173) concurrently, Ctrl-C stops both
```

Vite proxies `/api` and `/live` to the backend, so dev behaves like production.

### Test / lint / container

```sh
mise run test      # go test -race + vitest (in parallel)
mise run lint      # golangci-lint + eslint + prettier (in parallel)
mise run check     # lint + test together
mise run docker    # docker build (multi-stage, non-root final image)
mise run run       # build image and run it on :8080
mise run clean     # remove build artifacts
```

## mise notes

- `mise.toml` pins the toolchain in `[tools]` and defines project commands in
  `[tasks.*]`.
- `mise tasks` lists every task with its description.
- First run in a new repo requires `mise trust`.
- Inside the directory, `go`, `node`, `npm`, and `golangci-lint` on PATH are
  the pinned versions. Tasks with `depends` run in parallel.

## How the single-container pattern works

1. Vite builds the React app to `web/dist`.
2. `web/dist` is copied to `internal/ui/dist` (by `mise run build` and by the
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
mise.toml                pinned toolchain + all project tasks
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

## Design decisions

- Single binary/container: one deployable artifact, no CORS, no reverse proxy.
- stdlib `net/http` with `ServeMux` method patterns; a router can be added
  later if middleware-heavy routing is needed.
- No database: storage is app-specific.
- CI uses `actions/setup-go` / `actions/setup-node` reading the pinned
  versions rather than the mise action, so the pipeline does not depend on
  mise on GitHub runners.
