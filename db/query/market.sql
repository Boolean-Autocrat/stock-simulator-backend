-- name: UpdateBalance :exec
UPDATE users SET balance = balance + $1 WHERE id = $2;

-- name: AddSellOrder :exec
INSERT INTO sell_orders ("id", "user", "stock", "price", "quantity", "fulfilled") VALUES ($1, $2, $3, $4, $5, $6);

-- name: UpdateSellOrder :exec
UPDATE sell_orders SET "fulfilled" = "fulfilled" + $1 WHERE "id" = $2;

-- name: UpdateBuyOrder :exec
UPDATE buy_orders SET "fulfilled" = "fulfilled" + $1 WHERE "id" = $2;

-- name: AddBuyOrder :exec
INSERT INTO buy_orders ("id", "user", "stock", "price", "quantity", "fulfilled") VALUES ($1, $2, $3, $4, $5, $6);

-- name: AddTrade :exec
INSERT INTO trades ("id", "buy_order", "sell_order", "price", "quantity") VALUES ($1, $2, $3, $4, $5);

-- name: GetAllSellOrdersByUser :many
SELECT *, fulfilled < quantity AS is_complete FROM sell_orders WHERE "user" = $1;

-- name: GetAllBuyOrdersByUser :many
SELECT *, fulfilled < quantity AS is_complete FROM buy_orders WHERE "user" = $1;

-- name: ListBuyOrders :many
SELECT * FROM buy_orders WHERE "stock" = $1 AND "fulfilled" < "quantity" ORDER BY "price" DESC;

-- name: ListSellOrders :many
SELECT * FROM sell_orders WHERE "stock" = $1 AND "fulfilled" < "quantity" ORDER BY "price" ASC;