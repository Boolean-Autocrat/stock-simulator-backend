// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package db

import (
	"time"

	"github.com/google/uuid"
)

type AccessToken struct {
	ID        uuid.UUID `json:"id"`
	User      uuid.UUID `json:"user"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"createdAt"`
}

type Developer struct {
	ID         int32  `json:"id"`
	Name       string `json:"name"`
	Title      string `json:"title"`
	Picture    string `json:"picture"`
	Email      string `json:"email"`
	GithubLink string `json:"githubLink"`
}

type IpoHistory struct {
	ID        uuid.UUID `json:"id"`
	User      uuid.UUID `json:"user"`
	Stock     uuid.UUID `json:"stock"`
	Quantity  int32     `json:"quantity"`
	Price     float32   `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
}

type News struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Author    string    `json:"author"`
	Content   string    `json:"content"`
	Tag       string    `json:"tag"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
}

type NewsSentiment struct {
	ID      uuid.UUID `json:"id"`
	Article uuid.UUID `json:"article"`
	User    uuid.UUID `json:"user"`
	Like    bool      `json:"like"`
	Dislike bool      `json:"dislike"`
}

type Order struct {
	ID                uuid.UUID `json:"id"`
	User              uuid.UUID `json:"user"`
	Stock             uuid.UUID `json:"stock"`
	Quantity          int32     `json:"quantity"`
	FulfilledQuantity int32     `json:"fulfilledQuantity"`
	Price             float32   `json:"price"`
	IsBuy             bool      `json:"isBuy"`
	CreatedAt         time.Time `json:"createdAt"`
}

type Portfolio struct {
	ID       uuid.UUID `json:"id"`
	User     uuid.UUID `json:"user"`
	Stock    uuid.UUID `json:"stock"`
	Quantity int32     `json:"quantity"`
}

type PriceHistory struct {
	ID      uuid.UUID `json:"id"`
	Stock   uuid.UUID `json:"stock"`
	Price   float32   `json:"price"`
	PriceAt time.Time `json:"priceAt"`
}

type Stock struct {
	ID               uuid.UUID `json:"id"`
	Name             string    `json:"name"`
	Symbol           string    `json:"symbol"`
	Price            float32   `json:"price"`
	IpoQuantity      int32     `json:"ipoQuantity"`
	InCirculation    int32     `json:"inCirculation"`
	IsStock          bool      `json:"isStock"`
	IsCrypto         bool      `json:"isCrypto"`
	Trend            string    `json:"trend"`
	PercentageChange float32   `json:"percentageChange"`
}

type TradeHistory struct {
	ID       uuid.UUID `json:"id"`
	Stock    uuid.UUID `json:"stock"`
	Quantity int32     `json:"quantity"`
	Price    float32   `json:"price"`
	Buyer    uuid.UUID `json:"buyer"`
	Seller   uuid.UUID `json:"seller"`
	TradedAt time.Time `json:"tradedAt"`
}

type User struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"fullName"`
	Email    string    `json:"email"`
	Picture  string    `json:"picture"`
	Balance  float32   `json:"balance"`
}

type Watchlist struct {
	ID      uuid.UUID `json:"id"`
	User    uuid.UUID `json:"user"`
	Stock   uuid.UUID `json:"stock"`
	AddedAt time.Time `json:"addedAt"`
}
