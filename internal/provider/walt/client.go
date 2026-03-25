package walt

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	baseURL    string
	apiKey     string
	httpClient HTTPClient
}

func NewClient(baseURL, apiKey string, httpClient HTTPClient) *Client {
	if httpClient == nil {
		httpClient = http.DefaultClient
	}

	return &Client{
		baseURL:    strings.TrimRight(baseURL, "/"),
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

func (client *Client) CreateSession(ctx context.Context, request createSessionRequest) (createSessionResponse, []byte, error) {
	body, err := json.Marshal(request)
	if err != nil {
		return createSessionResponse{}, nil, err
	}

	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, client.baseURL+"/sessions", bytes.NewReader(body))
	if err != nil {
		return createSessionResponse{}, nil, err
	}

	httpRequest.Header.Set("Content-Type", "application/json")
	applyAPIKey(httpRequest, client.apiKey)

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return createSessionResponse{}, nil, err
	}
	defer httpResponse.Body.Close()

	rawResponse, err := io.ReadAll(httpResponse.Body)
	if err != nil {
		return createSessionResponse{}, nil, err
	}

	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		return createSessionResponse{}, rawResponse, fmt.Errorf("walt create session returned status %d", httpResponse.StatusCode)
	}

	var response createSessionResponse
	if err := json.Unmarshal(rawResponse, &response); err != nil {
		return createSessionResponse{}, rawResponse, err
	}

	return response, rawResponse, nil
}
