// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: portfolio.sql

package db

import (
	"context"

	"github.com/google/uuid"
)

const addOrUpdateStockToPortfolio = `-- name: AddOrUpdateStockToPortfolio :exec
INSERT INTO portfolio ("user", "stock", quantity) VALUES ($1, $2, $3) ON CONFLICT ("user", "stock") DO UPDATE SET quantity = portfolio.quantity + $3
`

type AddOrUpdateStockToPortfolioParams struct {
	User     uuid.UUID `json:"user"`
	Stock    uuid.UUID `json:"stock"`
	Quantity int32     `json:"quantity"`
}

func (q *Queries) AddOrUpdateStockToPortfolio(ctx context.Context, arg AddOrUpdateStockToPortfolioParams) error {
	_, err := q.db.ExecContext(ctx, addOrUpdateStockToPortfolio, arg.User, arg.Stock, arg.Quantity)
	return err
}

const getPortfolio = `-- name: GetPortfolio :many
SELECT p."stock", s.name, s.symbol, s.price, s.is_crypto, s.is_stock, s.trend, s.percentage_change, p.quantity 
FROM portfolio p
JOIN stocks s ON p."stock" = s.id
WHERE p."user" = $1
`

type GetPortfolioRow struct {
	Stock            uuid.UUID `json:"stock"`
	Name             string    `json:"name"`
	Symbol           string    `json:"symbol"`
	Price            float32   `json:"price"`
	IsCrypto         bool      `json:"isCrypto"`
	IsStock          bool      `json:"isStock"`
	Trend            string    `json:"trend"`
	PercentageChange float32   `json:"percentageChange"`
	Quantity         int32     `json:"quantity"`
}

func (q *Queries) GetPortfolio(ctx context.Context, user uuid.UUID) ([]GetPortfolioRow, error) {
	rows, err := q.db.QueryContext(ctx, getPortfolio, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPortfolioRow
	for rows.Next() {
		var i GetPortfolioRow
		if err := rows.Scan(
			&i.Stock,
			&i.Name,
			&i.Symbol,
			&i.Price,
			&i.IsCrypto,
			&i.IsStock,
			&i.Trend,
			&i.PercentageChange,
			&i.Quantity,
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

const getStockWithQuantity = `-- name: GetStockWithQuantity :one
SELECT SUM(quantity) AS quantity, "stock" FROM portfolio WHERE "user" = $1 AND "stock" = $2 GROUP BY "stock"
`

type GetStockWithQuantityParams struct {
	User  uuid.UUID `json:"user"`
	Stock uuid.UUID `json:"stock"`
}

type GetStockWithQuantityRow struct {
	Quantity int64     `json:"quantity"`
	Stock    uuid.UUID `json:"stock"`
}

func (q *Queries) GetStockWithQuantity(ctx context.Context, arg GetStockWithQuantityParams) (GetStockWithQuantityRow, error) {
	row := q.db.QueryRowContext(ctx, getStockWithQuantity, arg.User, arg.Stock)
	var i GetStockWithQuantityRow
	err := row.Scan(&i.Quantity, &i.Stock)
	return i, err
}

const getStocksAndQuantity = `-- name: GetStocksAndQuantity :many
SELECT SUM(quantity) AS quantity, "stock" FROM portfolio WHERE "user" = $1 GROUP BY "stock"
`

type GetStocksAndQuantityRow struct {
	Quantity int64     `json:"quantity"`
	Stock    uuid.UUID `json:"stock"`
}

func (q *Queries) GetStocksAndQuantity(ctx context.Context, user uuid.UUID) ([]GetStocksAndQuantityRow, error) {
	rows, err := q.db.QueryContext(ctx, getStocksAndQuantity, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetStocksAndQuantityRow
	for rows.Next() {
		var i GetStocksAndQuantityRow
		if err := rows.Scan(&i.Quantity, &i.Stock); err != nil {
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
