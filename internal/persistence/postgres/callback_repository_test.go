package postgres

import (
	"context"
	"testing"
	"time"

	pgxmock "github.com/pashagolub/pgxmock/v4"

	"github.com/nathang0147/Request-Service-API/internal/persistence/postgres/sqlc"
)

func TestCallbackRepositoryCreateEvent(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("expected pgxmock connection, got error: %v", err)
	}
	defer mock.Close(context.Background())

	repo := NewCallbackRepository(mock)
	requestID := testUUID(4)
	eventID := testUUID(5)
	createdAt := testTimestamp(2026, time.March, 25, 9, 0, 0)

	createArgs := sqlc.CreateVerificationEventParams{
		VerificationRequestID: requestID,
		Source:                "walt_callback",
		EventType:             "SESSION_VERIFIED",
		Payload:               []byte(`{"verified":true}`),
	}

	mock.ExpectQuery("INSERT INTO verification_events").
		WithArgs(
			createArgs.VerificationRequestID,
			createArgs.Source,
			createArgs.EventType,
			createArgs.Payload,
		).
		WillReturnRows(eventRows().
			AddRow(
				eventID,
				requestID,
				createArgs.Source,
				createArgs.EventType,
				createArgs.Payload,
				createdAt,
			))

	created, err := repo.CreateEvent(context.Background(), createArgs)
	if err != nil {
		t.Fatalf("expected event creation to succeed, got error: %v", err)
	}

	if created.EventType != "SESSION_VERIFIED" {
		t.Fatalf("expected event type SESSION_VERIFIED, got %q", created.EventType)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected mock expectations to be met, got error: %v", err)
	}
}

func TestCallbackRepositoryUpdateRequestStatus(t *testing.T) {
	mock, err := pgxmock.NewConn()
	if err != nil {
		t.Fatalf("expected pgxmock connection, got error: %v", err)
	}
	defer mock.Close(context.Background())

	repo := NewCallbackRepository(mock)
	requestID := testUUID(6)
	createdAt := testTimestamp(2026, time.March, 25, 7, 0, 0)
	updatedAt := testTimestamp(2026, time.March, 25, 9, 5, 0)
	reasonCode := textValue("POLICY_APPROVED")

	updateArgs := sqlc.UpdateVerificationRequestStatusParams{
		ID:         requestID,
		Status:     "VERIFIED",
		Verified:   true,
		ReasonCode: reasonCode,
	}

	mock.ExpectQuery("UPDATE verification_requests").
		WithArgs(
			updateArgs.ID,
			updateArgs.Status,
			updateArgs.Verified,
			updateArgs.ReasonCode,
		).
		WillReturnRows(requestRows().
			AddRow(
				requestID,
				"job-123",
				"cand-456",
				"walt",
				updateArgs.Status,
				updateArgs.Verified,
				reasonCode,
				createdAt,
				updatedAt,
			))

	updated, err := repo.UpdateRequestStatus(context.Background(), updateArgs)
	if err != nil {
		t.Fatalf("expected status update to succeed, got error: %v", err)
	}

	if !updated.Verified {
		t.Fatal("expected request to be verified after update")
	}

	if updated.ReasonCode.String != "POLICY_APPROVED" {
		t.Fatalf("expected reason code POLICY_APPROVED, got %q", updated.ReasonCode.String)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("expected mock expectations to be met, got error: %v", err)
	}
}

func eventRows() *pgxmock.Rows {
	return pgxmock.NewRows([]string{
		"id",
		"verification_request_id",
		"source",
		"event_type",
		"payload",
		"created_at",
	})
}
