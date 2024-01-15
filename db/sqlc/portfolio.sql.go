// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: portfolio.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const addStockToPortfolio = `-- name: AddStockToPortfolio :one
INSERT INTO portfolio (user_id, stock_id, purchase_price) VALUES ($1, $2, $3) RETURNING id, user_id, stock_id, purchase_price, purchased_at
`

type AddStockToPortfolioParams struct {
	UserID        uuid.UUID `json:"userId"`
	StockID       uuid.UUID `json:"stockId"`
	PurchasePrice string    `json:"purchasePrice"`
}

func (q *Queries) AddStockToPortfolio(ctx context.Context, arg AddStockToPortfolioParams) (Portfolio, error) {
	row := q.db.QueryRowContext(ctx, addStockToPortfolio, arg.UserID, arg.StockID, arg.PurchasePrice)
	var i Portfolio
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.StockID,
		&i.PurchasePrice,
		&i.PurchasedAt,
	)
	return i, err
}

const getPortfolio = `-- name: GetPortfolio :many
SELECT p.stock_id, p.purchase_price, p.purchased_at, s.name, s.symbol, s.price, s.is_crypto, s.is_stock
FROM portfolio p
JOIN stocks s ON p.stock_id = s.id
WHERE p.user_id = $1
ORDER BY p.purchased_at
LIMIT 10
OFFSET $2
`

type GetPortfolioParams struct {
	UserID uuid.UUID `json:"userId"`
	Offset int32     `json:"offset"`
}

type GetPortfolioRow struct {
	StockID       uuid.UUID `json:"stockId"`
	PurchasePrice string    `json:"purchasePrice"`
	PurchasedAt   time.Time `json:"purchasedAt"`
	Name          string    `json:"name"`
	Symbol        string    `json:"symbol"`
	Price         string    `json:"price"`
	IsCrypto      bool      `json:"isCrypto"`
	IsStock       bool      `json:"isStock"`
}

func (q *Queries) GetPortfolio(ctx context.Context, arg GetPortfolioParams) ([]GetPortfolioRow, error) {
	rows, err := q.db.QueryContext(ctx, getPortfolio, arg.UserID, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetPortfolioRow
	for rows.Next() {
		var i GetPortfolioRow
		if err := rows.Scan(
			&i.StockID,
			&i.PurchasePrice,
			&i.PurchasedAt,
			&i.Name,
			&i.Symbol,
			&i.Price,
			&i.IsCrypto,
			&i.IsStock,
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

const removeStockFromPortfolio = `-- name: RemoveStockFromPortfolio :exec
DELETE FROM portfolio WHERE user_id = $1 AND stock_id = $2
`

type RemoveStockFromPortfolioParams struct {
	UserID  uuid.UUID `json:"userId"`
	StockID uuid.UUID `json:"stockId"`
}

func (q *Queries) RemoveStockFromPortfolio(ctx context.Context, arg RemoveStockFromPortfolioParams) error {
	_, err := q.db.ExecContext(ctx, removeStockFromPortfolio, arg.UserID, arg.StockID)
	return err
}
