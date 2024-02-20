// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: market.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPendingOrder = `-- name: CreatePendingOrder :one
INSERT INTO orders ("user", stock, quantity, price, is_buy) VALUES ($1, $2, $3, $4, $5) RETURNING "id"
`

type CreatePendingOrderParams struct {
	User     uuid.UUID `json:"user"`
	Stock    uuid.UUID `json:"stock"`
	Quantity int32     `json:"quantity"`
	Price    float32   `json:"price"`
	IsBuy    bool      `json:"isBuy"`
}

func (q *Queries) CreatePendingOrder(ctx context.Context, arg CreatePendingOrderParams) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, createPendingOrder,
		arg.User,
		arg.Stock,
		arg.Quantity,
		arg.Price,
		arg.IsBuy,
	)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const createTrade = `-- name: CreateTrade :exec
INSERT INTO trade_history (stock, quantity, buyer, seller, price) VALUES ($1, $2, $3, $4, $5)
`

type CreateTradeParams struct {
	Stock    uuid.UUID `json:"stock"`
	Quantity int32     `json:"quantity"`
	Buyer    uuid.UUID `json:"buyer"`
	Seller   uuid.UUID `json:"seller"`
	Price    float32   `json:"price"`
}

func (q *Queries) CreateTrade(ctx context.Context, arg CreateTradeParams) error {
	_, err := q.db.ExecContext(ctx, createTrade,
		arg.Stock,
		arg.Quantity,
		arg.Buyer,
		arg.Seller,
		arg.Price,
	)
	return err
}

const getPendingOrders = `-- name: GetPendingOrders :many
SELECT stock, quantity, price, is_buy, fulfilled_quantity, "created_at" FROM orders WHERE "user" = $1
`

type GetPendingOrdersRow struct {
	Stock             uuid.UUID `json:"stock"`
	Quantity          int32     `json:"quantity"`
	Price             float32   `json:"price"`
	IsBuy             bool      `json:"isBuy"`
	FulfilledQuantity int32     `json:"fulfilledQuantity"`
	CreatedAt         time.Time `json:"createdAt"`
}

func (q *Queries) GetPendingOrders(ctx context.Context, user uuid.UUID) ([]GetPendingOrdersRow, error) {
	rows, err := q.db.QueryContext(ctx, getPendingOrders, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPendingOrdersRow
	for rows.Next() {
		var i GetPendingOrdersRow
		if err := rows.Scan(
			&i.Stock,
			&i.Quantity,
			&i.Price,
			&i.IsBuy,
			&i.FulfilledQuantity,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updatePendingOrder = `-- name: UpdatePendingOrder :exec
UPDATE orders SET fulfilled_quantity = fulfilled_quantity + $1 WHERE id = $2
`

type UpdatePendingOrderParams struct {
	FulfilledQuantity int32     `json:"fulfilledQuantity"`
	ID                uuid.UUID `json:"id"`
}

func (q *Queries) UpdatePendingOrder(ctx context.Context, arg UpdatePendingOrderParams) error {
	_, err := q.db.ExecContext(ctx, updatePendingOrder, arg.FulfilledQuantity, arg.ID)
	return err
}

const updateStockPrice = `-- name: UpdateStockPrice :exec
UPDATE stocks SET price = $1, trend = $2, percentage_change = $3 WHERE id = $4
`

type UpdateStockPriceParams struct {
	Price            float32   `json:"price"`
	Trend            string    `json:"trend"`
	PercentageChange float32   `json:"percentageChange"`
	ID               uuid.UUID `json:"id"`
}

func (q *Queries) UpdateStockPrice(ctx context.Context, arg UpdateStockPriceParams) error {
	_, err := q.db.ExecContext(ctx, updateStockPrice,
		arg.Price,
		arg.Trend,
		arg.PercentageChange,
		arg.ID,
	)
	return err
}
