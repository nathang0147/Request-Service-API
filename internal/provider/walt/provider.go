package walt

import (
	"context"
	"encoding/json"

	"github.com/nathang0147/Request-Service-API/internal/provider"
	"github.com/nathang0147/Request-Service-API/internal/verification"
)

type Provider struct {
	client *Client
}

func New(client *Client) *Provider {
	return &Provider{client: client}
}

func (adapter *Provider) CreateSession(ctx context.Context, input verification.ProviderSessionInput) (verification.ProviderSession, error) {
	response, rawResponse, err := adapter.client.CreateSession(ctx, toCreateSessionRequest(input, adapter.client.vcPolicyWebhookURL))
	if err != nil {
		return verification.ProviderSession{}, err
	}

	return toProviderSession(response, rawResponse), nil
}

func (adapter *Provider) ParseCallback(_ context.Context, body []byte) (provider.CallbackEvent, error) {
	var legacyPayload callbackPayload
	if err := json.Unmarshal(body, &legacyPayload); err == nil && legacyPayload.SessionID != "" {
		return toCallbackEvent(legacyPayload, body), nil
	}

	var statusPayload statusCallbackPayload
	if err := json.Unmarshal(body, &statusPayload); err != nil {
		return provider.CallbackEvent{}, err
	}

	return toStatusCallbackEvent(statusPayload, body), nil
}
