package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nathang0147/Request-Service-API/internal/persistence/postgres/sqlc"
)

type VerificationRepository struct {
	queries *sqlc.Queries
}

func NewVerificationRepository(db sqlc.DBTX) *VerificationRepository {
	return &VerificationRepository{
		queries: sqlc.New(db),
	}
}

func (repository *VerificationRepository) CreateRequest(ctx context.Context, params sqlc.CreateVerificationRequestParams) (sqlc.VerificationRequest, error) {
	return repository.queries.CreateVerificationRequest(ctx, params)
}

func (repository *VerificationRepository) GetRequestByID(ctx context.Context, id pgtype.UUID) (sqlc.VerificationRequest, error) {
	return repository.queries.GetVerificationRequestByID(ctx, id)
}

func (repository *VerificationRepository) UpdateRequestStatus(ctx context.Context, params sqlc.UpdateVerificationRequestStatusParams) (sqlc.VerificationRequest, error) {
	return repository.queries.UpdateVerificationRequestStatus(ctx, params)
}

func (repository *VerificationRepository) CreateSession(ctx context.Context, params sqlc.CreateVerificationSessionParams) (sqlc.VerificationSession, error) {
	return repository.queries.CreateVerificationSession(ctx, params)
}

func (repository *VerificationRepository) GetSessionByProviderSessionID(ctx context.Context, providerSessionID string) (sqlc.VerificationSession, error) {
	return repository.queries.GetVerificationSessionByProviderSessionID(ctx, providerSessionID)
}

func (repository *VerificationRepository) GetLatestSessionByRequestID(ctx context.Context, requestID pgtype.UUID) (sqlc.VerificationSession, error) {
	return repository.queries.GetLatestVerificationSessionByRequestID(ctx, requestID)
}
