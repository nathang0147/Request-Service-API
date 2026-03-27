package config

import (
	"strings"
	"testing"
)

func TestLoadAppliesDefaults(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("PORT", "")
	t.Setenv("LOG_LEVEL", "")
	t.Setenv("DEFAULT_PROVIDER", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected config to load, got error: %v", err)
	}

	if cfg.Port != "8080" {
		t.Fatalf("expected default port 8080, got %q", cfg.Port)
	}

	if cfg.LogLevel != "info" {
		t.Fatalf("expected default log level info, got %q", cfg.LogLevel)
	}

	if cfg.DefaultProvider != "walt" {
		t.Fatalf("expected default provider walt, got %q", cfg.DefaultProvider)
	}

	if cfg.WaltVerifierMode != "legacy" {
		t.Fatalf("expected default verifier mode legacy, got %q", cfg.WaltVerifierMode)
	}
}

func TestLoadFailsWhenRequiredEnvMissing(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("DATABASE_URL", "")

	_, err := Load()
	if err == nil {
		t.Fatal("expected missing DATABASE_URL to fail")
	}

	if !strings.Contains(err.Error(), "DATABASE_URL") {
		t.Fatalf("expected error to mention DATABASE_URL, got %v", err)
	}
}

func TestLoadAllowsMissingOptionalWaltBearerToken(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("WALT_BEARER_TOKEN", "")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected config to load without bearer token, got error: %v", err)
	}

	if cfg.WaltBearerToken != "" {
		t.Fatalf("expected empty bearer token, got %q", cfg.WaltBearerToken)
	}
}

func TestLoadReadsOptionalWaltVCPolicyWebhookURL(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("WALT_VC_POLICY_WEBHOOK_URL", "http://host.docker.internal:8787/api/verifier/policies/vc")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected config to load with webhook url, got error: %v", err)
	}

	if cfg.WaltVCPolicyWebhookURL != "http://host.docker.internal:8787/api/verifier/policies/vc" {
		t.Fatalf("expected webhook url to load, got %q", cfg.WaltVCPolicyWebhookURL)
	}
}

func TestLoadReadsOptionalPublicRedirectURLTemplate(t *testing.T) {
	setRequiredEnv(t)
	t.Setenv("PUBLIC_REDIRECT_URL_TEMPLATE", "http://localhost:7102/success/$id")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("expected config to load with redirect template, got error: %v", err)
	}

	if cfg.PublicRedirectTemplate != "http://localhost:7102/success/$id" {
		t.Fatalf("expected redirect template to load, got %q", cfg.PublicRedirectTemplate)
	}
}

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/request_service")
	t.Setenv("WALT_VERIFIER_BASE_URL", "http://host.docker.internal:7003")
	t.Setenv("WALT_BEARER_TOKEN", "test-bearer-token")
	t.Setenv("CALLBACK_BASE_URL", "https://request-service.example.com")
	t.Setenv("CALLBACK_AUTH_SECRET", "test-callback-secret")
}
