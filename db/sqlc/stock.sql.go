// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: stock.sql

package db

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const createPriceHistory = `-- name: CreatePriceHistory :exec
INSERT INTO price_history (stock_id, price) VALUES ($1, $2)
`

type CreatePriceHistoryParams struct {
	StockID uuid.UUID `json:"stockId"`
	Price   float32   `json:"price"`
}

func (q *Queries) CreatePriceHistory(ctx context.Context, arg CreatePriceHistoryParams) error {
	_, err := q.db.ExecContext(ctx, createPriceHistory, arg.StockID, arg.Price)
	return err
}

const createStock = `-- name: CreateStock :one
INSERT INTO stocks (name, symbol, price, quantity, is_crypto, is_stock) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, symbol, price, is_crypto, is_stock, quantity, trend, percent_change
`

type CreateStockParams struct {
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Price    float32 `json:"price"`
	Quantity int32   `json:"quantity"`
	IsCrypto bool    `json:"isCrypto"`
	IsStock  bool    `json:"isStock"`
}

func (q *Queries) CreateStock(ctx context.Context, arg CreateStockParams) (Stock, error) {
	row := q.db.QueryRowContext(ctx, createStock,
		arg.Name,
		arg.Symbol,
		arg.Price,
		arg.Quantity,
		arg.IsCrypto,
		arg.IsStock,
	)
	var i Stock
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Symbol,
		&i.Price,
		&i.IsCrypto,
		&i.IsStock,
		&i.Quantity,
		&i.Trend,
		&i.PercentChange,
	)
	return i, err
}

const getStock = `-- name: GetStock :one
SELECT id, name, symbol, price, is_crypto, is_stock, quantity, trend, percent_change FROM stocks WHERE name = $1 AND is_crypto = $2 AND is_stock = $3 AND symbol = $4 AND price = $5
`

type GetStockParams struct {
	Name     string  `json:"name"`
	IsCrypto bool    `json:"isCrypto"`
	IsStock  bool    `json:"isStock"`
	Symbol   string  `json:"symbol"`
	Price    float32 `json:"price"`
}

func (q *Queries) GetStock(ctx context.Context, arg GetStockParams) (Stock, error) {
	row := q.db.QueryRowContext(ctx, getStock,
		arg.Name,
		arg.IsCrypto,
		arg.IsStock,
		arg.Symbol,
		arg.Price,
	)
	var i Stock
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Symbol,
		&i.Price,
		&i.IsCrypto,
		&i.IsStock,
		&i.Quantity,
		&i.Trend,
		&i.PercentChange,
	)
	return i, err
}

const getStockById = `-- name: GetStockById :one
SELECT id, name, symbol, price, is_crypto, is_stock, quantity, trend, percent_change FROM stocks WHERE id = $1
`

func (q *Queries) GetStockById(ctx context.Context, id uuid.UUID) (Stock, error) {
	row := q.db.QueryRowContext(ctx, getStockById, id)
	var i Stock
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Symbol,
		&i.Price,
		&i.IsCrypto,
		&i.IsStock,
		&i.Quantity,
		&i.Trend,
		&i.PercentChange,
	)
	return i, err
}

const getStockPriceHistory = `-- name: GetStockPriceHistory :many
SELECT price, price_at FROM price_history WHERE stock_id = $1 ORDER BY price_at DESC
`

type GetStockPriceHistoryRow struct {
	Price   float32   `json:"price"`
	PriceAt time.Time `json:"priceAt"`
}

func (q *Queries) GetStockPriceHistory(ctx context.Context, stockID uuid.UUID) ([]GetStockPriceHistoryRow, error) {
	rows, err := q.db.QueryContext(ctx, getStockPriceHistory, stockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetStockPriceHistoryRow
	for rows.Next() {
		var i GetStockPriceHistoryRow
		if err := rows.Scan(&i.Price, &i.PriceAt); err != nil {
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

const getStocks = `-- name: GetStocks :many
SELECT id, name, symbol, price, is_crypto, is_stock, quantity, trend, percent_change FROM stocks
`

func (q *Queries) GetStocks(ctx context.Context) ([]Stock, error) {
	rows, err := q.db.QueryContext(ctx, getStocks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Stock
	for rows.Next() {
		var i Stock
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Symbol,
			&i.Price,
			&i.IsCrypto,
			&i.IsStock,
			&i.Quantity,
			&i.Trend,
			&i.PercentChange,
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

const getTrendingStocks = `-- name: GetTrendingStocks :many
SELECT stocks.id, stocks.name, stocks.symbol, stocks.price, stocks.quantity, stocks.is_crypto, stocks.is_stock, COUNT(price_history.id) AS price_history_count FROM stocks LEFT JOIN price_history ON stocks.id = price_history.stock_id GROUP BY stocks.id ORDER BY price_history_count DESC, stocks.name ASC LIMIT 10
`

type GetTrendingStocksRow struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Symbol            string    `json:"symbol"`
	Price             float32   `json:"price"`
	Quantity          int32     `json:"quantity"`
	IsCrypto          bool      `json:"isCrypto"`
	IsStock           bool      `json:"isStock"`
	PriceHistoryCount int64     `json:"priceHistoryCount"`
}

func (q *Queries) GetTrendingStocks(ctx context.Context) ([]GetTrendingStocksRow, error) {
	rows, err := q.db.QueryContext(ctx, getTrendingStocks)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetTrendingStocksRow
	for rows.Next() {
		var i GetTrendingStocksRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Symbol,
			&i.Price,
			&i.Quantity,
			&i.IsCrypto,
			&i.IsStock,
			&i.PriceHistoryCount,
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

const searchStocks = `-- name: SearchStocks :many
SELECT id, name, symbol, price, is_crypto, is_stock, quantity, trend, percent_change FROM stocks WHERE LOWER(name) LIKE '%' || LOWER($1) || '%'
`

func (q *Queries) SearchStocks(ctx context.Context, lower string) ([]Stock, error) {
	rows, err := q.db.QueryContext(ctx, searchStocks, lower)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Stock
	for rows.Next() {
		var i Stock
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Symbol,
			&i.Price,
			&i.IsCrypto,
			&i.IsStock,
			&i.Quantity,
			&i.Trend,
			&i.PercentChange,
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
