package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRequestIDAddsHeaderAndContextValue(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := RequestIDFromContext(r.Context())
		if requestID == "" {
			t.Fatal("expected request id in context")
		}

		if got := w.Header().Get(RequestIDHeader); got != requestID {
			t.Fatalf("expected response header to match context value, got header=%q context=%q", got, requestID)
		}

		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, rec.Code)
	}

	if rec.Header().Get(RequestIDHeader) == "" {
		t.Fatalf("expected response header %q to be set", RequestIDHeader)
	}
}

func TestRequestIDPreservesInboundHeader(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := RequestIDFromContext(r.Context()); got != "external-id" {
			t.Fatalf("expected request id %q in context, got %q", "external-id", got)
		}

		w.WriteHeader(http.StatusNoContent)
	}))

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(RequestIDHeader, "external-id")
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if got := rec.Header().Get(RequestIDHeader); got != "external-id" {
		t.Fatalf("expected response header %q, got %q", "external-id", got)
	}
}
