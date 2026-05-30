package verification

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestServiceCreateRequestStartsInCreatedState(t *testing.T) {
	repository := &stubRepository{
		createRequestFunc: func(_ context.Context, request VerificationRequest) (VerificationRequest, error) {
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
			ExpiresAt:         time.Date(2026, time.March, 25, 12, 0, 0, 0, time.UTC),
		},
	}

	service := NewService("walt", repository, provider)

	_, err := service.CreateRequest(context.Background(), CreateRequestInput{
		BusinessRef:  "job-123",
		CandidateRef: "cand-456",
	})
	if err != nil {
		t.Fatalf("expected create request to succeed, got error: %v", err)
	}

	if repository.createdRequest.Status != StatusCreated {
		t.Fatalf("expected initial status %q, got %q", StatusCreated, repository.createdRequest.Status)
	}

	if repository.createdRequest.Provider != "walt" {
		t.Fatalf("expected provider walt, got %q", repository.createdRequest.Provider)
	}
}

func TestServiceCreateRequestPersistsSessionAndReturnsPendingResponse(t *testing.T) {
	repository := &stubRepository{
		createRequestFunc: func(_ context.Context, request VerificationRequest) (VerificationRequest, error) {
			request.ID = "req-123"
			request.CreatedAt = time.Date(2026, time.March, 25, 8, 0, 0, 0, time.UTC)
			request.UpdatedAt = request.CreatedAt
			return request, nil
		},
		updateRequestStatusFunc: func(_ context.Context, request VerificationRequest) (VerificationRequest, error) {
			request.UpdatedAt = time.Date(2026, time.March, 25, 8, 1, 0, 0, time.UTC)
			return request, nil
		},
		createSessionFunc: func(_ context.Context, session VerificationSession) (VerificationSession, error) {
			session.ID = "sess-123"
			session.CreatedAt = time.Date(2026, time.March, 25, 8, 0, 30, 0, time.UTC)
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
			RawCreateResponse: []byte(`{"session":"ok"}`),
		},
	}

	service := NewService("walt", repository, provider)

	response, err := service.CreateRequest(context.Background(), CreateRequestInput{
		BusinessRef:     "job-123",
		CandidateRef:    "cand-456",
		CredentialTypes: []string{"DiplomaCredential", "TranscriptCredential"},
	})
	if err != nil {
		t.Fatalf("expected create request to succeed, got error: %v", err)
	}

	if len(provider.createSessionInput.CredentialTypes) != 2 {
		t.Fatalf("expected credential types to reach provider, got %+v", provider.createSessionInput.CredentialTypes)
	}

	if provider.createSessionInput.CredentialTypes[0] != "DiplomaCredential" || provider.createSessionInput.CredentialTypes[1] != "TranscriptCredential" {
		t.Fatalf("expected requested credential types to round-trip, got %+v", provider.createSessionInput.CredentialTypes)
	}

	if repository.createdSession.ProviderSessionID != "provider-session-123" {
		t.Fatalf("expected provider session to be persisted, got %q", repository.createdSession.ProviderSessionID)
	}

	if response.Status != StatusPending {
		t.Fatalf("expected response status %q, got %q", StatusPending, response.Status)
	}

	if response.Session.QRCodeURL != "https://example.com/qr" {
		t.Fatalf("expected qr code url to round-trip, got %q", response.Session.QRCodeURL)
	}
}

func TestServiceGetRequestReturnsNormalizedStatus(t *testing.T) {
	repository := &stubRepository{
		getRequestByIDFunc: func(_ context.Context, requestID string) (VerificationRequest, error) {
			return VerificationRequest{
				ID:           requestID,
				BusinessRef:  "job-123",
				CandidateRef: "cand-456",
				Provider:     "walt",
				Status:       StatusVerified,
				Verified:     true,
				ReasonCode:   "POLICY_APPROVED",
			}, nil
		},
		getLatestSessionByRequestIDFunc: func(_ context.Context, requestID string) (VerificationSession, error) {
			return VerificationSession{
				ID:                    "sess-123",
				VerificationRequestID: requestID,
				Provider:              "walt",
				ProviderSessionID:     "provider-session-123",
				QRCodeURL:             "https://example.com/qr",
				DeepLink:              "openid://example",
			}, nil
		},
	}

	service := NewService("walt", repository, &stubSessionProvider{})

	response, err := service.GetRequest(context.Background(), "req-123")
	if err != nil {
		t.Fatalf("expected get request to succeed, got error: %v", err)
	}

	if !response.Verified {
		t.Fatal("expected request to be verified")
	}

	if response.ReasonCode != "POLICY_APPROVED" {
		t.Fatalf("expected reason code POLICY_APPROVED, got %q", response.ReasonCode)
	}
}

