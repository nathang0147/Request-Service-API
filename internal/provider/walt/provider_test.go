package walt

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/nathang0147/Request-Service-API/internal/verification"
)

func TestProviderCreateSessionMapsRequestAndResponse(t *testing.T) {
	var capturedRequest legacyCreateSessionRequest
	httpClient := roundTripClient(func(request *http.Request) (*http.Response, error) {
		if request.Method != http.MethodPost {
			t.Fatalf("expected method %q, got %q", http.MethodPost, request.Method)
		}

		if request.URL.String() != "http://host.docker.internal:7003/openid4vc/verify" {
			t.Fatalf("expected url %q, got %q", "http://host.docker.internal:7003/openid4vc/verify", request.URL.String())
		}

		if request.Header.Get("Authorization") != "Bearer bearer-token" {
			t.Fatalf("expected bearer auth header, got %q", request.Header.Get("Authorization"))
		}

		if request.Header.Get("stateId") != "req-123" {
			t.Fatalf("expected stateId header req-123, got %q", request.Header.Get("stateId"))
		}

		if request.Header.Get("statusCallbackUri") != "https://callback.request-service.example.com/api/v1/callbacks/walt" {
			t.Fatalf("expected statusCallbackUri header, got %q", request.Header.Get("statusCallbackUri"))
		}

		if request.Header.Get("statusCallbackApiKey") != "callback-secret" {
			t.Fatalf("expected statusCallbackApiKey header, got %q", request.Header.Get("statusCallbackApiKey"))
		}

		expectedStatusURL := "https://wallet.example.com/success/req-123"
		if request.Header.Get("successRedirectUri") != expectedStatusURL {
			t.Fatalf("expected successRedirectUri header %q, got %q", expectedStatusURL, request.Header.Get("successRedirectUri"))
		}

		if request.Header.Get("errorRedirectUri") != expectedStatusURL {
			t.Fatalf("expected errorRedirectUri header %q, got %q", expectedStatusURL, request.Header.Get("errorRedirectUri"))
		}

		if err := json.NewDecoder(request.Body).Decode(&capturedRequest); err != nil {
			t.Fatalf("expected valid walt request body, got error: %v", err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Content-Type": []string{"text/plain; charset=UTF-8"},
			},
			Body: io.NopCloser(strings.NewReader("openid4vp://authorize?response_type=vp_token&state=req-123&presentation_definition_uri=https%3A%2F%2Fverifier.example.com%2Fopenid4vc%2Fpd%2Freq-123")),
		}, nil
	})

	client := NewClient(
		"http://host.docker.internal:7003",
		"legacy",
		"bearer-token",
		"http://host.docker.internal:8787/api/verifier/policies/vc",
		"https://callback.request-service.example.com",
		"https://public.request-service.example.com",
		"https://wallet.example.com/success/$id",
		"callback-secret",
		httpClient,
	)
	provider := New(client)

	session, err := provider.CreateSession(context.Background(), verification.ProviderSessionInput{
		RequestID:    "req-123",
		BusinessRef:  "job-123",
		CandidateRef: "cand-456",
		Provider:     "walt",
	})
	if err != nil {
		t.Fatalf("expected create session to succeed, got error: %v", err)
	}

	if len(capturedRequest.RequestCredentials) != 1 {
		t.Fatalf("expected one request credential, got %+v", capturedRequest.RequestCredentials)
	}

	if capturedRequest.RequestCredentials[0].Type != "UniversityDegree" {
		t.Fatalf("expected UniversityDegree type, got %+v", capturedRequest.RequestCredentials)
	}

	if capturedRequest.RequestCredentials[0].Format != "jwt_vc_json" {
		t.Fatalf("expected jwt_vc_json format, got %+v", capturedRequest.RequestCredentials)
	}

	inputDescriptor := capturedRequest.RequestCredentials[0].InputDescriptor
	if inputDescriptor.ID != "UniversityDegree" {
		t.Fatalf("expected input descriptor id UniversityDegree, got %+v", inputDescriptor)
	}

	if inputDescriptor.Format["jwt_vc_json"].Alg[0] != "EdDSA" {
		t.Fatalf("expected jwt_vc_json alg EdDSA, got %+v", inputDescriptor.Format)
	}

	if inputDescriptor.Constraints.Fields[0].Filter.Type != "array" {
		t.Fatalf("expected array filter for vc.type, got %+v", inputDescriptor.Constraints.Fields[0].Filter)
	}

	if inputDescriptor.Constraints.Fields[0].Filter.Contains == nil || inputDescriptor.Constraints.Fields[0].Filter.Contains.Const != "UniversityDegree" {
		t.Fatalf("expected contains.const UniversityDegree, got %+v", inputDescriptor.Constraints.Fields[0].Filter)
	}

	if len(capturedRequest.VPPolicies) != 2 {
		t.Fatalf("expected two vp policies, got %+v", capturedRequest.VPPolicies)
	}

	if got := capturedRequest.VPPolicies[0]; got != "signature" {
		t.Fatalf("expected first vp policy signature, got %#v", got)
	}

	if got := capturedRequest.VPPolicies[1]; got != "presentation-definition" {
		t.Fatalf("expected second vp policy presentation-definition, got %#v", got)
	}

	if len(capturedRequest.VCPolicies) != 2 {
		t.Fatalf("expected two vc policies, got %+v", capturedRequest.VCPolicies)
	}

	if got := capturedRequest.VCPolicies[0]; got != "signature" {
		t.Fatalf("expected first vc policy signature, got %#v", got)
	}

	webhookPolicy, ok := capturedRequest.VCPolicies[1].(map[string]interface{})
	if !ok {
		t.Fatalf("expected webhook policy object, got %#v", capturedRequest.VCPolicies[1])
	}

	if webhookPolicy["policy"] != "webhook" {
		t.Fatalf("expected webhook policy name, got %#v", webhookPolicy)
	}

	if webhookPolicy["args"] != "http://host.docker.internal:8787/api/verifier/policies/vc" {
		t.Fatalf("expected webhook policy url, got %#v", webhookPolicy)
	}

	if session.QRCodeURL != "openid4vp://authorize?response_type=vp_token&state=req-123&presentation_definition_uri=https%3A%2F%2Fverifier.example.com%2Fopenid4vc%2Fpd%2Freq-123" {
		t.Fatalf("expected legacy presentation url to round-trip, got %q", session.QRCodeURL)
	}

	if session.DeepLink != session.QRCodeURL {
		t.Fatalf("expected deep link to mirror presentation url, got qr=%q deep=%q", session.QRCodeURL, session.DeepLink)
	}

	if session.OfferURL != session.QRCodeURL {
		t.Fatalf("expected offer url to mirror presentation url, got offer=%q qr=%q", session.OfferURL, session.QRCodeURL)
	}

	if session.ProviderSessionID != "req-123" {
		t.Fatalf("expected provider session id req-123, got %q", session.ProviderSessionID)
	}
}

