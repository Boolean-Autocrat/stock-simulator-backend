// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: news.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const addArticle = `-- name: AddArticle :one
INSERT INTO news (title, author, content, tag) VALUES ($1, $2, $3, $4) RETURNING id, title, author, content, tag, image, created_at
`

type AddArticleParams struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Content string `json:"content"`
	Tag     string `json:"tag"`
}

func (q *Queries) AddArticle(ctx context.Context, arg AddArticleParams) (News, error) {
	row := q.db.QueryRowContext(ctx, addArticle,
		arg.Title,
		arg.Author,
		arg.Content,
		arg.Tag,
	)
	var i News
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Author,
		&i.Content,
		&i.Tag,
		&i.Image,
		&i.CreatedAt,
	)
	return i, err
}

const addArticleSentiment = `-- name: AddArticleSentiment :exec
INSERT INTO news_sentiment(article, "user", "like", "dislike")
VALUES ($1, $2, $3, $4)
ON CONFLICT (article, "user") DO UPDATE
SET "like" = excluded.like, "dislike" = excluded.dislike
`

type AddArticleSentimentParams struct {
	Article uuid.UUID `json:"article"`
	User    uuid.UUID `json:"user"`
	Like    bool      `json:"like"`
	Dislike bool      `json:"dislike"`
}

func (q *Queries) AddArticleSentiment(ctx context.Context, arg AddArticleSentimentParams) error {
	_, err := q.db.ExecContext(ctx, addArticleSentiment,
		arg.Article,
		arg.User,
		arg.Like,
		arg.Dislike,
	)
	return err
}

const deleteArticle = `-- name: DeleteArticle :exec
DELETE FROM news WHERE id = $1
`

func (q *Queries) DeleteArticle(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteArticle, id)
	return err
}

const getArticle = `-- name: GetArticle :one
SELECT id, title, author, content, tag, image, created_at FROM news WHERE id = $1
`

func (q *Queries) GetArticle(ctx context.Context, id uuid.UUID) (News, error) {
	row := q.db.QueryRowContext(ctx, getArticle, id)
	var i News
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Author,
		&i.Content,
		&i.Tag,
		&i.Image,
		&i.CreatedAt,
	)
	return i, err
}

const getArticleSentiment = `-- name: GetArticleSentiment :one
SELECT COUNT(CASE WHEN "like" THEN 1 END) AS like_count, COUNT(CASE WHEN "dislike" THEN 1 END) AS dislike_count FROM news_sentiment WHERE article = $1
`

type GetArticleSentimentRow struct {
	LikeCount    int64 `json:"likeCount"`
	DislikeCount int64 `json:"dislikeCount"`
}

func (q *Queries) GetArticleSentiment(ctx context.Context, article uuid.UUID) (GetArticleSentimentRow, error) {
	row := q.db.QueryRowContext(ctx, getArticleSentiment, article)
	var i GetArticleSentimentRow
	err := row.Scan(&i.LikeCount, &i.DislikeCount)
	return i, err
}

const getArticles = `-- name: GetArticles :many
SELECT id, title, author, content, tag, image, created_at FROM news
`

func (q *Queries) GetArticles(ctx context.Context) ([]News, error) {
	rows, err := q.db.QueryContext(ctx, getArticles)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []News
	for rows.Next() {
		var i News
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Author,
			&i.Content,
			&i.Tag,
			&i.Image,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUserSentiment = `-- name: GetUserSentiment :one
SELECT id, article, "user", "like", dislike FROM news_sentiment WHERE article = $1 AND "user" = $2
`

type GetUserSentimentParams struct {
	Article uuid.UUID `json:"article"`
	User    uuid.UUID `json:"user"`
}

func (q *Queries) GetUserSentiment(ctx context.Context, arg GetUserSentimentParams) (NewsSentiment, error) {
	row := q.db.QueryRowContext(ctx, getUserSentiment, arg.Article, arg.User)
	var i NewsSentiment
	err := row.Scan(
		&i.ID,
		&i.Article,
		&i.User,
		&i.Like,
		&i.Dislike,
	)
	return i, err
}

const updateArticle = `-- name: UpdateArticle :exec
UPDATE news SET title = $1, author = $2, content = $3, tag = $4 WHERE id = $5 RETURNING id, title, author, content, tag, image, created_at
`

type UpdateArticleParams struct {
	Title   string    `json:"title"`
	Author  string    `json:"author"`
	Content string    `json:"content"`
	Tag     string    `json:"tag"`
	ID      uuid.UUID `json:"id"`
}

func (q *Queries) UpdateArticle(ctx context.Context, arg UpdateArticleParams) error {
	_, err := q.db.ExecContext(ctx, updateArticle,
		arg.Title,
		arg.Author,
		arg.Content,
		arg.Tag,
		arg.ID,
	)
	return err
}
