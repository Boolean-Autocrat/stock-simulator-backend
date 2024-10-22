// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: token.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const createOrUpdateAccessToken = `-- name: CreateOrUpdateAccessToken :exec
INSERT INTO access_tokens ("user", token) VALUES ($1, $2) ON CONFLICT ("user") DO UPDATE SET token = $2
`

type CreateOrUpdateAccessTokenParams struct {
	User  uuid.UUID `json:"user"`
	Token string    `json:"token"`
}

func (q *Queries) CreateOrUpdateAccessToken(ctx context.Context, arg CreateOrUpdateAccessTokenParams) error {
	_, err := q.db.ExecContext(ctx, createOrUpdateAccessToken, arg.User, arg.Token)
	return err
}

const deleteAccessToken = `-- name: DeleteAccessToken :exec
DELETE FROM access_tokens WHERE "user" = $1
`

func (q *Queries) DeleteAccessToken(ctx context.Context, user uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteAccessToken, user)
	return err
}

const getTokenData = `-- name: GetTokenData :one
SELECT "user" FROM access_tokens WHERE token = $1
`

func (q *Queries) GetTokenData(ctx context.Context, token string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, getTokenData, token)
	var user uuid.UUID
	err := row.Scan(&user)
	return user, err
}
