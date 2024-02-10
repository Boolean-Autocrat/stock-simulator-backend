-- name: AddOrUpdateStockToPortfolio :one
INSERT INTO portfolio (user_id, stock_id, quantity) VALUES ($1, $2, $3) ON CONFLICT (stock_id) DO UPDATE SET quantity = portfolio.quantity + $3 RETURNING *;

-- name: UpdateStockQuantityPortfolio :exec
UPDATE portfolio SET quantity = $3 WHERE user_id = $1 AND stock_id = $2;

-- name: GetPortfolio :many
SELECT p.stock_id, s.name, s.symbol, s.price, s.is_crypto, s.is_stock
FROM portfolio p
JOIN stocks s ON p.stock_id = s.id
WHERE p.user_id = $1
LIMIT 10
OFFSET $2;

-- name: GetStocksAndQuantity :many
SELECT SUM(quantity) AS quantity, stock_id FROM portfolio WHERE user_id = $1 GROUP BY stock_id;

-- name: GetStockWithQuantity :one
SELECT SUM(quantity) AS quantity, stock_id FROM portfolio WHERE user_id = $1 AND stock_id = $2 GROUP BY stock_id;

-- name: RemoveStockFromPortfolio :exec
DELETE FROM portfolio WHERE user_id = $1 AND stock_id = $2;