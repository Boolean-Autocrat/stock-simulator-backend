-- name: CreateUser :one
INSERT INTO users
( full_name, email, picture) 
VALUES ( $1, $2, $3) RETURNING *;