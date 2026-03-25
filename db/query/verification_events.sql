-- name: CreateVerificationEvent :one
INSERT INTO verification_events (
    verification_request_id,
    source,
    event_type,
    payload
) VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: ListVerificationEventsByRequestID :many
SELECT *
FROM verification_events
WHERE verification_request_id = $1
ORDER BY created_at ASC;
