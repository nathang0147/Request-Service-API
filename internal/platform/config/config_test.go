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

func setRequiredEnv(t *testing.T) {
	t.Helper()
	t.Setenv("DATABASE_URL", "postgres://user:pass@localhost:5432/request_service")
	t.Setenv("WALT_BASE_URL", "https://walt.example.com")
	t.Setenv("WALT_API_KEY", "test-api-key")
	t.Setenv("CALLBACK_BASE_URL", "https://request-service.example.com")
	t.Setenv("CALLBACK_AUTH_SECRET", "test-callback-secret")
}
