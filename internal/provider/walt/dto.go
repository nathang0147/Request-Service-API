package walt

import (
	"encoding/json"
	"time"
)

type createSessionRequest struct {
	RequestID    string `json:"requestId"`
	BusinessRef  string `json:"businessRef"`
	CandidateRef string `json:"candidateRef"`
}

type createSessionResponse struct {
	SessionID string    `json:"sessionId"`
	QRCodeURL string    `json:"qrCodeUrl"`
	DeepLink  string    `json:"deepLink"`
	OfferURL  string    `json:"offerUrl"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type callbackPayload struct {
	SessionID  string          `json:"sessionId"`
	Status     string          `json:"status"`
	Verified   bool            `json:"verified"`
	ReasonCode string          `json:"reasonCode"`
	EventType  string          `json:"eventType"`
	Payload    json.RawMessage `json:"payload"`
}
