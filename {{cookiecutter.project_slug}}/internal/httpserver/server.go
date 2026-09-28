// Package httpserver wires routes: JSON API under /api/v1, health probes,
// and the embedded SPA with fallback to index.html.
package httpserver

import (
	"encoding/json"
	"io/fs"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"{{ cookiecutter.go_module_path }}/internal/ui"
)

// New returns the application http.Handler.
func New() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /live", handleLive)
	mux.HandleFunc("GET /api/v1/hello", handleHello)

	mux.HandleFunc("/", handleStatic(ui.Assets()))

	return requestLog(mux)
}

func handleLive(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleHello(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		name = "world"
	}
	writeJSON(w, http.StatusOK, map[string]string{"message": "hello, " + name})
}

// handleStatic serves files from the embedded frontend build and falls back
// to index.html for unknown non-API paths so client-side routing works (SPA).
func handleStatic(assets fs.FS) http.HandlerFunc {
	fileServer := http.FileServer(http.FS(assets))
	index, _ := fs.ReadFile(assets, "index.html")

	return func(w http.ResponseWriter, r *http.Request) {
		upath := strings.TrimPrefix(r.URL.Path, "/")
		if upath == "" {
			upath = "index.html"
		}

		if _, err := fs.Stat(assets, upath); err == nil {
			fileServer.ServeHTTP(w, r)
			return
		}

		// Unknown API routes stay JSON 404s, never the SPA shell.
		if strings.HasPrefix(r.URL.Path, "/api/") {
			writeJSON(w, http.StatusNotFound, map[string]string{"error": "not found"})
			return
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		http.ServeContent(w, r, "index.html", time.Time{}, strings.NewReader(string(index)))
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// requestLog is minimal access logging; swap for real middleware as needed.
func requestLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
		next.ServeHTTP(rec, r)
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration_ms", time.Since(start).Milliseconds(),
		)
	})
}

// statusRecorder captures the response code for access logging.
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}
