-- name: CreateVerificationSession :one
INSERT INTO verification_sessions (
    verification_request_id,
    provider,
    provider_session_id,
    qr_code_url,
    deep_link,
    offer_url,
    expires_at,
    raw_create_response
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
RETURNING *;

-- name: GetVerificationSessionByProviderSessionID :one
SELECT *
FROM verification_sessions
WHERE provider_session_id = $1
LIMIT 1;

-- name: GetLatestVerificationSessionByRequestID :one
SELECT *
FROM verification_sessions
WHERE verification_request_id = $1
ORDER BY created_at DESC
LIMIT 1;
