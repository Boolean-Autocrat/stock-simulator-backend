-- name: GetTokenData :one
SELECT user_id, expires_at FROM access_tokens WHERE token = $1;

-- name: CreateAccessToken :one
INSERT INTO access_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3)
ON CONFLICT (user_id) DO UPDATE
SET token = excluded.token, expires_at = excluded.expires_at
RETURNING *;

-- name: UpdateAccessToken :one
UPDATE access_tokens SET expires_at = $1 AND token = $2 WHERE user_id = $3 RETURNING *;

-- name: DeleteAccessToken :exec
DELETE FROM access_tokens WHERE user_id = $1;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE user_id = $1;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token, expires_at)
VALUES ($1, $2, $3)
ON CONFLICT (user_id) DO UPDATE
SET token = excluded.token, expires_at = excluded.expires_at
RETURNING *;

-- name: UpdateRefreshToken :one
UPDATE refresh_tokens SET expires_at = $1 AND token = $2 WHERE user_id = $3 RETURNING *;

-- name: DeleteRefreshToken :exec
DELETE FROM refresh_tokens WHERE user_id = $1;