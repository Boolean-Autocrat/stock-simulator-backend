// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: stock.sql

package db

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPriceHistory = `-- name: CreatePriceHistory :one
INSERT INTO price_history (stock_id, price) VALUES ($1, $2) RETURNING id, stock_id, price, price_at
`

type CreatePriceHistoryParams struct {
	StockID uuid.NullUUID  `json:"stockId"`
	Price   sql.NullString `json:"price"`
}

func (q *Queries) CreatePriceHistory(ctx context.Context, arg CreatePriceHistoryParams) (PriceHistory, error) {
	row := q.db.QueryRowContext(ctx, createPriceHistory, arg.StockID, arg.Price)
	var i PriceHistory
	err := row.Scan(
		&i.ID,
		&i.StockID,
		&i.Price,
		&i.PriceAt,
	)
	return i, err
}

const createStock = `-- name: CreateStock :one
INSERT INTO stocks (name, symbol, price, is_crypto, is_stock) VALUES ($1, $2, $3, $4, $5) RETURNING id, name, symbol, price, is_crypto, is_stock
`

type CreateStockParams struct {
	Name     sql.NullString `json:"name"`
	Symbol   sql.NullString `json:"symbol"`
	Price    sql.NullString `json:"price"`
	IsCrypto sql.NullBool   `json:"isCrypto"`
	IsStock  sql.NullBool   `json:"isStock"`
}

func (q *Queries) CreateStock(ctx context.Context, arg CreateStockParams) (Stock, error) {
	row := q.db.QueryRowContext(ctx, createStock,
		arg.Name,
		arg.Symbol,
		arg.Price,
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
	)
	return i, err
}

const getStock = `-- name: GetStock :one
SELECT id, name, symbol, price, is_crypto, is_stock FROM stocks WHERE name = $1 AND is_crypto = $2 AND is_stock = $3 AND symbol = $4 AND price = $5
`

type GetStockParams struct {
	Name     sql.NullString `json:"name"`
	IsCrypto sql.NullBool   `json:"isCrypto"`
	IsStock  sql.NullBool   `json:"isStock"`
	Symbol   sql.NullString `json:"symbol"`
	Price    sql.NullString `json:"price"`
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
	)
	return i, err
}

const getStockPriceHistory = `-- name: GetStockPriceHistory :many
SELECT id, stock_id, price, price_at FROM price_history WHERE stock_id = $1
`

func (q *Queries) GetStockPriceHistory(ctx context.Context, stockID uuid.NullUUID) ([]PriceHistory, error) {
	rows, err := q.db.QueryContext(ctx, getStockPriceHistory, stockID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PriceHistory
	for rows.Next() {
		var i PriceHistory
		if err := rows.Scan(
			&i.ID,
			&i.StockID,
			&i.Price,
			&i.PriceAt,
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

const getStockPriceHistoryByDate = `-- name: GetStockPriceHistoryByDate :many
SELECT id, stock_id, price, price_at FROM price_history WHERE stock_id = $1 AND price_at >= $2 AND price_at <= $3
`

type GetStockPriceHistoryByDateParams struct {
	StockID   uuid.NullUUID `json:"stockId"`
	PriceAt   time.Time     `json:"priceAt"`
	PriceAt_2 time.Time     `json:"priceAt2"`
}

func (q *Queries) GetStockPriceHistoryByDate(ctx context.Context, arg GetStockPriceHistoryByDateParams) ([]PriceHistory, error) {
	rows, err := q.db.QueryContext(ctx, getStockPriceHistoryByDate, arg.StockID, arg.PriceAt, arg.PriceAt_2)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PriceHistory
	for rows.Next() {
		var i PriceHistory
		if err := rows.Scan(
			&i.ID,
			&i.StockID,
			&i.Price,
			&i.PriceAt,
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

const getStocks = `-- name: GetStocks :many
SELECT id, name, symbol, price, is_crypto, is_stock FROM stocks WHERE is_crypto = $1 AND is_stock = $2
`

type GetStocksParams struct {
	IsCrypto sql.NullBool `json:"isCrypto"`
	IsStock  sql.NullBool `json:"isStock"`
}

func (q *Queries) GetStocks(ctx context.Context, arg GetStocksParams) ([]Stock, error) {
	rows, err := q.db.QueryContext(ctx, getStocks, arg.IsCrypto, arg.IsStock)
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
SELECT id, name, symbol, price, is_crypto, is_stock FROM stocks WHERE LOWER(name) LIKE '%' || LOWER($1) || '%'
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
