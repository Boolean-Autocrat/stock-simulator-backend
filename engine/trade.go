package engine

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Trade struct {
	TakerOrderID uuid.UUID `json:"taker_order_id"`
	MakerOrderID uuid.UUID `json:"maker_order_id"`
	Amount       uint64    `json:"amount"`
	Price        uint64    `json:"price"`
}

func (trade *Trade) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, trade)
}

func (trade *Trade) ToJSON() []byte {
	str, _ := json.Marshal(trade)
	return str
}
