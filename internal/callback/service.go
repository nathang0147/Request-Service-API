package callback

import (
	"context"

	"github.com/nathang0147/Request-Service-API/internal/provider"
	"github.com/nathang0147/Request-Service-API/internal/verification"
)

type Repository interface {
	GetSessionByProviderSessionID(context.Context, string) (SessionRecord, error)
	UpdateRequestStatus(context.Context, RequestStatusUpdate) (RequestStatusUpdate, error)
	CreateEvent(context.Context, EventRecord) (EventRecord, error)
}

type SessionRecord struct {
	VerificationRequestID string
	ProviderSessionID     string
	Provider              string
}

type RequestStatusUpdate struct {
	RequestID  string
	Status     verification.Status
	Verified   bool
	ReasonCode string
}

type EventRecord struct {
	RequestID string
	Source    string
	EventType string
	Payload   []byte
}

type EventParser interface {
	ParseCallback(context.Context, []byte) (provider.CallbackEvent, error)
}

type Service struct {
	repository Repository
	parser     EventParser
}

func NewService(repository Repository, parser EventParser) *Service {
	return &Service{
		repository: repository,
		parser:     parser,
	}
}

func (service *Service) HandleCallback(ctx context.Context, body []byte) (Result, error) {
	event, err := service.parser.ParseCallback(ctx, body)
	if err != nil {
		return Result{}, err
	}

	session, err := service.repository.GetSessionByProviderSessionID(ctx, event.ProviderSessionID)
	if err != nil {
		return Result{}, err
	}

	if _, err := service.repository.UpdateRequestStatus(ctx, RequestStatusUpdate{
		RequestID:  session.VerificationRequestID,
		Status:     event.Status,
		Verified:   event.Verified,
		ReasonCode: event.ReasonCode,
	}); err != nil {
		return Result{}, err
	}

	if _, err := service.repository.CreateEvent(ctx, EventRecord{
		RequestID: session.VerificationRequestID,
		Source:    "walt_callback",
		EventType: event.EventType,
		Payload:   event.Payload,
	}); err != nil {
		return Result{}, err
	}

	return toResult(), nil
}
