package verification

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	chi "github.com/go-chi/chi/v5"
)

func TestHandlerCreateRequest(t *testing.T) {
	service := &stubRequestService{
		createRequestFunc: func(_ context.Context, input CreateRequestInput) (CreateRequestResponse, error) {
			if input.BusinessRef != "job-123" || input.CandidateRef != "cand-456" {
				t.Fatalf("expected normalized create input, got %+v", input)
			}

			if len(input.CredentialTypes) != 2 || input.CredentialTypes[0] != "DiplomaCredential" || input.CredentialTypes[1] != "TranscriptCredential" {
				t.Fatalf("expected credential types to decode, got %+v", input.CredentialTypes)
			}

			return CreateRequestResponse{
				RequestID: "req-123",
				Status:    StatusPending,
				Verified:  false,
				Session: SessionDetails{
					QRCodeURL: "https://example.com/qr",
					DeepLink:  "openid://example",
					OfferURL:  "https://example.com/offer",
					ExpiresAt: time.Date(2026, time.March, 25, 12, 0, 0, 0, time.UTC),
				},
			}, nil
		},
	}

	router := chi.NewRouter()
	NewHandler(service).RegisterRoutes(router)

	body := bytes.NewBufferString(`{"businessRef":"job-123","candidateRef":"cand-456","credentialTypes":["DiplomaCredential","TranscriptCredential"]}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/verification-requests", body)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var response CreateRequestResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("expected valid json response, got error: %v", err)
	}

	if response.RequestID != "req-123" {
		t.Fatalf("expected request id req-123, got %q", response.RequestID)
	}

	if response.Session.QRCodeURL != "https://example.com/qr" {
		t.Fatalf("expected qr code url to round-trip, got %q", response.Session.QRCodeURL)
	}
}

func TestHandlerGetRequest(t *testing.T) {
	service := &stubRequestService{
		getRequestFunc: func(_ context.Context, requestID string) (GetRequestResponse, error) {
			if requestID != "req-123" {
				t.Fatalf("expected request id req-123, got %q", requestID)
			}

			return GetRequestResponse{
				RequestID: "req-123",
				Status:    StatusVerified,
				Verified:  true,
				Session: &SessionDetails{
					QRCodeURL: "https://example.com/qr",
				},
			}, nil
		},
	}

	router := chi.NewRouter()
	NewHandler(service).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/verification-requests/req-123", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response GetRequestResponse
	if err := json.NewDecoder(rec.Body).Decode(&response); err != nil {
		t.Fatalf("expected valid json response, got error: %v", err)
	}

	if !response.Verified {
		t.Fatal("expected response to be verified")
	}
}

func TestHandlerRejectsInvalidJSON(t *testing.T) {
	router := chi.NewRouter()
	NewHandler(&stubRequestService{}).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/verification-requests", bytes.NewBufferString(`{"businessRef":`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandlerRejectsMissingFields(t *testing.T) {
	router := chi.NewRouter()
	NewHandler(&stubRequestService{}).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/verification-requests", bytes.NewBufferString(`{"businessRef":"","candidateRef":""}`))
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rec.Code)
	}
}

func TestHandlerGetRequestReturnsNotFound(t *testing.T) {
	service := &stubRequestService{
		getRequestFunc: func(_ context.Context, _ string) (GetRequestResponse, error) {
			return GetRequestResponse{}, &Error{
				Code: ErrCodeRequestNotFound,
				Err:  errors.New("not found"),
			}
		},
	}

	router := chi.NewRouter()
	NewHandler(service).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/verification-requests/req-404", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status %d, got %d", http.StatusNotFound, rec.Code)
	}
}

func TestHandlerReturnsNormalizedInternalError(t *testing.T) {
	service := &stubRequestService{
		getRequestFunc: func(_ context.Context, _ string) (GetRequestResponse, error) {
			return GetRequestResponse{}, errors.New("database offline")
		},
	}

	router := chi.NewRouter()
	NewHandler(service).RegisterRoutes(router)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/verification-requests/req-500", nil)
	rec := httptest.NewRecorder()

	router.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
	}
}

type stubRequestService struct {
	createRequestFunc func(context.Context, CreateRequestInput) (CreateRequestResponse, error)
	getRequestFunc    func(context.Context, string) (GetRequestResponse, error)
}

func (service *stubRequestService) CreateRequest(ctx context.Context, input CreateRequestInput) (CreateRequestResponse, error) {
	if service.createRequestFunc == nil {
		return CreateRequestResponse{}, nil
	}

	return service.createRequestFunc(ctx, input)
}

func (service *stubRequestService) GetRequest(ctx context.Context, requestID string) (GetRequestResponse, error) {
	if service.getRequestFunc == nil {
		return GetRequestResponse{}, nil
	}

	return service.getRequestFunc(ctx, requestID)
}
