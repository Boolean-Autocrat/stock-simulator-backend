// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users
( full_name, email, picture) 
VALUES ( $1, $2, $3) RETURNING id, full_name, email, picture, balance
`

type CreateUserParams struct {
	FullName sql.NullString `json:"fullName"`
	Email    sql.NullString `json:"email"`
	Picture  sql.NullString `json:"picture"`
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser, arg.FullName, arg.Email, arg.Picture)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Picture,
		&i.Balance,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, full_name, email, picture, balance FROM users WHERE id = $1
`

func (q *Queries) GetUser(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Picture,
		&i.Balance,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, full_name, email, picture, balance FROM users WHERE email = $1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.FullName,
		&i.Email,
		&i.Picture,
		&i.Balance,
	)
	return i, err
}
