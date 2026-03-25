CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE verification_requests (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    business_ref TEXT NOT NULL,
    candidate_ref TEXT NOT NULL,
    provider TEXT NOT NULL,
    status TEXT NOT NULL CHECK (status IN (
        'CREATED',
        'SESSION_CREATED',
        'PENDING',
        'VERIFIED',
        'FAILED',
        'EXPIRED'
    )),
    verified BOOLEAN NOT NULL DEFAULT FALSE,
    reason_code TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE verification_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    verification_request_id UUID NOT NULL REFERENCES verification_requests(id) ON DELETE CASCADE,
    provider TEXT NOT NULL,
    provider_session_id TEXT NOT NULL UNIQUE,
    qr_code_url TEXT,
    deep_link TEXT,
    offer_url TEXT,
    expires_at TIMESTAMPTZ,
    raw_create_response JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE verification_events (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    verification_request_id UUID NOT NULL REFERENCES verification_requests(id) ON DELETE CASCADE,
    source TEXT NOT NULL,
    event_type TEXT NOT NULL,
    payload JSONB NOT NULL DEFAULT '{}'::JSONB,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX verification_requests_status_idx
    ON verification_requests (status);

CREATE INDEX verification_sessions_request_id_idx
    ON verification_sessions (verification_request_id);

CREATE INDEX verification_events_request_id_idx
    ON verification_events (verification_request_id, created_at);
