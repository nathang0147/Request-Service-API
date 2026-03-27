package walt

import "encoding/json"

type verifierMode string

const (
	legacyVerifierMode    verifierMode = "legacy"
	verifier2VerifierMode verifierMode = "verifier2"
)

type createSessionRequest struct {
	Legacy    *legacyCreateSessionRequest
	Verifier2 *verifier2CreateSessionRequest
}

type legacyCreateSessionRequest struct {
	StateID            string                      `json:"-"`
	VPPolicies         []any                       `json:"vp_policies,omitempty"`
	VCPolicies         []any                       `json:"vc_policies,omitempty"`
	RequestCredentials []legacyRequestedCredential `json:"request_credentials"`
}

type legacyRequestedCredential struct {
	Format          string                `json:"format"`
	Type            string                `json:"type"`
	InputDescriptor legacyInputDescriptor `json:"input_descriptor"`
}

type legacyInputDescriptor struct {
	ID          string                   `json:"id"`
	Format      map[string]legacyAlgSpec `json:"format"`
	Constraints legacyConstraints        `json:"constraints"`
}

type legacyAlgSpec struct {
	Alg []string `json:"alg,omitempty"`
}

type legacyConstraints struct {
	Fields          []legacyField `json:"fields"`
	LimitDisclosure string        `json:"limit_disclosure,omitempty"`
}

type legacyField struct {
	Path   []string          `json:"path"`
	Filter legacyFieldFilter `json:"filter"`
}

type legacyFieldFilter struct {
	Type     string            `json:"type"`
	Contains *legacyFieldConst `json:"contains,omitempty"`
	Pattern  string            `json:"pattern,omitempty"`
}

type legacyFieldConst struct {
	Const string `json:"const"`
}

type verifier2CreateSessionRequest struct {
	FlowType string                     `json:"flow_type"`
	CoreFlow verifier2CreateSessionCore `json:"core_flow"`
}

type verifier2CreateSessionCore struct {
	DCQLQuery verifier2CreateSessionDCQLQuery `json:"dcql_query"`
}

type verifier2CreateSessionDCQLQuery struct {
	Credentials []verifier2CreateSessionCredential `json:"credentials"`
}

type verifier2CreateSessionCredential struct {
	ID     string                               `json:"id"`
	Format string                               `json:"format"`
	Meta   verifier2CreateSessionCredentialMeta `json:"meta"`
	Claims []verifier2CreateSessionClaimQuery   `json:"claims"`
}

type verifier2CreateSessionCredentialMeta struct {
	TypeValues [][]string `json:"type_values"`
}

type verifier2CreateSessionClaimQuery struct {
	Path []string `json:"path"`
}

type createSessionResponse struct {
	LegacyURL string
	Verifier2 *verifier2CreateSessionResponse
}

type verifier2CreateSessionResponse struct {
	SessionID                        string `json:"sessionId"`
	BootstrapAuthorizationRequestURL string `json:"bootstrapAuthorizationRequestUrl"`
	FullAuthorizationRequestURL      string `json:"fullAuthorizationRequestUrl"`
	CreationTarget                   string `json:"creationTarget"`
}

type callbackPayload struct {
	SessionID  string          `json:"sessionId"`
	Status     string          `json:"status"`
	Verified   bool            `json:"verified"`
	ReasonCode string          `json:"reasonCode"`
	EventType  string          `json:"eventType"`
	Payload    json.RawMessage `json:"payload"`
}

type statusCallbackPayload struct {
	ID                 string                       `json:"id"`
	VerificationResult *bool                        `json:"verificationResult"`
	PolicyResults      *statusCallbackPolicyResults `json:"policyResults"`
}

type statusCallbackPolicyResults struct {
	Results []statusCallbackPolicyEntry `json:"results"`
}

type statusCallbackPolicyEntry struct {
	PolicyResults []statusCallbackPolicyResult `json:"policyResults"`
}

type statusCallbackPolicyResult struct {
	Policy    string          `json:"policy"`
	IsSuccess bool            `json:"is_success"`
	Result    json.RawMessage `json:"result"`
	Error     json.RawMessage `json:"error"`
}
