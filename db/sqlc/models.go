// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type News struct {
	ID          uuid.UUID      `json:"id"`
	Title       sql.NullString `json:"title"`
	Description sql.NullString `json:"description"`
	Photo       sql.NullString `json:"photo"`
	Likes       sql.NullInt32  `json:"likes"`
	Dislikes    sql.NullInt32  `json:"dislikes"`
}

type NewsSentiment struct {
	ID        uuid.UUID     `json:"id"`
	ArticleID uuid.NullUUID `json:"articleId"`
	UserID    uuid.NullUUID `json:"userId"`
	Like      sql.NullBool  `json:"like"`
	Dislike   sql.NullBool  `json:"dislike"`
}

type Portfolio struct {
	ID            uuid.UUID      `json:"id"`
	UserID        uuid.NullUUID  `json:"userId"`
	StockID       uuid.NullUUID  `json:"stockId"`
	PurchasePrice sql.NullString `json:"purchasePrice"`
	PurchasedAt   time.Time      `json:"purchasedAt"`
}

type PriceHistory struct {
	ID      uuid.UUID      `json:"id"`
	StockID uuid.NullUUID  `json:"stockId"`
	Price   sql.NullString `json:"price"`
	PriceAt time.Time      `json:"priceAt"`
}

type Stock struct {
	ID       uuid.UUID      `json:"id"`
	Name     sql.NullString `json:"name"`
	Symbol   sql.NullString `json:"symbol"`
	Price    sql.NullString `json:"price"`
	IsCrypto sql.NullBool   `json:"isCrypto"`
	IsStock  sql.NullBool   `json:"isStock"`
}

type User struct {
	ID       uuid.UUID      `json:"id"`
	FullName sql.NullString `json:"fullName"`
	Email    sql.NullString `json:"email"`
	Picture  sql.NullString `json:"picture"`
	Balance  sql.NullString `json:"balance"`
}
