package provider

import (
	"context"

	"github.com/nathang0147/Request-Service-API/internal/verification"
)

type Resolver struct {
	walt Provider
}

func NewResolver(walt Provider) *Resolver {
	return &Resolver{walt: walt}
}

func (resolver *Resolver) Resolve(_ context.Context, _ string) (verification.SessionProvider, error) {
	return resolver.walt, nil
}
