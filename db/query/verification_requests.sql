-- name: CreateVerificationRequest :one
INSERT INTO verification_requests (
    business_ref,
    candidate_ref,
    provider,
    status,
    verified,
    reason_code
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetVerificationRequestByID :one
SELECT *
FROM verification_requests
WHERE id = $1
LIMIT 1;

-- name: UpdateVerificationRequestStatus :one
UPDATE verification_requests
SET
    status = $2,
    verified = $3,
    reason_code = $4,
    updated_at = NOW()
WHERE id = $1
RETURNING *;
