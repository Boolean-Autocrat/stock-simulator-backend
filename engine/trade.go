package engine

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Trade struct {
	BuyerOrderID  uuid.UUID `json:"buyer_order_id"`
	BuyerID       uuid.UUID `json:"buyer_id"`
	SellerOrderID uuid.UUID `json:"seller_order_id"`
	SellerID      uuid.UUID `json:"seller_id"`
	Amount        int32     `json:"amount"`
	Price         float32   `json:"price"`
	Stock         uuid.UUID `json:"stock"`
}

func (trade *Trade) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, trade)
}

func (trade *Trade) ToJSON() []byte {
	str, _ := json.Marshal(trade)
	return str
}
