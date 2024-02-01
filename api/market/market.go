package market

import (
	"net/http"

	db "github.com/Boolean-Autocrat/stock-simulator-backend/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.POST("/market/sell", s.sellAsset)
	router.POST("/market/buy", s.buyAsset)
	router.GET("/market/status", s.GetOrderStats)
}

type Order struct {
	Stock    uuid.UUID `json:"stock"`
	Quantity int       `json:"quantity"`
	Price    float32   `json:"price"`
}

func (s *Service) sellAsset(c *gin.Context) {
	userId, _ := c.Get("userID")
	sellID, err := uuid.NewUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var req Order
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON payload!"})
		return
	}
	userStock, err := s.queries.GetStockWithQuantity(c, db.GetStockWithQuantityParams{
		UserID:  userId.(uuid.UUID),
		StockID: req.Stock,
	})
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"message": "You do not own enough of this stock!"})
		return
	}
	if int32(userStock.Quantity) < int32(req.Quantity) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Insufficient stock quantity to place order!"})
		return
	}
	buyOrders, err := s.queries.ListBuyOrders(c, req.Stock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	n := len(buyOrders)
	if n != 0 && buyOrders[n-1].Price <= req.Price {
		for i := n - 1; i >= 0; i-- {
			buyOrder := buyOrders[i]
			if buyOrder.Price < req.Price {
				break
			}
			if buyOrder.Quantity-buyOrder.Fulfilled >= int32(req.Quantity) {
				tradeID, err := uuid.NewUUID()
				if err != nil {
					c.JSON(http.StatusInternalServerError, err)
					return
				}
				s.queries.AddTrade(c, db.AddTradeParams{
					ID:        tradeID,
					BuyOrder:  buyOrder.ID,
					SellOrder: sellID,
					Quantity:  int32(req.Quantity),
					Price:     req.Price,
				})
				s.queries.UpdateBuyOrder(c, db.UpdateBuyOrderParams{
					ID:        buyOrder.ID,
					Fulfilled: buyOrder.Fulfilled + int32(req.Quantity),
				})
				s.queries.AddSellOrder(c, db.AddSellOrderParams{
					ID:        sellID,
					User:      userId.(uuid.UUID),
					Stock:     req.Stock,
					Price:     req.Price,
					Quantity:  int32(req.Quantity),
					Fulfilled: int32(req.Quantity),
				})
				s.queries.UpdateBalance(c, db.UpdateBalanceParams{
					ID:      userId.(uuid.UUID),
					Balance: req.Price * float32(req.Quantity),
				})
				c.JSON(http.StatusOK, gin.H{"message": "Order successfully processed and matched!"})
				return
			}
			if buyOrder.Quantity-buyOrder.Fulfilled < int32(req.Quantity) {
				tradeID, err := uuid.NewUUID()
				if err != nil {
					c.JSON(http.StatusInternalServerError, err)
					return
				}
				s.queries.AddTrade(c, db.AddTradeParams{
					ID:        tradeID,
					BuyOrder:  buyOrder.ID,
					SellOrder: sellID,
					Quantity:  buyOrder.Quantity - buyOrder.Fulfilled,
					Price:     req.Price,
				})
				s.queries.UpdateBuyOrder(c, db.UpdateBuyOrderParams{
					ID:        buyOrder.ID,
					Fulfilled: buyOrder.Quantity,
				})
				s.queries.UpdateBalance(c, db.UpdateBalanceParams{
					ID:      userId.(uuid.UUID),
					Balance: req.Price * float32(buyOrder.Quantity-buyOrder.Fulfilled),
				})
				req.Quantity -= int(buyOrder.Quantity - buyOrder.Fulfilled)
				continue
			}
		}
	}
	s.queries.AddSellOrder(c, db.AddSellOrderParams{
		ID:        sellID,
		User:      userId.(uuid.UUID),
		Price:     req.Price,
		Stock:     req.Stock,
		Quantity:  int32(req.Quantity),
		Fulfilled: 0,
	})
	if req.Quantity != 0 {
		c.JSON(http.StatusOK, gin.H{"message": "Order added to queue successfully!"})
		return
	} else {
		c.JSON(http.StatusOK, gin.H{"message": "Order successfully processed and matched!"})
	}
}

