package verification

import "context"

type Repository interface {
	CreateRequest(context.Context, VerificationRequest) (VerificationRequest, error)
	UpdateRequestStatus(context.Context, VerificationRequest) (VerificationRequest, error)
	GetRequestByID(context.Context, string) (VerificationRequest, error)
	CreateSession(context.Context, VerificationSession) (VerificationSession, error)
	GetLatestSessionByRequestID(context.Context, string) (VerificationSession, error)
}

type SessionProvider interface {
	CreateSession(context.Context, ProviderSessionInput) (ProviderSession, error)
}

type ProviderResolver interface {
	Resolve(context.Context, string) (SessionProvider, error)
}
