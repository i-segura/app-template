# syntax=docker/dockerfile:1

FROM node:{{ cookiecutter.node_version }}-alpine AS web
WORKDIR /app/web
COPY web/package.json web/package-lock.json ./
RUN npm ci
COPY web/ ./
RUN npm run build

FROM golang:{{ cookiecutter.go_version }}-alpine AS build
RUN apk add --no-cache ca-certificates git
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
# Replaces the committed placeholder so go:embed picks up the real frontend build.
COPY --from=web /app/web/dist ./internal/ui/dist
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/app ./cmd/app

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tini tzdata \
 && adduser -D -u 65532 appuser
USER appuser
ARG PORT=8080
ENV PORT=${PORT}
EXPOSE ${PORT}
HEALTHCHECK --interval=30s --timeout=3s CMD wget -qO- "http://127.0.0.1:$PORT/live" || exit 1
COPY --from=build /out/app /usr/local/bin/app
ENTRYPOINT ["tini", "--"]
CMD ["app"]
