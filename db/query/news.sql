-- name: AddArticle :one
INSERT INTO news (title, description, photo) VALUES ($1, $2, $3) RETURNING *;

-- name: GetArticle :one
SELECT * FROM news WHERE id = $1;

-- name: GetArticles :many
SELECT * FROM news;

-- name: GetArticleSentiment :many
SELECT * FROM news_sentiment WHERE article_id = $1;

-- name: AddArticleSentiment :one
INSERT INTO news_sentiment (article_id, user_id, "like", "dislike") VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateArticleSentiment :one
UPDATE news_sentiment SET "like" = $1,"dislike" = $2 WHERE article_id = $3 AND user_id = $4 RETURNING *;

-- name: GetArticleSentimentByUser :one
SELECT * FROM news_sentiment WHERE article_id = $1 AND user_id = $2;