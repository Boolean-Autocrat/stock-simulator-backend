-- name: CreateUser :one
INSERT INTO users
( full_name, email, picture) 
VALUES ( $1, $2, $3) RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;