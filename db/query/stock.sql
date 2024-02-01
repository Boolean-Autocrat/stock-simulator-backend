-- name: CreateStock :one
INSERT INTO stocks (name, symbol, price, quantity, is_crypto, is_stock) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetStock :one
SELECT * FROM stocks WHERE name = $1 AND is_crypto = $2 AND is_stock = $3 AND symbol = $4 AND price = $5;

-- name: GetStockById :one
SELECT * FROM stocks WHERE id = $1;

-- name: GetStocks :many
SELECT * FROM stocks;

-- name: SearchStocks :many
SELECT * FROM stocks WHERE LOWER(name) LIKE '%' || LOWER($1) || '%';

-- name: GetStockPriceHistory :many
SELECT price, price_at FROM price_history WHERE stock_id = $1 ORDER BY price_at DESC;

-- name: CreatePriceHistory :exec
INSERT INTO price_history (stock_id, price) VALUES ($1, $2);

-- name: GetTrendingStocks :many
SELECT stocks.id, stocks.name, stocks.symbol, stocks.price, stocks.quantity, stocks.is_crypto, stocks.is_stock, COUNT(price_history.id) AS price_history_count FROM stocks LEFT JOIN price_history ON stocks.id = price_history.stock_id GROUP BY stocks.id ORDER BY price_history_count DESC, stocks.name ASC LIMIT 10;