package stocks

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

func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.GET("/stocks", s.GetStocks)
	router.GET("/stocks/:id", s.GetStock)
	router.GET("/stocks/trending", s.GetTrendingStocks)
	router.GET("/stocks/search", s.SearchStocks)
	router.GET("/stocks/:id/price_history", s.GetStockPriceHistory)
}

type Stock struct {
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"price"`
	IsCrypto bool    `json:"is_crypto"`
	IsStock  bool    `json:"is_stock"`
}

func (s *Service) GetStock(c *gin.Context) {
	idStr := c.Param("id")
	stockID, err := uuid.Parse(idStr)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}
	stock, err := s.queries.GetStockById(c, stockID)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stock)
}

func (s *Service) GetStocks(c *gin.Context) {
	stocks, err := s.queries.GetStocks(c)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stocks)
}

func (s *Service) GetTrendingStocks(c *gin.Context) {
	stocks, err := s.queries.GetTrendingStocks(c)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, stocks)
}

func (s *Service) SearchStocks(c *gin.Context) {
	query := c.Query("query")

	stocks, err := s.queries.SearchStocks(c, query)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stocks)
}

func (s *Service) GetStockPriceHistory(c *gin.Context) {
	idStr := c.Param("id")
	stockID, err := uuid.Parse(idStr)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}

	priceHistory, err := s.queries.GetStockPriceHistory(c, stockID)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, priceHistory)
}