func (s *Service) buyAsset(c *gin.Context) {
	userId, _ := c.Get("userID")
	userBalance, err := s.queries.GetUserBalance(c, userId.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	buyID, err := uuid.NewUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	var req Order
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid JSON payload!"})
		return
	}
	if userBalance < req.Price*float32(req.Quantity) {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Insufficient balance to place order!"})
		return
	}
	sellOrders, err := s.queries.ListSellOrders(c, req.Stock)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	s.queries.UpdateBalance(c, db.UpdateBalanceParams{
		ID:      userId.(uuid.UUID),
		Balance: -req.Price * float32(req.Quantity),
	})
	n := len(sellOrders)
	if n != 0 && sellOrders[n-1].Price >= req.Price {
		for i := n - 1; i >= 0; i-- {
			sellOrder := sellOrders[i]
			if sellOrder.Price > req.Price {
				break
			}
			if sellOrder.Quantity-sellOrder.Fulfilled >= int32(req.Quantity) {
				tradeID, err := uuid.NewUUID()
				if err != nil {
					c.JSON(http.StatusInternalServerError, err)
					return
				}
				s.queries.AddTrade(c, db.AddTradeParams{
					ID:        tradeID,
					BuyOrder:  buyID,
					SellOrder: sellOrder.ID,
					Quantity:  int32(req.Quantity),
					Price:     req.Price,
				})
				s.queries.UpdateSellOrder(c, db.UpdateSellOrderParams{
					ID:        sellOrder.ID,
					Fulfilled: sellOrder.Fulfilled + int32(req.Quantity),
				})
				s.queries.AddBuyOrder(c, db.AddBuyOrderParams{
					ID:        buyID,
					User:      userId.(uuid.UUID),
					Price:     req.Price,
					Stock:     req.Stock,
					Quantity:  int32(req.Quantity),
					Fulfilled: int32(req.Quantity),
				})
				s.queries.UpdateBalance(c, db.UpdateBalanceParams{
					ID:      sellOrder.User,
					Balance: req.Price * float32(req.Quantity),
				})
				c.JSON(http.StatusOK, gin.H{"message": "Order successfully processed and matched!"})
				return
			}
			if sellOrder.Quantity-sellOrder.Fulfilled < int32(req.Quantity) {
				tradeID, err := uuid.NewUUID()
				if err != nil {
					c.JSON(http.StatusInternalServerError, err)
					return
				}
				s.queries.AddTrade(c, db.AddTradeParams{
					ID:        tradeID,
					BuyOrder:  buyID,
					SellOrder: sellOrder.ID,
					Quantity:  sellOrder.Quantity - sellOrder.Fulfilled,
					Price:     req.Price,
				})
				s.queries.UpdateSellOrder(c, db.UpdateSellOrderParams{
					ID:        sellOrder.ID,
					Fulfilled: sellOrder.Quantity,
				})
				s.queries.UpdateBalance(c, db.UpdateBalanceParams{
					ID:      sellOrder.User,
					Balance: req.Price * float32(sellOrder.Quantity-sellOrder.Fulfilled),
				})
				req.Quantity -= int(sellOrder.Quantity - sellOrder.Fulfilled)
				continue
			}
		}
	}
	s.queries.AddBuyOrder(c, db.AddBuyOrderParams{
		ID:        buyID,
		User:      userId.(uuid.UUID),
		Stock:     req.Stock,
		Price:     req.Price,
		Quantity:  int32(req.Quantity),
		Fulfilled: 0,
	})
	c.JSON(http.StatusOK, gin.H{"message": "Order added to orderbook successfully!"})
}

func (s *Service) GetOrderStats(c *gin.Context) {
	userId, _ := c.Get("userID")
	buyOrders, err := s.queries.GetAllBuyOrdersByUser(c, userId.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	sellOrders, err := s.queries.GetAllSellOrdersByUser(c, userId.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"buyOrders": buyOrders, "sellOrders": sellOrders})
}
