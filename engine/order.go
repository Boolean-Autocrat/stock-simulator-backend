package engine

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Order struct {
	Amount uint64    `json:"amount"`
	Price  uint64    `json:"price"`
	ID     uuid.UUID `json:"id"`
	Side   int8      `json:"side"`
}

func (order *Order) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, order)
}

func (order *Order) ToJSON() []byte {
	str, _ := json.Marshal(order)
	return str
}
