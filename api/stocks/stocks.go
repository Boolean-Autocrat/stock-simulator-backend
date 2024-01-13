package stocks

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

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
	router.POST("/stocks", s.CreateStock)
	router.GET("/stocks/search/:query", s.SearchStocks)
	router.GET("/stocks/:id/price_history", s.GetStockPriceHistory)
	router.GET("/stocks/:id/price_history/:start/:end", s.GetStockPriceHistoryByDate)
	router.POST("/stocks/:id/price_history", s.CreatePriceHistory)
}

type Stock struct {
	Name     string  `json:"name"`
	Symbol   string  `json:"symbol"`
	Price    float64 `json:"price"`
	IsCrypto bool    `json:"is_crypto"`
	IsStock  bool    `json:"is_stock"`
}

func (s *Service) CreateStock(c *gin.Context) {
	var params db.CreateStockParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stock, err := s.queries.CreateStock(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, stock)
}

func (s *Service) GetStock(c *gin.Context) {
	var params db.GetStockParams
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stock, err := s.queries.GetStock(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stock)
}

func (s *Service) GetStocks(c *gin.Context) {
	isCrypto, _ := strconv.ParseBool(c.Query("is_crypto"))
	isStock, _ := strconv.ParseBool(c.Query("is_stock"))
	params := db.GetStocksParams{
		sql.NullBool{Bool: isCrypto, Valid: true},
		sql.NullBool{Bool: isStock, Valid: true},
	}
	stocks, err := s.queries.GetStocks(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stocks)
}

func (s *Service) SearchStocks(c *gin.Context) {
	query := c.Param("query")

	stocks, err := s.queries.SearchStocks(c, query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stocks)
}

func (s *Service) GetStockPriceHistory(c *gin.Context) {
	idStr := c.Param("id")
	stockID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}

	priceHistory, err := s.queries.GetStockPriceHistory(c, stockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, priceHistory)
}

func (s *Service) GetStockPriceHistoryByDate(c *gin.Context) {
	idStr := c.Param("id")
	stockID, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}

	startDate, err := time.Parse(time.RFC3339, c.Param("start"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
		return
	}

	endDate, err := time.Parse(time.RFC3339, c.Param("end"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
		return
	}

	params := db.GetStockPriceHistoryByDateParams{
		stockID,
		startDate,
		endDate,
	}

	priceHistory, err := s.queries.GetStockPriceHistoryByDate(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, priceHistory)
}

func (s *Service) CreatePriceHistory(c *gin.Context) {
	stockID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid stock ID"})
		return
	}

	var input struct {
		Price float64 `json:"price"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	priceHistory, err := s.queries.CreatePriceHistory(c, stockID, input.Price)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, priceHistory)
}
