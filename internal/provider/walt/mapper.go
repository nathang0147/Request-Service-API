package walt

import (
	"encoding/json"
	"strings"

	"github.com/nathang0147/Request-Service-API/internal/provider"
	"github.com/nathang0147/Request-Service-API/internal/verification"
)

func toCreateSessionRequest(input verification.ProviderSessionInput, vcPolicyWebhookURL string) createSessionRequest {
	const credentialType = "UniversityDegree"

	vcPolicies := []any{"signature"}
	if vcPolicyWebhookURL != "" {
		vcPolicies = append(vcPolicies, map[string]string{
			"policy": "webhook",
			"args":   vcPolicyWebhookURL,
		})
	}

	return createSessionRequest{
		Legacy: &legacyCreateSessionRequest{
			StateID: input.RequestID,
			VPPolicies: []any{
				"signature",
				"presentation-definition",
			},
			VCPolicies: vcPolicies,
			RequestCredentials: []legacyRequestedCredential{
				{
					Format: "jwt_vc_json",
					Type:   credentialType,
					InputDescriptor: legacyInputDescriptor{
						ID: credentialType,
						Format: map[string]legacyAlgSpec{
							"jwt_vc_json": {Alg: []string{"EdDSA"}},
						},
						Constraints: legacyConstraints{
							Fields: []legacyField{
								{
									Path: []string{"$.vc.type"},
									Filter: legacyFieldFilter{
										Type: "array",
										Contains: &legacyFieldConst{
											Const: credentialType,
										},
									},
								},
							},
							LimitDisclosure: "required",
						},
					},
				},
			},
		},
		Verifier2: &verifier2CreateSessionRequest{
			FlowType: "cross_device",
			CoreFlow: verifier2CreateSessionCore{
				DCQLQuery: verifier2CreateSessionDCQLQuery{
					Credentials: []verifier2CreateSessionCredential{
						{
							ID:     input.RequestID,
							Format: "jwt_vc_json",
							Meta: verifier2CreateSessionCredentialMeta{
								TypeValues: [][]string{
									{"VerifiableCredential", credentialType},
								},
							},
							Claims: []verifier2CreateSessionClaimQuery{
								{Path: []string{"name"}},
							},
						},
					},
				},
			},
		},
	}
}

func toProviderSession(response createSessionResponse, raw []byte) verification.ProviderSession {
	if response.Verifier2 == nil {
		providerSessionID := extractLegacyStateFromURL(response.LegacyURL)
		if providerSessionID == "" {
			providerSessionID = response.LegacyURL
		}

		rawCreateResponse, err := json.Marshal(map[string]string{
			"presentationUrl": response.LegacyURL,
		})
		if err != nil {
			rawCreateResponse = raw
		}

		return verification.ProviderSession{
			ProviderSessionID: providerSessionID,
			QRCodeURL:         response.LegacyURL,
			DeepLink:          response.LegacyURL,
			OfferURL:          response.LegacyURL,
			RawCreateResponse: rawCreateResponse,
		}
	}

	return verification.ProviderSession{
		ProviderSessionID: response.Verifier2.SessionID,
		QRCodeURL:         response.Verifier2.BootstrapAuthorizationRequestURL,
		DeepLink:          response.Verifier2.BootstrapAuthorizationRequestURL,
		OfferURL:          response.Verifier2.FullAuthorizationRequestURL,
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

func toStatusCallbackEvent(payload statusCallbackPayload, raw []byte) provider.CallbackEvent {
	event := provider.CallbackEvent{
		ProviderSessionID: payload.ID,
		Status:            verification.StatusPending,
		Verified:          false,
		ReasonCode:        "",
		EventType:         "SESSION_PENDING",
		Payload:           raw,
	}

	if payload.VerificationResult == nil {
		return event
	}

	if *payload.VerificationResult {
		event.Status = verification.StatusVerified
		event.Verified = true
		event.EventType = "SESSION_VERIFIED"
		return event
	}

	event.Status = verification.StatusFailed
	event.EventType = "SESSION_FAILED"
	event.ReasonCode = deriveFailureReasonCode(payload.PolicyResults)

	if strings.Contains(event.ReasonCode, "EXPIRED") {
		event.Status = verification.StatusExpired
		event.EventType = "SESSION_EXPIRED"
	}

	return event
}

func deriveFailureReasonCode(policyResults *statusCallbackPolicyResults) string {
	if policyResults == nil {
		return "VERIFICATION_FAILED"
	}

	for _, entry := range policyResults.Results {
		for _, result := range entry.PolicyResults {
			if result.IsSuccess {
				continue
			}

			switch strings.ToLower(strings.TrimSpace(result.Policy)) {
			case "expired":
				return "CREDENTIAL_EXPIRED"
			case "":
				return "VERIFICATION_FAILED"
			default:
				slug := strings.ToUpper(strings.ReplaceAll(strings.TrimSpace(result.Policy), "-", "_"))
				return "POLICY_" + slug + "_FAILED"
			}
		}
	}

	return "VERIFICATION_FAILED"
}
