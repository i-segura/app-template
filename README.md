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
- **mise**: `mise.toml` pins the toolchain (`[tools]`: Go, Node,
  golangci-lint) **and** defines all tasks (`[tasks]`: `dev`, `build`, `test`,
  `lint`, `fmt`, `check`, `docker`, `run`, `clean`). One tool replaces the
  Makefile *and* the "install these versions yourself" README section.
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
  both are template variables and both land in the generated `mise.toml`
  `[tools]`, so the pinned versions are *enforced*, not just documented.
- **mise over make** — mise is make plus a toolchain manager: `[tools]` pins
  Go/Node/golangci-lint per directory (`mise install` reproduces the exact
  environment), `[tasks]` replaces the Makefile, and task `depends` gives
  free parallelism (`test` runs Go and frontend suites concurrently). One
  manifest, no drift between "what the README says to install" and "what the
  CI image has".
- **Vitest over Jest** — first-class Vite integration, no babel config.

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
