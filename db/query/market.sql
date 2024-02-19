-- name: CreateTrade :exec
INSERT INTO trade_history (stock, quantity, buyer, seller) VALUES ($1, $2, $3, $4);

-- name: UpdateStockPrice :exec
UPDATE stocks SET price = $1, trend = $2, percentage_change = $3 WHERE id = $4;