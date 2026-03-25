package walt

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/nathang0147/Request-Service-API/internal/verification"
)

func TestProviderCreateSessionMapsRequestAndResponse(t *testing.T) {
	var capturedRequest createSessionRequest
	httpClient := roundTripClient(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodPost {
			t.Fatalf("expected method %q, got %q", http.MethodPost, request.Method)
		}

		if request.URL.String() != "https://walt.example.com/sessions" {
			t.Fatalf("expected url %q, got %q", "https://walt.example.com/sessions", request.URL.String())
		}

		if request.Header.Get("Authorization") != "Bearer secret-token" {
			t.Fatalf("expected bearer auth header, got %q", request.Header.Get("Authorization"))
		}

		if err := json.NewDecoder(request.Body).Decode(&capturedRequest); err != nil {
			t.Fatalf("expected valid walt request body, got error: %v", err)
		}

		body, err := json.Marshal(createSessionResponse{
			SessionID: "provider-session-123",
			QRCodeURL: "https://example.com/qr",
			DeepLink:  "openid://example",
			OfferURL:  "https://example.com/offer",
			ExpiresAt: time.Date(2026, time.March, 25, 12, 0, 0, 0, time.UTC),
		})
		if err != nil {
			t.Fatalf("expected response marshal to succeed, got error: %v", err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(strings.NewReader(string(body))),
		}, nil
	})

	client := NewClient("https://walt.example.com", "secret-token", httpClient)
	provider := New(client)

	session, err := provider.CreateSession(context.Background(), verification.ProviderSessionInput{
		RequestID:    "req-123",
		BusinessRef:  "job-123",
		CandidateRef: "cand-456",
		Provider:     "walt",
	})
	if err != nil {
		t.Fatalf("expected create session to succeed, got error: %v", err)
	}

	if capturedRequest.RequestID != "req-123" {
		t.Fatalf("expected request id req-123, got %q", capturedRequest.RequestID)
	}

	if capturedRequest.BusinessRef != "job-123" || capturedRequest.CandidateRef != "cand-456" {
		t.Fatalf("expected business/candidate refs to be mapped, got %+v", capturedRequest)
	}

	if session.QRCodeURL != "https://example.com/qr" {
		t.Fatalf("expected qr code url to round-trip, got %q", session.QRCodeURL)
	}

	if session.ProviderSessionID != "provider-session-123" {
		t.Fatalf("expected provider session id provider-session-123, got %q", session.ProviderSessionID)
	}
}

type roundTripClient func(*http.Request) (*http.Response, error)

func (client roundTripClient) Do(request *http.Request) (*http.Response, error) {
	return client(request)
}

func TestProviderParseCallbackNormalizesEvent(t *testing.T) {
	provider := New(NewClient("https://walt.example.com", "secret-token", http.DefaultClient))

	event, err := provider.ParseCallback(context.Background(), []byte(`{
		"sessionId": "provider-session-123",
		"status": "VERIFIED",
		"verified": true,
		"reasonCode": "POLICY_APPROVED",
		"eventType": "SESSION_VERIFIED",
		"payload": {"verified": true}
	}`))
	if err != nil {
		t.Fatalf("expected callback parse to succeed, got error: %v", err)
	}

	if event.ProviderSessionID != "provider-session-123" {
		t.Fatalf("expected provider session id provider-session-123, got %q", event.ProviderSessionID)
	}

	if event.Status != verification.StatusVerified {
		t.Fatalf("expected normalized status %q, got %q", verification.StatusVerified, event.Status)
	}

	if !event.Verified {
		t.Fatal("expected normalized callback to be verified")
	}
}
