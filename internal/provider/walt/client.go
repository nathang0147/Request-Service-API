package walt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	verifierBaseURL           string
	verifierMode              verifierMode
	bearerToken               string
	vcPolicyWebhookURL        string
	callbackBaseURL           string
	publicBaseURL             string
	publicRedirectURLTemplate string
	callbackAuthSecret        string
	httpClient                HTTPClient
}

func NewClient(verifierBaseURL, mode, bearerToken, vcPolicyWebhookURL, callbackBaseURL, publicBaseURL, publicRedirectURLTemplate, callbackAuthSecret string, httpClient HTTPClient) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	verifierMode := verifierMode(strings.ToLower(strings.TrimSpace(mode)))
	if verifierMode == "" {
		verifierMode = legacyVerifierMode
	}

	return &Client{
		verifierBaseURL:           strings.TrimRight(verifierBaseURL, "/"),
		verifierMode:              verifierMode,
		bearerToken:               bearerToken,
		vcPolicyWebhookURL:        strings.TrimSpace(vcPolicyWebhookURL),
		callbackBaseURL:           strings.TrimRight(strings.TrimSpace(callbackBaseURL), "/"),
		publicBaseURL:             strings.TrimRight(strings.TrimSpace(publicBaseURL), "/"),
		publicRedirectURLTemplate: strings.TrimSpace(publicRedirectURLTemplate),
		callbackAuthSecret:        strings.TrimSpace(callbackAuthSecret),
		httpClient:                httpClient,
	}
}

func (client *Client) CreateSession(ctx context.Context, request createSessionRequest) (createSessionResponse, []byte, error) {
	switch client.verifierMode {
	case verifier2VerifierMode:
		return client.createVerifier2Session(ctx, request)
	case legacyVerifierMode:
		return client.createLegacySession(ctx, request)
	default:
		return createSessionResponse{}, nil, fmt.Errorf("unsupported walt verifier mode %q", client.verifierMode)
	}
}

func (client *Client) createLegacySession(ctx context.Context, request createSessionRequest) (createSessionResponse, []byte, error) {
	if request.Legacy == nil {
		return createSessionResponse{}, nil, fmt.Errorf("missing legacy create session request")
	}

	body, err := json.Marshal(request.Legacy)
	if err != nil {
		return createSessionResponse{}, nil, err
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, client.verifierBaseURL+"/openid4vc/verify", bytes.NewReader(body))
	if err != nil {
		return createSessionResponse{}, nil, err
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	if stateID := deriveLegacyStateID(request.Legacy); stateID != "" {
		httpRequest.Header.Set("stateId", stateID)
	}
	if requestStatusURL := client.requestStatusURL(request.Legacy); requestStatusURL != "" {
		httpRequest.Header.Set("successRedirectUri", requestStatusURL)
		httpRequest.Header.Set("errorRedirectUri", requestStatusURL)
	}
	if callbackURL := client.statusCallbackURL(); callbackURL != "" {
		httpRequest.Header.Set("statusCallbackUri", callbackURL)
	}
	if client.callbackAuthSecret != "" {
		httpRequest.Header.Set("statusCallbackApiKey", client.callbackAuthSecret)
	}
	applyBearerToken(httpRequest, client.bearerToken)

	rawResponse, err := client.do(httpRequest)
	if err != nil {
		return createSessionResponse{}, nil, err
	}

	return createSessionResponse{LegacyURL: strings.TrimSpace(string(rawResponse))}, rawResponse, nil
}

func (client *Client) createVerifier2Session(ctx context.Context, request createSessionRequest) (createSessionResponse, []byte, error) {
	if request.Verifier2 == nil {
		return createSessionResponse{}, nil, fmt.Errorf("missing verifier2 create session request")
	}

	body, err := json.Marshal(request.Verifier2)
	if err != nil {
		return createSessionResponse{}, nil, err
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, client.verifierBaseURL+"/verification-session/create", bytes.NewReader(body))
	if err != nil {
		return createSessionResponse{}, nil, err
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	applyBearerToken(httpRequest, client.bearerToken)

	rawResponse, err := client.do(httpRequest)
	if err != nil {
		return createSessionResponse{}, nil, err
	}

	var response verifier2CreateSessionResponse
	if err := json.Unmarshal(rawResponse, &response); err != nil {
		return createSessionResponse{}, rawResponse, err
	}

	return createSessionResponse{Verifier2: &response}, rawResponse, nil
}

func (client *Client) do(httpRequest *http.Request) ([]byte, error) {
	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return nil, err
	}
	defer httpResponse.Body.Close()

	rawResponse, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return nil, err
	}

	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		return rawResponse, fmt.Errorf("walt create session returned status %d", httpResponse.StatusCode)
	}

	return rawResponse, nil
}

func deriveLegacyStateID(request *legacyCreateSessionRequest) string {
	if request == nil {
		return ""
	}

	return strings.TrimSpace(request.StateID)
}

func extractLegacyStateFromURL(rawURL string) string {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return ""
	}

	return parsedURL.Query().Get("state")
}

func (client *Client) statusCallbackURL() string {
	if client.callbackBaseURL == "" {
		return ""
	}

	return client.callbackBaseURL + "/api/v1/callbacks/walt"
}

func (client *Client) requestStatusURL(request *legacyCreateSessionRequest) string {
	if templateURL := strings.TrimSpace(client.publicRedirectURLTemplate); templateURL != "" {
		stateID := ""
		if request != nil {
			stateID = strings.TrimSpace(request.StateID)
		}
		return strings.ReplaceAll(templateURL, "$id", url.PathEscape(stateID))
	}

	baseURL := client.publicBaseURL
	if baseURL == "" {
		baseURL = client.callbackBaseURL
	}

	if baseURL == "" || request == nil {
		return ""
	}

	stateID := strings.TrimSpace(request.StateID)
	if stateID == "" {
		return ""
	}

	return baseURL + "/api/v1/verification-requests/" + url.PathEscape(stateID)
}
