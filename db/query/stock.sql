-- name: CreateStock :one
INSERT INTO stocks (symbol, price, is_crypto, is_stock) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: GetStock :one
SELECT * FROM stocks WHERE symbol = $1;

-- name: GetStocks :many
SELECT * FROM stocks WHERE is_crypto = $1 AND is_stock = $2;

-- name: GetStocksBySymbol :many
SELECT * FROM stocks WHERE symbol LIKE $1;

-- name: GetStocksByName :many
SELECT * FROM stocks WHERE symbol LIKE $1;