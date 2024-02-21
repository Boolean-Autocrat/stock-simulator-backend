package ipo

import (
	"log"
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

func (s *Service) RegisterHandlers(router *gin.RouterGroup) {
	router.POST("/ipo/buy", s.ipoBuy)
}

type ipoBuyRequest struct {
	StockID uuid.UUID `json:"stock" binding:"required"`
	Amount  int       `json:"amount" binding:"required"`
}

func (s *Service) ipoBuy(c *gin.Context) {
	userID := c.MustGet("userID").(uuid.UUID)
	var req ipoBuyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		log.Print(err.Error())
		c.JSON(400, gin.H{"message": "Invalid request."})
		return
	}
	stock, err := s.queries.GetStockById(c, req.StockID)
	if err != nil {
		log.Print(err.Error())
		c.JSON(400, gin.H{"message": "Invalid request."})
		return
	}
	if stock.IpoQuantity-stock.InCirculation < int32(req.Amount) {
		c.JSON(400, gin.H{"message": "Not enough stocks available for purchase"})
		return
	}
	userBalance, _ := s.queries.GetUserBalance(c, userID)
	if userBalance < float32(req.Amount)*stock.Price {
		c.JSON(400, gin.H{"message": "Not enough balance to purchase stock"})
		return
	}
	buyErr := s.queries.BuyStock(c, db.BuyStockParams{
		InCirculation: int32(req.Amount),
		ID:            req.StockID,
	})
	if buyErr != nil {
		log.Print(buyErr.Error())
		c.JSON(400, gin.H{"message": "Invalid request."})
		return
	}
	s.queries.BeginTransaction(c)
	addStockErr := s.queries.AddOrUpdateStockToPortfolio(c, db.AddOrUpdateStockToPortfolioParams{
		Stock:    req.StockID,
		User:     userID,
		Quantity: int32(req.Amount),
	})
	if addStockErr != nil {
		log.Print(addStockErr.Error())
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}
	balanceErr := s.queries.UpdateBalance(c, db.UpdateBalanceParams{
		ID:      userID,
		Balance: -(float32(req.Amount) * stock.Price),
	})
	if balanceErr != nil {
		log.Print(balanceErr.Error())
		c.JSON(400, gin.H{"message": "Invalid request."})
		return
	}
	s.queries.EndTransaction(c)
	addIpoHistoryErr := s.queries.AddToIpoHistory(c, db.AddToIpoHistoryParams{
		Stock:    req.StockID,
		User:     userID,
		Quantity: int32(req.Amount),
		Price:    stock.Price,
	})
	if addIpoHistoryErr != nil {
		log.Print(addIpoHistoryErr.Error())
		c.JSON(500, gin.H{"message": "Internal server error"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Stock purchased successfully",
	})
}
