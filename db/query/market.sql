-- name: CreateTrade :exec
INSERT INTO trade_history (stock, quantity, buyer, seller, price) VALUES ($1, $2, $3, $4, $5);

-- name: UpdateStockPrice :exec
UPDATE stocks SET price = $1, trend = $2, percentage_change = $3 WHERE id = $4;

-- name: CreatePendingOrder :one
INSERT INTO orders ("user", stock, quantity, price, is_buy) VALUES ($1, $2, $3, $4, $5) RETURNING "id";

-- name: UpdatePendingOrder :exec
UPDATE orders SET fulfilled_quantity = fulfilled_quantity + $1 WHERE id = $2;

-- name: GetPendingOrders :many
SELECT s.id AS stock_id, s.name AS stock, o.quantity, o.price, o.is_buy, o.fulfilled_quantity, o.created_at
FROM orders o
JOIN stocks s ON o.stock = s.id
WHERE o.user = $1;

