package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewServerUsesConfiguredPort(t *testing.T) {
	server := newServer("8080", http.NewServeMux())

	if server.Addr != ":8080" {
		t.Fatalf("expected server address %q, got %q", ":8080", server.Addr)
	}
}

func TestNewHandlerServesHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	newHandler(nil).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if strings.TrimSpace(rec.Body.String()) != "ok" {
		t.Fatalf("expected body %q, got %q", "ok", rec.Body.String())
	}
}
