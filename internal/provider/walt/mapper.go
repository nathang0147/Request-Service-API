package walt

import (
	"github.com/nathang0147/Request-Service-API/internal/provider"
	"github.com/nathang0147/Request-Service-API/internal/verification"
)

func toCreateSessionRequest(input verification.ProviderSessionInput) createSessionRequest {
	return createSessionRequest{
		RequestID:    input.RequestID,
		BusinessRef:  input.BusinessRef,
		CandidateRef: input.CandidateRef,
	}
}

func toProviderSession(response createSessionResponse, raw []byte) verification.ProviderSession {
	return verification.ProviderSession{
		ProviderSessionID: response.SessionID,
		QRCodeURL:         response.QRCodeURL,
		DeepLink:          response.DeepLink,
		OfferURL:          response.OfferURL,
		ExpiresAt:         response.ExpiresAt,
		RawCreateResponse: raw,
	}
}

func toCallbackEvent(payload callbackPayload, raw []byte) provider.CallbackEvent {
	status := verification.StatusPending
	switch payload.Status {
	case string(verification.StatusVerified):
		status = verification.StatusVerified
	case string(verification.StatusFailed):
		status = verification.StatusFailed
	case string(verification.StatusExpired):
		status = verification.StatusExpired
	}

	return provider.CallbackEvent{
		ProviderSessionID: payload.SessionID,
		Status:            status,
		Verified:          payload.Verified,
		ReasonCode:        payload.ReasonCode,
		EventType:         payload.EventType,
		Payload:           raw,
	}
}
