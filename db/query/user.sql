-- name: CreateUser :one
INSERT INTO users
  (full_name, email, picture) 
VALUES 
  ($1, $2, $3) 
ON CONFLICT (email) DO UPDATE
  SET full_name = excluded.full_name, picture = excluded.picture
RETURNING *;

-- name: GetUser :one
SELECT full_name, email, picture, balance FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: GetLeaderboard :many
SELECT id, full_name, picture, balance FROM users ORDER BY balance DESC, full_name ASC LIMIT 10;

-- name: GetUserPosition :one
SELECT id, full_name, picture, balance, position FROM (
  SELECT id, full_name, picture, balance, row_number() OVER (ORDER BY balance DESC) AS position FROM users
) AS users_with_position WHERE id = $1;

-- name: GetUserBalance :one
SELECT balance FROM users WHERE id = $1;

-- name: UpdateBalance :exec
UPDATE users SET balance = balance + $1 WHERE id = $2;

-- name: GetDevelopers :many
SELECT * FROM developers ORDER BY id;