func TestProviderCreateSessionMapsRequestedCredentialTypes(t *testing.T) {
	var capturedRequest legacyCreateSessionRequest
	httpClient := roundTripClient(func(request *http.Request) (*http.Response, error) {
		if err := json.NewDecoder(request.Body).Decode(&capturedRequest); err != nil {
			t.Fatalf("expected valid walt request body, got error: %v", err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Content-Type": []string{"text/plain; charset=UTF-8"},
			},
			Body: io.NopCloser(strings.NewReader("openid4vp://authorize?response_type=vp_token&state=req-123")),
		}, nil
	})

	client := NewClient("http://host.docker.internal:7003", "legacy", "", "", "", "", "", "", httpClient)
	provider := New(client)

	_, err := provider.CreateSession(context.Background(), verification.ProviderSessionInput{
		RequestID:        "req-123",
		BusinessRef:      "job-123",
		CandidateRef:     "cand-456",
		Provider:         "walt",
		CredentialTypes:  []string{"DiplomaCredential", "TranscriptCredential"},
	})
	if err != nil {
		t.Fatalf("expected create session to succeed, got error: %v", err)
	}

	if len(capturedRequest.RequestCredentials) != 2 {
		t.Fatalf("expected two request credentials, got %+v", capturedRequest.RequestCredentials)
	}

	for index, expectedType := range []string{"DiplomaCredential", "TranscriptCredential"} {
		credential := capturedRequest.RequestCredentials[index]
		if credential.Type != expectedType {
			t.Fatalf("expected credential %d type %q, got %+v", index, expectedType, credential)
		}
		if credential.InputDescriptor.ID != expectedType {
			t.Fatalf("expected descriptor %d id %q, got %+v", index, expectedType, credential.InputDescriptor)
		}
		field := credential.InputDescriptor.Constraints.Fields[0]
		if field.Filter.Contains == nil || field.Filter.Contains.Const != expectedType {
			t.Fatalf("expected descriptor %d to filter type %q, got %+v", index, expectedType, field.Filter)
		}
	}
}

