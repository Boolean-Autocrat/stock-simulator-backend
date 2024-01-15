-- name: AddStockToPortfolio :one
INSERT INTO portfolio (user_id, stock_id, purchase_price) VALUES ($1, $2, $3) RETURNING *;

-- name: GetPortfolio :many
SELECT p.stock_id, p.purchase_price, p.purchased_at, s.name, s.symbol, s.price, s.is_crypto, s.is_stock
FROM portfolio p
JOIN stocks s ON p.stock_id = s.id
WHERE p.user_id = $1
ORDER BY p.purchased_at
LIMIT 10
OFFSET $2;

-- name: RemoveStockFromPortfolio :exec
DELETE FROM portfolio WHERE user_id = $1 AND stock_id = $2;