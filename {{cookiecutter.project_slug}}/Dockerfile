# syntax=docker/dockerfile:1

# ---------- Stage 1: build the React frontend ----------
FROM node:{{ cookiecutter.node_version }}-alpine AS web
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

# ---------- Stage 2: build the Go backend (embeds the frontend) ----------
FROM golang:{{ cookiecutter.go_version }}-alpine AS build
RUN apk add --no-cache ca-certs git
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Replace the placeholder with the real frontend build before compiling,
# so go:embed picks it up.
COPY --from=web /app/web/dist ./internal/ui/dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/app

# ---------- Final: slim non-root image ----------
FROM alpine:3.21
RUN apk add --no-cache ca-certs tini tzdata \
 && adduser -D -u 65532 appuser
USER appuser
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- http://127.0.0.1:8080/live || exit 1
COPY --from=build /out/app /usr/local/bin/app
ENTRYPOINT ["tini", "--"]
CMD ["app"]
