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