package callback

import (
	"context"
	"testing"

	"github.com/nathang0147/Request-Service-API/internal/provider"
	"github.com/nathang0147/Request-Service-API/internal/verification"
)

func TestServiceHandleCallbackParsesPayloadUpdatesRequestAndPersistsEvent(t *testing.T) {
	repository := &stubRepository{
		getSessionByProviderSessionIDFunc: func(_ context.Context, providerSessionID string) (SessionRecord, error) {
			if providerSessionID != "provider-session-123" {
				t.Fatalf("expected provider session id provider-session-123, got %q", providerSessionID)
			}

			return SessionRecord{
				VerificationRequestID: "req-123",
				ProviderSessionID:     providerSessionID,
				Provider:              "walt",
			}, nil
		},
		updateRequestStatusFunc: func(_ context.Context, update RequestStatusUpdate) (RequestStatusUpdate, error) {
			return update, nil
		},
		createEventFunc: func(_ context.Context, event EventRecord) (EventRecord, error) {
			return event, nil
		},
	}
	providerAdapter := &stubProvider{
		parseCallbackResult: provider.CallbackEvent{
			ProviderSessionID: "provider-session-123",
			Status:            verification.StatusVerified,
			Verified:          true,
			ReasonCode:        "POLICY_APPROVED",
			EventType:         "SESSION_VERIFIED",
			Payload:           []byte(`{"verified":true}`),
		},
	}

	service := NewService(repository, providerAdapter)

	result, err := service.HandleCallback(context.Background(), []byte(`{"sessionId":"provider-session-123"}`))
	if err != nil {
		t.Fatalf("expected callback handling to succeed, got error: %v", err)
	}

	if repository.updatedRequest.RequestID != "req-123" {
		t.Fatalf("expected request id req-123 to be updated, got %q", repository.updatedRequest.RequestID)
	}

	if repository.updatedRequest.Status != verification.StatusVerified {
		t.Fatalf("expected verified status to be persisted, got %q", repository.updatedRequest.Status)
	}

	if repository.createdEvent.EventType != "SESSION_VERIFIED" {
		t.Fatalf("expected audit event SESSION_VERIFIED, got %q", repository.createdEvent.EventType)
	}

	if result.Status != "accepted" {
		t.Fatalf("expected callback result accepted, got %q", result.Status)
	}
}

type stubRepository struct {
	updatedRequest                    RequestStatusUpdate
	createdEvent                      EventRecord
	getSessionByProviderSessionIDFunc func(context.Context, string) (SessionRecord, error)
	updateRequestStatusFunc           func(context.Context, RequestStatusUpdate) (RequestStatusUpdate, error)
	createEventFunc                   func(context.Context, EventRecord) (EventRecord, error)
}

func (repository *stubRepository) GetSessionByProviderSessionID(ctx context.Context, providerSessionID string) (SessionRecord, error) {
	return repository.getSessionByProviderSessionIDFunc(ctx, providerSessionID)
}

func (repository *stubRepository) UpdateRequestStatus(ctx context.Context, update RequestStatusUpdate) (RequestStatusUpdate, error) {
	repository.updatedRequest = update
	return repository.updateRequestStatusFunc(ctx, update)
}

func (repository *stubRepository) CreateEvent(ctx context.Context, event EventRecord) (EventRecord, error) {
	repository.createdEvent = event
	return repository.createEventFunc(ctx, event)
}

type stubProvider struct {
	parseCallbackResult provider.CallbackEvent
	parseCallbackErr    error
}

func (adapter *stubProvider) CreateSession(context.Context, verification.ProviderSessionInput) (verification.ProviderSession, error) {
	return verification.ProviderSession{}, nil
}

func (adapter *stubProvider) ParseCallback(_ context.Context, _ []byte) (provider.CallbackEvent, error) {
	if adapter.parseCallbackErr != nil {
		return provider.CallbackEvent{}, adapter.parseCallbackErr
	}

	return adapter.parseCallbackResult, nil
}
