-- name: AddStockToPortfolio :one
INSERT INTO portfolio (user_id, stock_id, purchase_price) VALUES ($1, $2, $3) RETURNING *;

-- name: GetPortfolio :many
SELECT * FROM portfolio WHERE user_id = $1;

-- name: GetPortfolioStocks :many
SELECT * FROM portfolio WHERE user_id = $1 AND stock_id = $2;

-- name: GetPortfolioStock :one
SELECT * FROM portfolio WHERE user_id = $1 AND stock_id = $2;

-- name: RemoveStockFromPortfolio :exec
DELETE FROM portfolio WHERE user_id = $1 AND stock_id = $2;