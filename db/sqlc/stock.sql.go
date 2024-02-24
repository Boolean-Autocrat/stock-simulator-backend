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

const addToIpoHistory = `-- name: AddToIpoHistory :exec
INSERT INTO ipo_history ("user", stock, quantity, price) VALUES ($1, $2, $3, $4)
`

type AddToIpoHistoryParams struct {
	User     uuid.UUID `json:"user"`
	Stock    uuid.UUID `json:"stock"`
	Quantity int32     `json:"quantity"`
	Price    float32   `json:"price"`
}

func (q *Queries) AddToIpoHistory(ctx context.Context, arg AddToIpoHistoryParams) error {
	_, err := q.db.ExecContext(ctx, addToIpoHistory,
		arg.User,
		arg.Stock,
		arg.Quantity,
		arg.Price,
	)
	return err
}

const addToWatchlist = `-- name: AddToWatchlist :exec
INSERT INTO watchlist ("user", stock) VALUES ($1, $2)
`

type AddToWatchlistParams struct {
	User  uuid.UUID `json:"user"`
	Stock uuid.UUID `json:"stock"`
}

func (q *Queries) AddToWatchlist(ctx context.Context, arg AddToWatchlistParams) error {
	_, err := q.db.ExecContext(ctx, addToWatchlist, arg.User, arg.Stock)
	return err
}

const buyStock = `-- name: BuyStock :exec
UPDATE stocks SET in_circulation = in_circulation + $1 WHERE id = $2
`

type BuyStockParams struct {
	InCirculation int32     `json:"inCirculation"`
	ID            uuid.UUID `json:"id"`
}

func (q *Queries) BuyStock(ctx context.Context, arg BuyStockParams) error {
	_, err := q.db.ExecContext(ctx, buyStock, arg.InCirculation, arg.ID)
	return err
}

const checkWatchlist = `-- name: CheckWatchlist :one
SELECT id, "user", stock, added_at FROM watchlist WHERE "user" = $1 AND stock = $2
`

type CheckWatchlistParams struct {
	User  uuid.UUID `json:"user"`
	Stock uuid.UUID `json:"stock"`
}

func (q *Queries) CheckWatchlist(ctx context.Context, arg CheckWatchlistParams) (Watchlist, error) {
	row := q.db.QueryRowContext(ctx, checkWatchlist, arg.User, arg.Stock)
	var i Watchlist
	err := row.Scan(
		&i.ID,
		&i.User,
		&i.Stock,
		&i.AddedAt,
	)
	return i, err
}

const createPriceHistory = `-- name: CreatePriceHistory :exec
INSERT INTO price_history (stock, price) VALUES ($1, $2)
`

type CreatePriceHistoryParams struct {
	Stock uuid.UUID `json:"stock"`
	Price float32   `json:"price"`
}

func (q *Queries) CreatePriceHistory(ctx context.Context, arg CreatePriceHistoryParams) error {
	_, err := q.db.ExecContext(ctx, createPriceHistory, arg.Stock, arg.Price)
	return err
}

const createStock = `-- name: CreateStock :one
INSERT INTO stocks (name, symbol, price, ipo_quantity, is_crypto, is_stock) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id, name, symbol, price, ipo_quantity, in_circulation, is_stock, is_crypto, trend, percentage_change
`

type CreateStockParams struct {
	Name        string  `json:"name"`
	Symbol      string  `json:"symbol"`
	Price       float32 `json:"price"`
	IpoQuantity int32   `json:"ipoQuantity"`
	IsCrypto    bool    `json:"isCrypto"`
	IsStock     bool    `json:"isStock"`
}

func (q *Queries) CreateStock(ctx context.Context, arg CreateStockParams) (Stock, error) {
	row := q.db.QueryRowContext(ctx, createStock,
		arg.Name,
		arg.Symbol,
		arg.Price,
		arg.IpoQuantity,
		arg.IsCrypto,
		arg.IsStock,
	)
	var i Stock
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Symbol,
		&i.Price,
		&i.IpoQuantity,
		&i.InCirculation,
		&i.IsStock,
		&i.IsCrypto,
		&i.Trend,
		&i.PercentageChange,
	)
	return i, err
}

const getIpoHistory = `-- name: GetIpoHistory :many
SELECT s.name, i.quantity, i.price, i.created_at FROM ipo_history i JOIN stocks s ON s.id = i.stock WHERE "user" = $1
`

