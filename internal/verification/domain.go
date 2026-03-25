package verification

import "time"

type Status string

const (
	StatusCreated        Status = "CREATED"
	StatusSessionCreated Status = "SESSION_CREATED"
	StatusPending        Status = "PENDING"
	StatusVerified       Status = "VERIFIED"
	StatusFailed         Status = "FAILED"
	StatusExpired        Status = "EXPIRED"
)

type VerificationRequest struct {
	ID           string
	BusinessRef  string
	CandidateRef string
	Provider     string
	Status       Status
	Verified     bool
	ReasonCode   string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type VerificationSession struct {
	ID                    string
	VerificationRequestID string
	Provider              string
	ProviderSessionID     string
	QRCodeURL             string
	DeepLink              string
	OfferURL              string
	ExpiresAt             time.Time
	RawCreateResponse     []byte
	CreatedAt             time.Time
}

type ProviderSessionInput struct {
	RequestID    string
	BusinessRef  string
	CandidateRef string
	Provider     string
}

type ProviderSession struct {
	ProviderSessionID string
	QRCodeURL         string
	DeepLink          string
	OfferURL          string
	ExpiresAt         time.Time
	RawCreateResponse []byte
}
