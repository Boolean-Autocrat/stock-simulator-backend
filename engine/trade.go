package engine

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Trade struct {
	BuyerOrderID  uuid.UUID `json:"buyer_order_id" binding:"required"`
	BuyerID       uuid.UUID `json:"buyer_id" binding:"required"`
	SellerOrderID uuid.UUID `json:"seller_order_id" binding:"required"`
	SellerID      uuid.UUID `json:"seller_id" binding:"required"`
	Amount        int32     `json:"amount" binding:"required"`
	Price         float32   `json:"price" binding:"required"`
	Stock         uuid.UUID `json:"stock" binding:"required"`
}

func (trade *Trade) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, trade)
}

func (trade *Trade) ToJSON() []byte {
	str, _ := json.Marshal(trade)
	return str
}
