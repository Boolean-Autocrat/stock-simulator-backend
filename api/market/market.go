package market

import (
	db "github.com/Boolean-Autocrat/stock-simulator-backend/db/sqlc"
	"github.com/Boolean-Autocrat/stock-simulator-backend/engine"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service struct {
	queries  *db.Queries
	producer *kafka.Producer
}

func NewService(queries *db.Queries, producer *kafka.Producer) *Service {
	return &Service{queries: queries, producer: producer}
}

func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.POST("/market/sell", s.sellAsset)
	router.POST("/market/buy", s.buyAsset)
	router.GET("/market/status", s.GetOrderStats)
}

type Order struct {
	Stock  uuid.UUID `json:"stock" binding:"required"`
	Amount uint64    `json:"amount" binding:"required"`
	Price  float32   `json:"price" binding:"required"`
}

func (s *Service) sellAsset(c *gin.Context) {
	// TODO: Check if user has enough stock to sell
	// TODO: Check if stock exists
	userID, _ := c.Get("userID")
	var order *Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var parsedOrder = engine.Order{
		UserID: userID.(uuid.UUID),
		Stock:  order.Stock,
		Amount: order.Amount,
		Price:  order.Price,
		Side:   1,
	}
	var tradesTopic = "orders"
	s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &tradesTopic,
			Partition: kafka.PartitionAny,
		},
		Value: parsedOrder.ToJSON(),
	}, nil)
	c.JSON(200, gin.H{"message": "Order placed successfully"})
}

func (s *Service) buyAsset(c *gin.Context) {
	// TODO: Check if user has enough balance
	// TODO: Check if stock exists
	userID, _ := c.Get("userID")
	var order *Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	var parsedOrder = engine.Order{
		UserID: userID.(uuid.UUID),
		Stock:  order.Stock,
		Amount: order.Amount,
		Price:  order.Price,
		Side:   0,
	}
	var tradesTopic = "orders"
	s.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &tradesTopic,
			Partition: kafka.PartitionAny,
		},
		Value: parsedOrder.ToJSON(),
	}, nil)
	c.JSON(200, gin.H{"message": "Order placed successfully"})
}

func (s *Service) GetOrderStats(c *gin.Context) {
}
