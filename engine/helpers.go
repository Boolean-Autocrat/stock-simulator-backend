package engine

import (
	"context"
	"log"

	db "github.com/Boolean-Autocrat/stock-simulator-backend/db/sqlc"
	amqp "github.com/rabbitmq/amqp091-go"
)

func AddPendingOrders(queries *db.Queries, channelMQ *amqp.Channel) {
	pendingOrders, err := queries.GetUnfulfilledOrders(context.Background())
	if err != nil {
		panic(err)
	}
	for _, order := range pendingOrders {
		var side int8
		if order.IsBuy {
			side = 1
		} else {
			side = 0
		}
		var parsedOrder = Order{
			OrderID: order.ID,
			UserID:  order.User,
			Stock:   order.Stock,
			Amount:  order.Quantity - order.FulfilledQuantity,
			Price:   order.Price,
			Side:    side,
		}
		message := amqp.Publishing{
			ContentType: "text/plain",
			Body:        parsedOrder.ToJSON(),
		}
		if err := channelMQ.PublishWithContext(
			context.Background(), // context
			"",                   // exchange
			"orders",             // queue name
			false,                // mandatory
			false,                // immediate
			message,              // message to publish
		); err != nil {
			log.Println(err)
			panic(err)
		}
	}
}
