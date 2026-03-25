package postgres

import (
	"context"

	"github.com/nathang0147/Request-Service-API/internal/persistence/postgres/sqlc"
	"github.com/nathang0147/Request-Service-API/internal/verification"
)

type VerificationStore struct {
	repository *VerificationRepository
}

func NewVerificationStore(db sqlc.DBTX) *VerificationStore {
	return &VerificationStore{
		repository: NewVerificationRepository(db),
	}
}

func (store *VerificationStore) CreateRequest(ctx context.Context, request verification.VerificationRequest) (verification.VerificationRequest, error) {
	created, err := store.repository.CreateRequest(ctx, toCreateRequestParams(request))
	if err != nil {
		return verification.VerificationRequest{}, err
	}

	return toDomainRequest(created), nil
}

func (store *VerificationStore) UpdateRequestStatus(ctx context.Context, request verification.VerificationRequest) (verification.VerificationRequest, error) {
	updated, err := store.repository.UpdateRequestStatus(ctx, toUpdateRequestStatusParams(request))
	if err != nil {
		return verification.VerificationRequest{}, err
	}

	return toDomainRequest(updated), nil
}

func (store *VerificationStore) GetRequestByID(ctx context.Context, requestID string) (verification.VerificationRequest, error) {
	id, err := parseUUID(requestID)
	if err != nil {
		return verification.VerificationRequest{}, err
	}

	request, err := store.repository.GetRequestByID(ctx, id)
	if err != nil {
		return verification.VerificationRequest{}, err
	}

	return toDomainRequest(request), nil
}

func (store *VerificationStore) CreateSession(ctx context.Context, session verification.VerificationSession) (verification.VerificationSession, error) {
	created, err := store.repository.CreateSession(ctx, toCreateSessionParams(session))
	if err != nil {
		return verification.VerificationSession{}, err
	}

	return toDomainSession(created), nil
}

func (store *VerificationStore) GetLatestSessionByRequestID(ctx context.Context, requestID string) (verification.VerificationSession, error) {
	id, err := parseUUID(requestID)
	if err != nil {
		return verification.VerificationSession{}, err
	}

	session, err := store.repository.GetLatestSessionByRequestID(ctx, id)
	if err != nil {
		return verification.VerificationSession{}, err
	}

	return toDomainSession(session), nil
}

var _ verification.Repository = (*VerificationStore)(nil)
