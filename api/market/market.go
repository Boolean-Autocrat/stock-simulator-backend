package market

import (
	"database/sql"
	"log"

	db "github.com/Boolean-Autocrat/stock-simulator-backend/db/sqlc"
	"github.com/Boolean-Autocrat/stock-simulator-backend/engine"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Service struct {
	queries   *db.Queries
	channelMQ *amqp.Channel
}

func NewService(queries *db.Queries, channelMQ *amqp.Channel) *Service {
	return &Service{queries: queries, channelMQ: channelMQ}
}

func (s *Service) RegisterHandlers(router *gin.RouterGroup) {
	router.POST("/market/sell", s.sellAsset)
	router.POST("/market/buy", s.buyAsset)
	router.GET("/market/status", s.GetOrderStats)
}

type Order struct {
	Stock  uuid.UUID `json:"stock" binding:"required"`
	Amount int32     `json:"amount" binding:"required"`
	Price  float32   `json:"price" binding:"required"`
}

func (s *Service) sellAsset(c *gin.Context) {
	userID, _ := c.Get("userID")
	var order *Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	portfolio, err := s.queries.GetPortfolio(c, userID.(uuid.UUID))
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(400, gin.H{"error": "You don't have enough of this stock to sell."})
			return
		}
		log.Println(err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	var stockExists bool
	for _, stock := range portfolio {
		if stock.Stock == order.Stock {
			if stock.Quantity < order.Amount {
				c.JSON(400, gin.H{"error": "You don't have enough of this stock to sell."})
				return
			} else {
				stockExists = true
				break
			}
		}
	}
	if !stockExists {
		c.JSON(400, gin.H{"error": "You don't have enough of this stock to sell."})
		return
	}
	s.queries.BeginTransaction(c)
	updatePortfolioErr := s.queries.AddOrUpdateStockToPortfolio(c, db.AddOrUpdateStockToPortfolioParams{
		User:     userID.(uuid.UUID),
		Stock:    order.Stock,
		Quantity: -order.Amount,
	})
	if updatePortfolioErr != nil {
		log.Println(updatePortfolioErr)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	s.queries.EndTransaction(c)
	pendingOrder, err := s.queries.CreatePendingOrder(c, db.CreatePendingOrderParams{
		User:     userID.(uuid.UUID),
		Stock:    order.Stock,
		Quantity: order.Amount,
		Price:    order.Price,
		IsBuy:    false,
	})
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	var parsedOrder = engine.Order{
		OrderID: pendingOrder,
		UserID:  userID.(uuid.UUID),
		Stock:   order.Stock,
		Amount:  order.Amount,
		Price:   order.Price,
		Side:    0,
	}
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        parsedOrder.ToJSON(),
	}
	if err := s.channelMQ.PublishWithContext(
		c,        // context
		"",       // exchange
		"orders", // queue name
		false,    // mandatory
		false,    // immediate
		message,  // message to publish
	); err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"message": "Order placed successfully"})
}

func (s *Service) buyAsset(c *gin.Context) {
	userID, _ := c.Get("userID")
	var order *Order
	err := c.ShouldBindJSON(&order)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	stock, err := s.queries.GetStockById(c, order.Stock)
	if err != nil {
		c.JSON(400, gin.H{"error": "Stock does not exist"})
		return
	}
	balance, _ := s.queries.GetUserBalance(c, userID.(uuid.UUID))
	if balance < float32(order.Amount)*order.Price {
		c.JSON(400, gin.H{"error": "Not enough balance to place this stock order."})
		return
	}
	s.queries.BeginTransaction(c)
	updateBalanceErr := s.queries.UpdateBalance(c, db.UpdateBalanceParams{
		ID:      userID.(uuid.UUID),
		Balance: -float32(order.Amount) * order.Price,
	})
	if updateBalanceErr != nil {
		log.Println(updateBalanceErr)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	s.queries.EndTransaction(c)
	pendingOrder, err := s.queries.CreatePendingOrder(c, db.CreatePendingOrderParams{
		User:     userID.(uuid.UUID),
		Stock:    order.Stock,
		Quantity: order.Amount,
		Price:    order.Price,
		IsBuy:    true,
	})
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	var parsedOrder = engine.Order{
		OrderID: pendingOrder,
		UserID:  userID.(uuid.UUID),
		Stock:   stock.ID,
		Amount:  order.Amount,
		Price:   order.Price,
		Side:    1,
	}
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        parsedOrder.ToJSON(),
	}
	if err := s.channelMQ.PublishWithContext(
		c,        // context
		"",       // exchange
		"orders", // queue name
		false,    // mandatory
		false,    // immediate
		message,  // message to publish
	); err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"message": "Order placed successfully"})
}

func (s *Service) GetOrderStats(c *gin.Context) {
	userID, _ := c.Get("userID")
	pendingOrders, err := s.queries.GetPendingOrders(c, userID.(uuid.UUID))
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	var pendingOrdersResponse []gin.H
	for _, order := range pendingOrders {
		var orderType string
		if order.IsBuy {
			orderType = "buy"
		} else {
			orderType = "sell"
		}
		var orderStatus string
		if order.FulfilledQuantity == order.Quantity {
			orderStatus = "fulfilled"
		} else {
			orderStatus = "pending"
		}
		pendingOrdersResponse = append(pendingOrdersResponse, gin.H{
			"stock":      order.Stock,
			"quantity":   order.Quantity,
			"price":      order.Price,
			"type":       orderType,
			"fulfilled":  order.FulfilledQuantity,
			"status":     orderStatus,
			"created_at": order.CreatedAt,
		})
	}
	ipoHistory, err := s.queries.GetIpoHistory(c, userID.(uuid.UUID))
	if err != nil {
		log.Println(err)
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.JSON(200, gin.H{"orders": pendingOrdersResponse, "ipo": ipoHistory})
}
