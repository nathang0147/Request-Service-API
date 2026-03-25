package postgres

import (
	"context"
	"testing"
	"time"

	pgxmock "github.com/pashagolub/pgxmock/v4"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/nathang0147/Request-Service-API/internal/persistence/postgres/sqlc"
)

func TestVerificationRepositoryCreateAndGetRequest(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("expected pgxmock connection, got error: %v", err)
	}
	defer mock.Close(context.Background())

	repo := NewVerificationRepository(mock)
	requestID := testUUID(1)
	createdAt := testTimestamp(2026, time.March, 25, 8, 0, 0)
	updatedAt := testTimestamp(2026, time.March, 25, 8, 1, 0)
	reasonCode := testTextValue("")

	createArgs := sqlc.CreateVerificationRequestParams{
		BusinessRef:  "job-123",
		CandidateRef: "cand-456",
		Provider:     "walt",
		Status:       "CREATED",
		Verified:     false,
		ReasonCode:   reasonCode,
	}

	mock.ExpectQuery("INSERT INTO verification_requests").
		WithArgs(
			createArgs.BusinessRef,
			createArgs.CandidateRef,
			createArgs.Provider,
			createArgs.Status,
			createArgs.Verified,
			createArgs.ReasonCode,
		).
		WillReturnRows(requestRows().
			AddRow(
				requestID,
				createArgs.BusinessRef,
				createArgs.CandidateRef,
				createArgs.Provider,
				createArgs.Status,
				createArgs.Verified,
				reasonCode,
				createdAt,
				updatedAt,
			))

	created, err := repo.CreateRequest(context.Background(), createArgs)
	if err != nil {
		t.Fatalf("expected request creation to succeed, got error: %v", err)
	}

	if created.Status != "CREATED" {
		t.Fatalf("expected status CREATED, got %q", created.Status)
	}

	mock.ExpectQuery("SELECT id, business_ref, candidate_ref, provider, status, verified, reason_code, created_at, updated_at FROM verification_requests").
		WithArgs(requestID).
		WillReturnRows(requestRows().
			AddRow(
				requestID,
				createArgs.BusinessRef,
				createArgs.CandidateRef,
				createArgs.Provider,
				createArgs.Status,
				createArgs.Verified,
				reasonCode,
				createdAt,
				updatedAt,
			))

	fetched, err := repo.GetRequestByID(context.Background(), requestID)
	if err != nil {
		t.Fatalf("expected request lookup to succeed, got error: %v", err)
	}

	if fetched.BusinessRef != "job-123" {
		t.Fatalf("expected business ref job-123, got %q", fetched.BusinessRef)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected mock expectations to be met, got error: %v", err)
	}
}

func TestVerificationRepositoryCreateSessionAndLookupByProviderSessionID(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("expected pgxmock connection, got error: %v", err)
	}
	defer mock.Close(context.Background())

	repo := NewVerificationRepository(mock)
	requestID := testUUID(2)
	sessionID := testUUID(3)
	expiresAt := testTimestamp(2026, time.March, 25, 10, 0, 0)
	createdAt := testTimestamp(2026, time.March, 25, 8, 2, 0)

	createArgs := sqlc.CreateVerificationSessionParams{
		VerificationRequestID: requestID,
		Provider:              "walt",
		ProviderSessionID:     "provider-session-123",
		QrCodeUrl:             testTextValue("https://example.com/qr"),
		DeepLink:              testTextValue("openid://example"),
		OfferUrl:              testTextValue("https://example.com/offer"),
		ExpiresAt:             expiresAt,
		RawCreateResponse:     []byte(`{"session":"ok"}`),
	}

	mock.ExpectQuery("INSERT INTO verification_sessions").
		WithArgs(
			createArgs.VerificationRequestID,
			createArgs.Provider,
			createArgs.ProviderSessionID,
			createArgs.QrCodeUrl,
			createArgs.DeepLink,
			createArgs.OfferUrl,
			createArgs.ExpiresAt,
			createArgs.RawCreateResponse,
		).
		WillReturnRows(sessionRows().
			AddRow(
				sessionID,
				requestID,
				createArgs.Provider,
				createArgs.ProviderSessionID,
				createArgs.QrCodeUrl,
				createArgs.DeepLink,
				createArgs.OfferUrl,
				createArgs.ExpiresAt,
				createArgs.RawCreateResponse,
				createdAt,
			))

	created, err := repo.CreateSession(context.Background(), createArgs)
	if err != nil {
		t.Fatalf("expected session creation to succeed, got error: %v", err)
	}

	if created.ProviderSessionID != "provider-session-123" {
		t.Fatalf("expected provider session id provider-session-123, got %q", created.ProviderSessionID)
	}

	mock.ExpectQuery("SELECT id, verification_request_id, provider, provider_session_id, qr_code_url, deep_link, offer_url, expires_at, raw_create_response, created_at FROM verification_sessions").
		WithArgs("provider-session-123").
		WillReturnRows(sessionRows().
			AddRow(
				sessionID,
				requestID,
				createArgs.Provider,
				createArgs.ProviderSessionID,
				createArgs.QrCodeUrl,
				createArgs.DeepLink,
				createArgs.OfferUrl,
				createArgs.ExpiresAt,
				createArgs.RawCreateResponse,
				createdAt,
			))

	fetched, err := repo.GetSessionByProviderSessionID(context.Background(), "provider-session-123")
	if err != nil {
		t.Fatalf("expected session lookup to succeed, got error: %v", err)
	}

	if fetched.Provider != "walt" {
		t.Fatalf("expected provider walt, got %q", fetched.Provider)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected mock expectations to be met, got error: %v", err)
	}
}

func requestRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"business_ref",
		"candidate_ref",
		"provider",
		"status",
		"verified",
		"reason_code",
		"created_at",
		"updated_at",
	})
}

func sessionRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"verification_request_id",
		"provider",
		"provider_session_id",
		"qr_code_url",
		"deep_link",
		"offer_url",
		"expires_at",
		"raw_create_response",
		"created_at",
	})
}

func testTextValue(value string) pgtype.Text {
	return pgtype.Text{String: value, Valid: value != ""}
}

func testUUID(seed byte) pgtype.UUID {
	var bytes [16]byte
	bytes[15] = seed

	return pgtype.UUID{Bytes: bytes, Valid: true}
}

func testTimestamp(year int, month time.Month, day, hour, minute, second int) pgtype.Timestamptz {
	return pgtype.Timestamptz{
		Time:  time.Date(year, month, day, hour, minute, second, 0, time.UTC),
		Valid: true,
	}
}
