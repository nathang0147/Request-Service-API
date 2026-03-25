package postgres

import (
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/nathang0147/Request-Service-API/internal/callback"
	"github.com/nathang0147/Request-Service-API/internal/persistence/postgres/sqlc"
	"github.com/nathang0147/Request-Service-API/internal/verification"
)

func parseUUID(value string) (pgtype.UUID, error) {
	var id pgtype.UUID
	if err := id.Scan(value); err != nil {
		return pgtype.UUID{}, err
	}

	return id, nil
}

func toCreateRequestParams(request verification.VerificationRequest) sqlc.CreateVerificationRequestParams {
	return sqlc.CreateVerificationRequestParams{
		BusinessRef:  request.BusinessRef,
		CandidateRef: request.CandidateRef,
		Provider:     request.Provider,
		Status:       string(request.Status),
		Verified:     request.Verified,
		ReasonCode:   textValue(request.ReasonCode),
	}
}

func toUpdateRequestStatusParams(request verification.VerificationRequest) sqlc.UpdateVerificationRequestStatusParams {
	return sqlc.UpdateVerificationRequestStatusParams{
		ID:         mustUUID(request.ID),
		Status:     string(request.Status),
		Verified:   request.Verified,
		ReasonCode: textValue(request.ReasonCode),
	}
}

func toCreateSessionParams(session verification.VerificationSession) sqlc.CreateVerificationSessionParams {
	return sqlc.CreateVerificationSessionParams{
		VerificationRequestID: mustUUID(session.VerificationRequestID),
		Provider:              session.Provider,
		ProviderSessionID:     session.ProviderSessionID,
		QrCodeUrl:             textValue(session.QRCodeURL),
		DeepLink:              textValue(session.DeepLink),
		OfferUrl:              textValue(session.OfferURL),
		ExpiresAt:             timestampValue(session.ExpiresAt),
		RawCreateResponse:     session.RawCreateResponse,
	}
}

func toCreateEventParams(event callback.EventRecord) sqlc.CreateVerificationEventParams {
	return sqlc.CreateVerificationEventParams{
		VerificationRequestID: mustUUID(event.RequestID),
		Source:                event.Source,
		EventType:             event.EventType,
		Payload:               event.Payload,
	}
}

func toCallbackStatusParams(update callback.RequestStatusUpdate) sqlc.UpdateVerificationRequestStatusParams {
	return sqlc.UpdateVerificationRequestStatusParams{
		ID:         mustUUID(update.RequestID),
		Status:     string(update.Status),
		Verified:   update.Verified,
		ReasonCode: textValue(update.ReasonCode),
	}
}

func toDomainRequest(request sqlc.VerificationRequest) verification.VerificationRequest {
	return verification.VerificationRequest{
		ID:           request.ID.String(),
		BusinessRef:  request.BusinessRef,
		CandidateRef: request.CandidateRef,
		Provider:     request.Provider,
		Status:       verification.Status(request.Status),
		Verified:     request.Verified,
		ReasonCode:   request.ReasonCode.String,
		CreatedAt:    request.CreatedAt.Time,
		UpdatedAt:    request.UpdatedAt.Time,
	}
}

func toDomainSession(session sqlc.VerificationSession) verification.VerificationSession {
	return verification.VerificationSession{
		ID:                    session.ID.String(),
		VerificationRequestID: session.VerificationRequestID.String(),
		Provider:              session.Provider,
		ProviderSessionID:     session.ProviderSessionID,
		QRCodeURL:             session.QrCodeUrl.String,
		DeepLink:              session.DeepLink.String,
		OfferURL:              session.OfferUrl.String,
		ExpiresAt:             session.ExpiresAt.Time,
		RawCreateResponse:     session.RawCreateResponse,
		CreatedAt:             session.CreatedAt.Time,
	}
}

func callbackStatus(status string) verification.Status {
	return verification.Status(status)
}

func textValue(value string) pgtype.Text {
	return pgtype.Text{String: value, Valid: value != ""}
}

func timestampValue(value time.Time) pgtype.Timestamptz {
	if value.IsZero() {
		return pgtype.Timestamptz{}
	}

	return pgtype.Timestamptz{
		Time:  value,
		Valid: true,
	}
}

func mustUUID(value string) pgtype.UUID {
	id, err := parseUUID(value)
	if err != nil {
		panic(fmt.Sprintf("invalid uuid %q: %v", value, err))
	}

	return id
}
