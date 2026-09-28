# app-template

Cookiecutter template for a **Go backend + React frontend shipped as a single
container**. The Vite build is embedded into the Go binary with `go:embed`;
the final Docker image is a slim Alpine running as a non-root user.

## Usage

```sh
curl https://mise.run | sh          # install mise once
git clone gh:i-segura/app-template && cd app-template
mise install                        # pinned cookiecutter + Go + Node for this repo
mise run render                     # render with defaults into /tmp/app-template-render
```

Or point cookiecutter at the repo directly (`cookiecutter gh:i-segura/app-template`)
with any cookiecutter install. Answer the prompts (or use `--no-input` with
defaults) and you get a project that builds with one `mise run build` and runs
with one binary.

## Template variables

| Variable         | Default                                  | Notes                                            |
| ---------------- | ---------------------------------------- | ------------------------------------------------ |
| `project_name`   | `My App`                                 | Human-readable name (used in README, page title) |
| `project_slug`   | slugified `project_name`                 | Directory, npm package, Docker image name        |
| `description`    | generic                                  | One-liner for README and the app UI              |
| `go_module_path` | `github.com/example/<slug>`              | Change to your real repo path after generating   |
| `node_version`   | `24`                                     | Feeds mise.toml `[tools]` + Dockerfile + npm engines |
| `go_version`     | `1.26`                                   | Feeds mise.toml `[tools]` + go.mod + Dockerfile   |

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
- **mise**: `mise.toml` pins the toolchain (`[tools]`: Go, Node,
  golangci-lint) and defines all tasks (`[tasks]`: `dev`, `build`, `test`,
  `lint`, `fmt`, `check`, `docker`, `run`, `clean`).
- **Quality gates**: `go test` with `httptest` handler tests, Vitest +
  Testing Library example tests, ESLint (flat config) + Prettier,
  golangci-lint v2 config.
- **GitHub Actions CI**: Go lint+test (race, coverage), frontend lint+test,
  and a `docker build` validation job on push/PR.

## Design decisions

- Single binary/container: one deployable artifact, no CORS, no reverse proxy.
- stdlib `net/http` with `ServeMux` method patterns; a router can be added
  later if middleware-heavy routing is needed.
- No database: storage is app-specific.
- mise instead of make: `[tools]` pins the toolchain per directory,
  `[tasks]` defines the commands, `depends` runs tasks in parallel.
- Vitest for frontend tests: first-class Vite integration.

## Hacking on the template

This repo dogfoods mise: its own `mise.toml` pins cookiecutter plus the same
Go/Node lines the generated project gets. Verify changes end to end:

```sh
mise install
mise run verify    # render + mise run build/test/lint inside the rendered project
```

Or step by step:

```sh
cookiecutter --no-input . -o /tmp/render-test project_name="Demo App"
cd /tmp/render-test/demo-app
mise trust -y && mise install
mise run build && mise run test && mise run lint
```
