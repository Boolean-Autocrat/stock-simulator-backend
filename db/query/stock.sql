-- name: CreateStock :one
INSERT INTO stocks (name, symbol, price, quantity, is_crypto, is_stock) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetStock :one
SELECT * FROM stocks WHERE name = $1 AND is_crypto = $2 AND is_stock = $3 AND symbol = $4 AND price = $5;

-- name: GetStocks :many
SELECT * FROM stocks;

-- name: SearchStocks :many
SELECT * FROM stocks WHERE LOWER(name) LIKE '%' || LOWER($1) || '%';

-- name: GetStockPriceHistory :many
SELECT * FROM price_history WHERE stock_id = $1;

-- name: GetStockPriceHistoryByDate :many
SELECT * FROM price_history WHERE stock_id = $1 AND price_at >= $2 AND price_at <= $3;

-- name: CreatePriceHistory :one
INSERT INTO price_history (stock_id, price) VALUES ($1, $2) RETURNING *;