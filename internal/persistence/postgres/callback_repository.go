package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nathang0147/Request-Service-API/internal/persistence/postgres/sqlc"
)

type CallbackRepository struct {
	queries *sqlc.Queries
}

func NewCallbackRepository(db sqlc.DBTX) *CallbackRepository {
	return &CallbackRepository{
		queries: sqlc.New(db),
	}
}

func (repository *CallbackRepository) GetSessionByProviderSessionID(ctx context.Context, providerSessionID string) (sqlc.VerificationSession, error) {
	return repository.queries.GetVerificationSessionByProviderSessionID(ctx, providerSessionID)
}

func (repository *CallbackRepository) CreateEvent(ctx context.Context, params sqlc.CreateVerificationEventParams) (sqlc.VerificationEvent, error) {
	return repository.queries.CreateVerificationEvent(ctx, params)
}

func (repository *CallbackRepository) ListEventsByRequestID(ctx context.Context, requestID pgtype.UUID) ([]sqlc.VerificationEvent, error) {
	return repository.queries.ListVerificationEventsByRequestID(ctx, requestID)
}

func (repository *CallbackRepository) UpdateRequestStatus(ctx context.Context, params sqlc.UpdateVerificationRequestStatusParams) (sqlc.VerificationRequest, error) {
	return repository.queries.UpdateVerificationRequestStatus(ctx, params)
}