func TestProviderCreateSessionOmitsAuthorizationHeaderWhenBearerTokenMissing(t *testing.T) {
	httpClient := roundTripClient(func(request *http.Request) (*http.Response, error) {
		if got := request.Header.Get("Authorization"); got != "" {
			t.Fatalf("expected no authorization header, got %q", got)
		}

		if got := request.Header.Get("successRedirectUri"); got != "" {
			t.Fatalf("expected no successRedirectUri header, got %q", got)
		}

		if got := request.Header.Get("errorRedirectUri"); got != "" {
			t.Fatalf("expected no errorRedirectUri header, got %q", got)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Content-Type": []string{"text/plain; charset=UTF-8"},
			},
			Body: io.NopCloser(strings.NewReader("openid4vp://authorize?response_type=vp_token&state=req-123&presentation_definition_uri=https%3A%2F%2Fverifier.example.com%2Fopenid4vc%2Fpd%2Freq-123")),
		}, nil
	})

	client := NewClient("http://host.docker.internal:7003", "legacy", "", "", "", "", "", "", httpClient)
	provider := New(client)

	_, err := provider.CreateSession(context.Background(), verification.ProviderSessionInput{
		RequestID:    "req-123",
		BusinessRef:  "job-123",
		CandidateRef: "cand-456",
		Provider:     "walt",
	})
	if err != nil {
		t.Fatalf("expected create session to succeed without bearer token, got error: %v", err)
	}
}

func TestProviderCreateSessionVerifier2ModeStillSupported(t *testing.T) {
	var capturedRequest verifier2CreateSessionRequest
	httpClient := roundTripClient(func(request *http.Request) (*http.Response, error) {
		if request.URL.String() != "http://host.docker.internal:7004/verification-session/create" {
			t.Fatalf("expected verifier2 url, got %q", request.URL.String())
		}

		if err := json.NewDecoder(request.Body).Decode(&capturedRequest); err != nil {
			t.Fatalf("expected valid verifier2 body, got error: %v", err)
		}

		return &http.Response{
			StatusCode: http.StatusOK,
			Header: http.Header{
				"Content-Type": []string{"application/json"},
			},
			Body: io.NopCloser(strings.NewReader(`{
				"sessionId": "provider-session-456",
				"bootstrapAuthorizationRequestUrl": "openid4vp://authorize?client_id=verifier2&request_uri=https%3A%2F%2Fverifier.example.com%2Frequest",
				"fullAuthorizationRequestUrl": "openid4vp://authorize?response_type=vp_token&client_id=verifier2"
			}`)),
		}, nil
	})

	client := NewClient("http://host.docker.internal:7004", "verifier2", "", "", "", "", "", "", httpClient)
	provider := New(client)

	session, err := provider.CreateSession(context.Background(), verification.ProviderSessionInput{
		RequestID:    "req-456",
		BusinessRef:  "job-123",
		CandidateRef: "cand-456",
		Provider:     "walt",
	})
	if err != nil {
		t.Fatalf("expected verifier2 create session to succeed, got error: %v", err)
	}

	if capturedRequest.CoreFlow.DCQLQuery.Credentials[0].ID != "req-456" {
		t.Fatalf("expected dcql credential id req-456, got %+v", capturedRequest.CoreFlow.DCQLQuery.Credentials)
	}

	if capturedRequest.CoreFlow.DCQLQuery.Credentials[0].Meta.TypeValues[0][1] != "UniversityDegree" {
		t.Fatalf("expected verifier2 type UniversityDegree, got %+v", capturedRequest.CoreFlow.DCQLQuery.Credentials[0].Meta.TypeValues)
	}

	if session.ProviderSessionID != "provider-session-456" {
		t.Fatalf("expected verifier2 provider session id provider-session-456, got %q", session.ProviderSessionID)
	}
}

