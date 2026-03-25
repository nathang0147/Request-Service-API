package verification

import "time"

type CreateRequestInput struct {
	BusinessRef  string `json:"businessRef"`
	CandidateRef string `json:"candidateRef"`
}

type SessionDetails struct {
	QRCodeURL string    `json:"qrCodeUrl"`
	DeepLink  string    `json:"deepLink"`
	OfferURL  string    `json:"offerUrl"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type CreateRequestResponse struct {
	RequestID  string         `json:"requestId"`
	Status     Status         `json:"status"`
	Verified   bool           `json:"verified"`
	ReasonCode string         `json:"reasonCode,omitempty"`
	Session    SessionDetails `json:"session"`
}

type GetRequestResponse struct {
	RequestID  string          `json:"requestId"`
	Status     Status          `json:"status"`
	Verified   bool            `json:"verified"`
	ReasonCode string          `json:"reasonCode,omitempty"`
	Session    *SessionDetails `json:"session,omitempty"`
}
