package httpserver

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLive(t *testing.T) {
	srv := httptest.NewServer(New())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/live")
	if err != nil {
		t.Fatalf("GET /live: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["status"] != "ok" {
		t.Fatalf("status field = %q, want %q", body["status"], "ok")
	}
}

func TestHello(t *testing.T) {
	srv := httptest.NewServer(New())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/hello?name=ada")
	if err != nil {
		t.Fatalf("GET /api/v1/hello: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	var body map[string]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}
	if body["message"] != "hello, ada" {
		t.Fatalf("message = %q, want %q", body["message"], "hello, ada")
	}
}

func TestUnknownAPIRouteIsJSON404(t *testing.T) {
	srv := httptest.NewServer(New())
	defer srv.Close()

	resp, err := http.Get(srv.URL + "/api/v1/nope")
	if err != nil {
		t.Fatalf("GET /api/v1/nope: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("status = %d, want 404", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Fatalf("content-type = %q, want JSON", ct)
	}
}

func TestSPAFallbackServesIndex(t *testing.T) {
	srv := httptest.NewServer(New())
	defer srv.Close()

	// A client-side route that is not a real file must return index.html.
	resp, err := http.Get(srv.URL + "/some/client/route")
	if err != nil {
		t.Fatalf("GET /some/client/route: %v", err)
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}
	if ct := resp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/html") {
		t.Fatalf("content-type = %q, want HTML", ct)
	}
}