type roundTripClient func(*http.Request) (*http.Response, error)

func (client roundTripClient) Do(request *http.Request) (*http.Response, error) {
	return client(request)
}

func TestProviderParseCallbackNormalizesEvent(t *testing.T) {
	provider := New(NewClient("http://host.docker.internal:7003", "legacy", "bearer-token", "", "", "", "", "", http.DefaultClient))

	event, err := provider.ParseCallback(context.Background(), []byte(`{
		"id": "provider-session-123",
		"verificationResult": true,
		"policyResults": {
			"results": [
				{
					"policyResults": [
						{
							"policy": "webhook",
							"is_success": true,
							"result": {"decision":"accept"}
						}
					]
				}
			]
		}
	}`))
	if err != nil {
		t.Fatalf("expected callback parse to succeed, got error: %v", err)
	}

	if event.ProviderSessionID != "provider-session-123" {
		t.Fatalf("expected provider session id provider-session-123, got %q", event.ProviderSessionID)
	}

	if event.Status != verification.StatusVerified {
		t.Fatalf("expected normalized status %q, got %q", verification.StatusVerified, event.Status)
	}

	if !event.Verified {
		t.Fatal("expected normalized callback to be verified")
	}

	if event.EventType != "SESSION_VERIFIED" {
		t.Fatalf("expected SESSION_VERIFIED event type, got %q", event.EventType)
	}
}

func TestProviderParseCallbackMapsFailedPolicyResult(t *testing.T) {
	provider := New(NewClient("http://host.docker.internal:7003", "legacy", "bearer-token", "", "", "", "", "", http.DefaultClient))

	event, err := provider.ParseCallback(context.Background(), []byte(`{
		"id": "provider-session-456",
		"verificationResult": false,
		"policyResults": {
			"results": [
				{
					"policyResults": [
						{
							"policy": "webhook",
							"is_success": false,
							"error": "Issuer is not trusted by IU policy"
						}
					]
				}
			]
		}
	}`))
	if err != nil {
		t.Fatalf("expected callback parse to succeed, got error: %v", err)
	}

	if event.ProviderSessionID != "provider-session-456" {
		t.Fatalf("expected provider session id provider-session-456, got %q", event.ProviderSessionID)
	}

	if event.Status != verification.StatusFailed {
		t.Fatalf("expected normalized status %q, got %q", verification.StatusFailed, event.Status)
	}

	if event.Verified {
		t.Fatal("expected normalized callback to be unverified")
	}

	if event.ReasonCode != "POLICY_WEBHOOK_FAILED" {
		t.Fatalf("expected POLICY_WEBHOOK_FAILED reason code, got %q", event.ReasonCode)
	}

	if event.EventType != "SESSION_FAILED" {
		t.Fatalf("expected SESSION_FAILED event type, got %q", event.EventType)
	}
}
