package engine

import (
	"encoding/json"

	"github.com/google/uuid"
)

type Order struct {
	Amount uint64    `json:"amount" binding:"required"`
	Price  float32   `json:"price" binding:"required"`
	UserID uuid.UUID `json:"id" binding:"required"`
	Side   int8      `json:"side" binding:"required"`
	Stock  uuid.UUID `json:"stock" binding:"required"`
}

func (order *Order) FromJSON(msg []byte) error {
	return json.Unmarshal(msg, order)
}

func (order *Order) ToJSON() []byte {
	str, _ := json.Marshal(order)
	return str
}
