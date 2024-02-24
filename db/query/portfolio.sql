-- name: AddOrUpdateStockToPortfolio :exec
INSERT INTO portfolio ("user", "stock", quantity) VALUES ($1, $2, $3) ON CONFLICT ("user", "stock") DO UPDATE SET quantity = portfolio.quantity + $3;

-- name: GetPortfolio :many
SELECT p."stock", s.name, s.symbol, s.price, s.is_crypto, s.is_stock, s.trend, s.percentage_change, p.quantity 
FROM portfolio p
JOIN stocks s ON p."stock" = s.id
WHERE p."user" = $1;

-- name: GetStocksAndQuantity :many
SELECT SUM(quantity) AS quantity, "stock" FROM portfolio WHERE "user" = $1 GROUP BY "stock";

-- name: GetStockWithQuantity :one
SELECT SUM(quantity) AS quantity, "stock" FROM portfolio WHERE "user" = $1 AND "stock" = $2 GROUP BY "stock";