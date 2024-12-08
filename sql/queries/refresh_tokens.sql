-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1,
    NOW(),
    NOW(),
    $2,
    $3
)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens
    WHERE token = $1;

-- name: SetRevokedAt :exec
UPDATE refresh_tokens
SET revoked_at = $1, 
    updated_at = $2
WHERE token = $3;
