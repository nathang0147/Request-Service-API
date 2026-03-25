package callback

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	chi "github.com/go-chi/chi/v5"
)

func TestHandlerRejectsUnauthorizedCallback(t *testing.T) {
	router := chi.NewRouter()
	NewHandler(&stubCallbackService{}, rejectingAuthenticator{}).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/callbacks/walt", bytes.NewBufferString(`{}`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("expected status %d, got %d", http.StatusUnauthorized, rec.Code)
	}
}

func TestHandlerReturnsAcceptedResponse(t *testing.T) {
	router := chi.NewRouter()
	NewHandler(&stubCallbackService{
		handleCallbackFunc: func(_ context.Context, body []byte) (Result, error) {
			if string(body) != `{"sessionId":"provider-session-123"}` {
				t.Fatalf("expected callback body to round-trip, got %q", string(body))
			}

			return Result{Status: "accepted"}, nil
		},
	}, staticAuthenticator{secret: "secret"}).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/callbacks/walt", bytes.NewBufferString(`{"sessionId":"provider-session-123"}`))
	req.Header.Set(callbackSecretHeader, "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusAccepted {
		t.Fatalf("expected status %d, got %d", http.StatusAccepted, rec.Code)
	}
}

func TestHandlerReturnsInternalErrorOnServiceFailure(t *testing.T) {
	router := chi.NewRouter()
	NewHandler(&stubCallbackService{
		handleCallbackFunc: func(_ context.Context, _ []byte) (Result, error) {
			return Result{}, errors.New("database offline")
		},
	}, staticAuthenticator{secret: "secret"}).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/callbacks/walt", bytes.NewBufferString(`{"sessionId":"provider-session-123"}`))
	req.Header.Set(callbackSecretHeader, "secret")
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

type stubCallbackService struct {
	handleCallbackFunc func(context.Context, []byte) (Result, error)
}

func (service *stubCallbackService) HandleCallback(ctx context.Context, body []byte) (Result, error) {
	if service.handleCallbackFunc == nil {
		return Result{}, nil
	}

	return service.handleCallbackFunc(ctx, body)
}

type rejectingAuthenticator struct{}

func (rejectingAuthenticator) Authenticate(*http.Request) error {
	return errors.New("unauthorized")
}
