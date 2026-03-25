package verification

import (
	"context"
	"testing"
	"time"
)

func TestServiceEndToEndCreateRequestUsesResolverAndReturnsPendingSession(t *testing.T) {
	repository := &stubRepository{
		createRequestFunc: func(_ context.Context, request VerificationRequest) (VerificationRequest, error) {
			request.ID = "req-123"
			return request, nil
		},
		updateRequestStatusFunc: func(_ context.Context, request VerificationRequest) (VerificationRequest, error) {
			return request, nil
		},
		createSessionFunc: func(_ context.Context, session VerificationSession) (VerificationSession, error) {
			return session, nil
		},
	}

	provider := &stubSessionProvider{
		createSessionResult: ProviderSession{
			ProviderSessionID: "provider-session-123",
			QRCodeURL:         "https://example.com/qr",
			DeepLink:          "openid://example",
			OfferURL:          "https://example.com/offer",
			ExpiresAt:         time.Date(2026, time.March, 25, 12, 0, 0, 0, time.UTC),
		},
	}

	resolver := &stubProviderResolver{
		provider: provider,
	}

	service := NewServiceWithResolver("walt", repository, resolver)

	response, err := service.CreateRequest(context.Background(), CreateRequestInput{
		BusinessRef:  "job-123",
		CandidateRef: "cand-456",
	})
	if err != nil {
		t.Fatalf("expected end-to-end create request to succeed, got error: %v", err)
	}

	if resolver.requestedProvider != "walt" {
		t.Fatalf("expected resolver to receive provider walt, got %q", resolver.requestedProvider)
	}

	if repository.createdSession.ProviderSessionID != "provider-session-123" {
		t.Fatalf("expected persisted provider session id provider-session-123, got %q", repository.createdSession.ProviderSessionID)
	}

	if response.Status != StatusPending {
		t.Fatalf("expected response status %q, got %q", StatusPending, response.Status)
	}

	if response.Session.DeepLink != "openid://example" {
		t.Fatalf("expected deep link to round-trip, got %q", response.Session.DeepLink)
	}
}

type stubProviderResolver struct {
	requestedProvider string
	provider          SessionProvider
}

func (resolver *stubProviderResolver) Resolve(_ context.Context, provider string) (SessionProvider, error) {
	resolver.requestedProvider = provider
	return resolver.provider, nil
}
