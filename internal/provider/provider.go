package provider

import (
	"context"

	"github.com/nathang0147/Request-Service-API/internal/verification"
)

type CallbackEvent struct {
	ProviderSessionID string
	Status            verification.Status
	Verified          bool
	ReasonCode        string
	EventType         string
	Payload           []byte
}

type Provider interface {
	CreateSession(context.Context, verification.ProviderSessionInput) (verification.ProviderSession, error)
	ParseCallback(context.Context, []byte) (CallbackEvent, error)
}
