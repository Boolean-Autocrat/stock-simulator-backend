-- name: GetTokenData :one
SELECT "user" FROM access_tokens WHERE token = $1;

-- name: CreateOrUpdateAccessToken :exec
INSERT INTO access_tokens ("user", token) VALUES ($1, $2) ON CONFLICT ("user") DO UPDATE SET token = $2;

-- name: DeleteAccessToken :exec
DELETE FROM access_tokens WHERE "user" = $1;