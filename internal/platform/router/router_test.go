package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewServesHealthz(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	New(nil).ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	if strings.TrimSpace(rec.Body.String()) != "ok" {
		t.Fatalf("expected body %q, got %q", "ok", rec.Body.String())
	}
}
