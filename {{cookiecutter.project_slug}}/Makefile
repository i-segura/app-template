# {{ cookiecutter.project_name }} — task runner
# Requires: go {{ cookiecutter.go_version }}+, node {{ cookiecutter.node_version }}+ (see README).

APP := $(notdir $(CURDIR))
BIN := bin/app

.PHONY: all dev build test lint fmt docker run clean web-deps

all: build

## dev: run backend + frontend dev servers (two processes; see README)
dev:
	@echo "Run these in two terminals:"
	@echo "  1. go run ./cmd/app"
	@echo "  2. cd web && npm run dev"
	@echo "Frontend on http://localhost:5173 (proxies /api and /live to :8080)"

## web-deps: install frontend dependencies
web-deps:
	cd web && npm ci

## build: build frontend, copy into internal/ui/dist, then build the Go binary
build:
	cd web && npm ci && npm run build
	find internal/ui/dist -mindepth 1 ! -name .gitignore -delete
	cp -r web/dist/* internal/ui/dist/
	go build -trimpath -ldflags "-s -w" -o $(BIN) ./cmd/app

## test: run Go tests and frontend tests
test:
	go test -race ./cmd/... ./internal/...
	cd web && npm ci && npm test

## lint: run golangci-lint (if installed) + eslint/prettier
lint:
	@command -v golangci-lint >/dev/null 2>&1 && golangci-lint run ./... || echo "golangci-lint not installed, skipping Go lint"
	cd web && npm ci && npm run lint

## fmt: format Go code
fmt:
	gofmt -w .

## docker: build the single-container image
docker:
	docker build -t $(APP):latest .

## run: build and run the container on :8080
run: docker
	docker run --rm -p 8080:8080 $(APP):latest

## clean: remove build artifacts
clean:
	rm -rf bin web/dist internal/ui/dist/assets