type GetIpoHistoryRow struct {
	Name      string    `json:"name"`
	Quantity  int32     `json:"quantity"`
	Price     float32   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

func (q *Queries) GetIpoHistory(ctx context.Context, user uuid.UUID) ([]GetIpoHistoryRow, error) {
	rows, err := q.db.QueryContext(ctx, getIpoHistory, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetIpoHistoryRow
	for rows.Next() {
		var i GetIpoHistoryRow
		if err := rows.Scan(
			&i.Name,
			&i.Quantity,
			&i.Price,
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

const getStock = `-- name: GetStock :one
SELECT id, name, symbol, price, ipo_quantity, in_circulation, is_stock, is_crypto, trend, percentage_change FROM stocks WHERE name = $1 AND is_crypto = $2 AND is_stock = $3 AND symbol = $4 AND price = $5
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
		&i.IpoQuantity,
		&i.InCirculation,
		&i.IsStock,
		&i.IsCrypto,
		&i.Trend,
		&i.PercentageChange,
	)
	return i, err
}

const getStockById = `-- name: GetStockById :one
SELECT id, name, symbol, price, ipo_quantity, in_circulation, is_stock, is_crypto, trend, percentage_change FROM stocks WHERE id = $1
`

func (q *Queries) GetStockById(ctx context.Context, id uuid.UUID) (Stock, error) {
	row := q.db.QueryRowContext(ctx, getStockById, id)
	var i Stock
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Symbol,
		&i.Price,
		&i.IpoQuantity,
		&i.InCirculation,
		&i.IsStock,
		&i.IsCrypto,
		&i.Trend,
		&i.PercentageChange,
	)
	return i, err
}

const getStockPriceHistory = `-- name: GetStockPriceHistory :many
SELECT price, price_at FROM price_history WHERE stock = $1 ORDER BY price_at DESC
`

type GetStockPriceHistoryRow struct {
	Price   float32   `json:"price"`
	PriceAt time.Time `json:"priceAt"`
}

func (q *Queries) GetStockPriceHistory(ctx context.Context, stock uuid.UUID) ([]GetStockPriceHistoryRow, error) {
	rows, err := q.db.QueryContext(ctx, getStockPriceHistory, stock)
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
SELECT id, name, symbol, price, ipo_quantity, in_circulation, is_stock, is_crypto, trend, percentage_change FROM stocks
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
			&i.IpoQuantity,
			&i.InCirculation,
			&i.IsStock,
			&i.IsCrypto,
			&i.Trend,
			&i.PercentageChange,
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
SELECT stocks.id, stocks.name, stocks.symbol, stocks.price, stocks.ipo_quantity, stocks.is_crypto, stocks.is_stock, COUNT(price_history.id) AS price_history_count FROM stocks LEFT JOIN price_history ON stocks.id = price_history.stock GROUP BY stocks.id ORDER BY price_history_count DESC, stocks.name ASC LIMIT 10
`

type GetTrendingStocksRow struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Symbol            string    `json:"symbol"`
	Price             float32   `json:"price"`
	IpoQuantity       int32     `json:"ipoQuantity"`
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
			&i.IpoQuantity,
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

const getWatchlist = `-- name: GetWatchlist :many
SELECT stocks.name, stocks.symbol, stocks.price, stocks.is_crypto, stocks.is_stock, watchlist.added_at FROM stocks JOIN watchlist ON stocks.id = watchlist.stock WHERE watchlist."user" = $1
`

type GetWatchlistRow struct {
	Name     string    `json:"name"`
	Symbol   string    `json:"symbol"`
	Price    float32   `json:"price"`
	IsCrypto bool      `json:"isCrypto"`
	IsStock  bool      `json:"isStock"`
	AddedAt  time.Time `json:"addedAt"`
}

func (q *Queries) GetWatchlist(ctx context.Context, user uuid.UUID) ([]GetWatchlistRow, error) {
	rows, err := q.db.QueryContext(ctx, getWatchlist, user)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetWatchlistRow
	for rows.Next() {
		var i GetWatchlistRow
		if err := rows.Scan(
			&i.Name,
			&i.Symbol,
			&i.Price,
			&i.IsCrypto,
			&i.IsStock,
			&i.AddedAt,
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
SELECT id, name, symbol, price, ipo_quantity, in_circulation, is_stock, is_crypto, trend, percentage_change FROM stocks WHERE LOWER(name) LIKE '%' || LOWER($1) || '%'
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
			&i.IpoQuantity,
			&i.InCirculation,
			&i.IsStock,
			&i.IsCrypto,
			&i.Trend,
			&i.PercentageChange,
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
