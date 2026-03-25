package provider

import (
	"context"
	"testing"

	"github.com/nathang0147/Request-Service-API/internal/verification"
)

func TestResolverAlwaysReturnsWaltInV1(t *testing.T) {
	waltProvider := &stubProvider{}
	resolver := NewResolver(waltProvider)

	resolved, err := resolver.Resolve(context.Background(), "anything")
	if err != nil {
		t.Fatalf("expected resolver to succeed, got error: %v", err)
	}

	if resolved != waltProvider {
		t.Fatal("expected resolver to always return the configured walt provider")
	}
}

type stubProvider struct{}

func (provider *stubProvider) CreateSession(context.Context, verification.ProviderSessionInput) (verification.ProviderSession, error) {
	return verification.ProviderSession{}, nil
}

func (provider *stubProvider) ParseCallback(context.Context, []byte) (CallbackEvent, error) {
	return CallbackEvent{}, nil
}
