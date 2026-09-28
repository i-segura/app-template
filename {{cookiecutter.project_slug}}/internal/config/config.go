// Package config loads application configuration from environment variables.
package config

import (
	"log/slog"
	"os"
	"strings"
)

// Config holds runtime configuration.
type Config struct {
	Port     string
	LogLevel slog.Level
}

// Load reads configuration from the environment.
//
//	PORT       listen port (default "8080")
//	LOG_LEVEL  debug | info | warn | error (default "info")
func Load() Config {
	cfg := Config{
		Port:     getenv("PORT", "8080"),
		LogLevel: parseLevel(getenv("LOG_LEVEL", "info")),
	}
	return cfg
}

func getenv(key, fallback string) string {
	if v := strings.TrimSpace(os.Getenv(key)); v != "" {
		return v
	}
	return fallback
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
