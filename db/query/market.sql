-- name: CreateTrade :exec
INSERT INTO trade_history (stock, quantity, buyer, seller, price) VALUES ($1, $2, $3, $4, $5);

-- name: UpdateStockPrice :exec
UPDATE stocks SET price = $1, trend = $2, percentage_change = $3 WHERE id = $4;

-- name: CreatePendingOrder :one
INSERT INTO orders ("user", stock, quantity, price, is_buy) VALUES ($1, $2, $3, $4, $5) RETURNING "id";

-- name: UpdatePendingOrder :exec
UPDATE orders SET fulfilled_quantity = fulfilled_quantity + $1 WHERE id = $2;

-- name: GetPendingOrders :many
SELECT stock, quantity, price, is_buy, fulfilled_quantity, "created_at" FROM orders WHERE "user" = $1;