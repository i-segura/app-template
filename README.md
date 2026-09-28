# app-template

Cookiecutter template for a **Go backend + React frontend shipped as a single
container**. The Vite build is embedded into the Go binary with `go:embed`;
the final Docker image is a slim Alpine running as a non-root user.

## Usage

```sh
pipx install cookiecutter   # or: uv tool install cookiecutter
cookiecutter gh:i-segura/app-template
```

Answer the prompts (or use `--no-input` with defaults) and you get a project
that builds with one `make build` and runs with one binary.

## Template variables

| Variable         | Default                                  | Notes                                            |
| ---------------- | ---------------------------------------- | ------------------------------------------------ |
| `project_name`   | `My App`                                 | Human-readable name (used in README, page title) |
| `project_slug`   | slugified `project_name`                 | Directory, npm package, Docker image name        |
| `description`    | generic                                  | One-liner for README and the app UI              |
| `go_module_path` | `github.com/example/<slug>`              | Change to your real repo path after generating   |
| `node_version`   | `24`                                     | LTS line used in Dockerfile + npm engines        |
| `go_version`     | `1.26`                                   | Used in go.mod + Dockerfile                      |

Kept intentionally minimal — opinionated extras (databases, auth, CI/CD to a
registry) belong in the generated project, not the template.

## What the generated project includes

- **Go backend** (`cmd/app` + `internal/{config,httpserver,ui}`): stdlib
  `net/http` with `http.ServeMux` patterns, JSON health endpoint `/live`,
  example API `/api/v1/hello`, SPA fallback to `index.html` for client-side
  routes (JSON 404s under `/api/*`), slog JSON logging, graceful shutdown on
  SIGINT/SIGTERM.
- **React frontend** (`web/`): React 19 + TypeScript + Vite, dev-server proxy
  for `/api` and `/live` so dev behaves like production.
- **Embed pipeline**: `web/dist` → `internal/ui/dist` → `//go:embed all:dist`.
  A committed placeholder keeps `go build`/`go test` green before any
  frontend build.
- **Dockerfile**: 3 stages (node builds frontend → go builds binary →
  alpine + tini, non-root uid 65532, `HEALTHCHECK` on `/live`).
- **Makefile**: `dev`, `build`, `test`, `lint`, `docker`, `run`, `clean`.
- **Quality gates**: `go test` with `httptest` handler tests, Vitest +
  Testing Library example tests, ESLint (flat config) + Prettier,
  golangci-lint v2 config.
- **GitHub Actions CI**: Go lint+test (race, coverage), frontend lint+test,
  and a `docker build` validation job on push/PR.

## Why these defaults

- **Single binary/container** — one artifact to deploy, no CORS, no nginx
  sidecar; the embed pattern is proven in production apps.
- **stdlib `net/http` over gin** — a template should teach the platform;
  ServeMux method patterns cover what a small API needs. Swap in a router
  later if you need middleware-heavy routing.
- **No database in the template** — every app picks its own storage; the
  README shows where to add one.
- **Node 24 / Go 1.26** — current LTS/stable lines at template-writing time;
  both are template variables.
- **Vitest over Jest** — first-class Vite integration, no babel config.

## Hacking on the template

Render locally to verify changes:

```sh
cookiecutter --no-input -o /tmp/render-test project_name="Demo App"
cd /tmp/render-test/demo-app
go build ./... && go test ./...
cd web && npm ci && npm run build && npm test && npm run lint
```
