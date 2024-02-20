-- name: CreateStock :one
INSERT INTO stocks (name, symbol, price, ipo_quantity, is_crypto, is_stock) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;

-- name: GetStock :one
SELECT * FROM stocks WHERE name = $1 AND is_crypto = $2 AND is_stock = $3 AND symbol = $4 AND price = $5;

-- name: GetStockById :one
SELECT * FROM stocks WHERE id = $1;

-- name: GetStocks :many
SELECT * FROM stocks;

-- name: SearchStocks :many
SELECT * FROM stocks WHERE LOWER(name) LIKE '%' || LOWER($1) || '%';

-- name: GetStockPriceHistory :many
SELECT price, price_at FROM price_history WHERE stock = $1 ORDER BY price_at DESC;

-- name: CreatePriceHistory :exec
INSERT INTO price_history (stock, price) VALUES ($1, $2);

-- name: GetTrendingStocks :many
SELECT stocks.id, stocks.name, stocks.symbol, stocks.price, stocks.ipo_quantity, stocks.is_crypto, stocks.is_stock, COUNT(price_history.id) AS price_history_count FROM stocks LEFT JOIN price_history ON stocks.id = price_history.stock GROUP BY stocks.id ORDER BY price_history_count DESC, stocks.name ASC LIMIT 10;

-- name: BuyStock :exec
UPDATE stocks SET in_circulation = in_circulation + $1 WHERE id = $2;

-- name: AddToIpoHistory :exec
INSERT INTO ipo_history ("user", stock, quantity, price) VALUES ($1, $2, $3, $4);

-- name: GetIpoHistory :many
SELECT s.name, i.quantity, i.price, i.created_at FROM ipo_history i JOIN stocks s ON s.id = i.stock WHERE "user" = $1;

-- name: GetWatchlist :many
SELECT stocks.name, stocks.symbol, stocks.price, stocks.is_crypto, stocks.is_stock, watchlist.added_at FROM stocks JOIN watchlist ON stocks.id = watchlist.stock WHERE watchlist."user" = $1;

-- name: AddToWatchlist :exec
INSERT INTO watchlist ("user", stock) VALUES ($1, $2);