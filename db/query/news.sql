-- name: AddArticle :one
INSERT INTO news (title, author, content, tag) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetArticle :one
SELECT * FROM news WHERE id = $1;

-- name: GetUserSentiment :one
SELECT * FROM news_sentiment WHERE article_id = $1 AND user_id = $2;

-- name: GetArticleSentiment :one
SELECT COUNT(CASE WHEN "like" THEN 1 END) AS like_count, COUNT(CASE WHEN "dislike" THEN 1 END) AS dislike_count FROM news_sentiment WHERE article_id = $1;

-- name: GetArticles :many
SELECT * FROM news;

-- name: AddArticleSentiment :exec
INSERT INTO news_sentiment(article_id, user_id, "like", "dislike")
VALUES ($1, $2, $3, $4)
ON CONFLICT (user_id) DO UPDATE
SET "like" = excluded.like, "dislike" = excluded.dislike;