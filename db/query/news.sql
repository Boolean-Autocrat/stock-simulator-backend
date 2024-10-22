-- name: AddArticle :one
INSERT INTO news (title, author, content, tag) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: UpdateArticle :exec
UPDATE news SET title = $1, author = $2, content = $3, tag = $4 WHERE id = $5 RETURNING *;

-- name: GetArticle :one
SELECT * FROM news WHERE id = $1;

-- name: GetUserSentiment :one
SELECT * FROM news_sentiment WHERE article = $1 AND "user" = $2;

-- name: GetArticleSentiment :one
SELECT COUNT(CASE WHEN "like" THEN 1 END) AS like_count, COUNT(CASE WHEN "dislike" THEN 1 END) AS dislike_count FROM news_sentiment WHERE article = $1;

-- name: GetArticles :many
SELECT * FROM news;

-- name: AddArticleSentiment :exec
INSERT INTO news_sentiment(article, "user", "like", "dislike")
VALUES ($1, $2, $3, $4)
ON CONFLICT (article, "user") DO UPDATE
SET "like" = excluded.like, "dislike" = excluded.dislike;

-- name: DeleteArticle :exec
DELETE FROM news WHERE id = $1;