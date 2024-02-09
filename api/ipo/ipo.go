package ipo

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
	router.POST("/ipo/buy", s.ipoBuy)
}

type ipoBuyRequest struct {
	StockID uuid.UUID `json:"stock" binding:"required"`
	Amount  int       `json:"amount" binding:"required"`
}

func (s *Service) ipoBuy(c *gin.Context) {
	var req ipoBuyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	stock, err := s.queries.GetStockById(c, req.StockID)
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	if stock.Quantity-stock.InCirculation < int32(req.Amount) {
		c.JSON(http.StatusBadRequest, "Not enough stock available")
		return
	}
	buyErr := s.queries.BuyStock(c, db.BuyStockParams{
		InCirculation: stock.Quantity,
		ID:            req.StockID,
	})
	if buyErr != nil {
		c.JSON(http.StatusBadRequest, buyErr.Error())
		return
	}
	c.JSON(http.StatusOK, "Stock purchased successfully")
}
