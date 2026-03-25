package postgres

import (
	"context"

	"github.com/nathang0147/Request-Service-API/internal/callback"
	"github.com/nathang0147/Request-Service-API/internal/persistence/postgres/sqlc"
)

type CallbackStore struct {
	repository *CallbackRepository
}

func NewCallbackStore(db sqlc.DBTX) *CallbackStore {
	return &CallbackStore{
		repository: NewCallbackRepository(db),
	}
}

func (store *CallbackStore) GetSessionByProviderSessionID(ctx context.Context, providerSessionID string) (callback.SessionRecord, error) {
	session, err := store.repository.GetSessionByProviderSessionID(ctx, providerSessionID)
	if err != nil {
		return callback.SessionRecord{}, err
	}

	return callback.SessionRecord{
		VerificationRequestID: session.VerificationRequestID.String(),
		ProviderSessionID:     session.ProviderSessionID,
		Provider:              session.Provider,
	}, nil
}

func (store *CallbackStore) UpdateRequestStatus(ctx context.Context, update callback.RequestStatusUpdate) (callback.RequestStatusUpdate, error) {
	updated, err := store.repository.UpdateRequestStatus(ctx, toCallbackStatusParams(update))
	if err != nil {
		return callback.RequestStatusUpdate{}, err
	}

	return callback.RequestStatusUpdate{
		RequestID:  updated.ID.String(),
		Status:     callbackStatus(updated.Status),
		Verified:   updated.Verified,
		ReasonCode: updated.ReasonCode.String,
	}, nil
}

func (store *CallbackStore) CreateEvent(ctx context.Context, event callback.EventRecord) (callback.EventRecord, error) {
	created, err := store.repository.CreateEvent(ctx, toCreateEventParams(event))
	if err != nil {
		return callback.EventRecord{}, err
	}

	return callback.EventRecord{
		RequestID: created.VerificationRequestID.String(),
		Source:    created.Source,
		EventType: created.EventType,
		Payload:   created.Payload,
	}, nil
}

var _ callback.Repository = (*CallbackStore)(nil)