func TestServiceCreateRequestMapsProviderFailure(t *testing.T) {
	repository := &stubRepository{
		createRequestFunc: func(_ context.Context, request VerificationRequest) (VerificationRequest, error) {
			request.ID = "req-123"
			return request, nil
		},
		updateRequestStatusFunc: func(_ context.Context, request VerificationRequest) (VerificationRequest, error) {
			return request, nil
		},
	}
	provider := &stubSessionProvider{
		createSessionErr: errors.New("upstream create session failed"),
	}

	service := NewService("walt", repository, provider)

	_, err := service.CreateRequest(context.Background(), CreateRequestInput{
		BusinessRef:  "job-123",
		CandidateRef: "cand-456",
	})
	if err == nil {
		t.Fatal("expected provider failure to return an error")
	}

	var serviceError *Error
	if !errors.As(err, &serviceError) {
		t.Fatalf("expected typed verification error, got %T", err)
	}

	if serviceError.Code != ErrCodeProviderSessionCreateFailed {
		t.Fatalf("expected error code %q, got %q", ErrCodeProviderSessionCreateFailed, serviceError.Code)
	}

	if repository.updatedRequest.Status != StatusFailed {
		t.Fatalf("expected failed status to be persisted, got %q", repository.updatedRequest.Status)
	}

	if repository.updatedRequest.ReasonCode != ErrCodeProviderSessionCreateFailed {
		t.Fatalf("expected failure reason code %q, got %q", ErrCodeProviderSessionCreateFailed, repository.updatedRequest.ReasonCode)
	}
}

type stubRepository struct {
	createdRequest                  VerificationRequest
	updatedRequest                  VerificationRequest
	createdSession                  VerificationSession
	createRequestFunc               func(context.Context, VerificationRequest) (VerificationRequest, error)
	updateRequestStatusFunc         func(context.Context, VerificationRequest) (VerificationRequest, error)
	getRequestByIDFunc              func(context.Context, string) (VerificationRequest, error)
	createSessionFunc               func(context.Context, VerificationSession) (VerificationSession, error)
	getLatestSessionByRequestIDFunc func(context.Context, string) (VerificationSession, error)
}

func (repository *stubRepository) CreateRequest(ctx context.Context, request VerificationRequest) (VerificationRequest, error) {
	repository.createdRequest = request
	return repository.createRequestFunc(ctx, request)
}

func (repository *stubRepository) UpdateRequestStatus(ctx context.Context, request VerificationRequest) (VerificationRequest, error) {
	repository.updatedRequest = request
	return repository.updateRequestStatusFunc(ctx, request)
}

func (repository *stubRepository) GetRequestByID(ctx context.Context, requestID string) (VerificationRequest, error) {
	return repository.getRequestByIDFunc(ctx, requestID)
}

func (repository *stubRepository) CreateSession(ctx context.Context, session VerificationSession) (VerificationSession, error) {
	repository.createdSession = session
	return repository.createSessionFunc(ctx, session)
}

func (repository *stubRepository) GetLatestSessionByRequestID(ctx context.Context, requestID string) (VerificationSession, error) {
	return repository.getLatestSessionByRequestIDFunc(ctx, requestID)
}

type stubSessionProvider struct {
	createSessionResult ProviderSession
	createSessionErr    error
	createSessionInput  ProviderSessionInput
}

func (provider *stubSessionProvider) CreateSession(_ context.Context, input ProviderSessionInput) (ProviderSession, error) {
	provider.createSessionInput = input

	if provider.createSessionErr != nil {
		return ProviderSession{}, provider.createSessionErr
	}

	return provider.createSessionResult, nil
}
