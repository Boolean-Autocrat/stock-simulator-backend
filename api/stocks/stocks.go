package stocks

import (
	"net/http"
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
	router.GET("/stocks/search", s.SearchStocks)
	router.GET("/stocks/:id/price_history", s.GetStockPriceHistory)
	// router.POST("/stocks/:id/price_history", s.CreatePriceHistory)
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
	// isCrypto, _ := strconv.ParseBool(c.Query("is_crypto"))
	// isStock, _ := strconv.ParseBool(c.Query("is_stock"))
	// params := db.GetStocksParams{
	// 	IsCrypto: isCrypto,
	// 	IsStock:  isStock,
	// }
	stocks, err := s.queries.GetStocks(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stocks)
}

func (s *Service) SearchStocks(c *gin.Context) {
	query := c.Query("query")

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

	var startDate, endDate time.Time

	if c.Query("start_date") != "" && c.Query("end_date") != "" {
		startDate, err = time.Parse(time.RFC3339, c.Param("start"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid start date"})
			return
		}

		endDate, err = time.Parse(time.RFC3339, c.Param("end"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid end date"})
			return
		}
	} else {
		startDate = time.Now().AddDate(-3, 0, 0)
		endDate = time.Now().AddDate(0, 0, 1)
	}

	params := db.GetStockPriceHistoryByDateParams{
		StockID:   stockID,
		PriceAt:   startDate,
		PriceAt_2: endDate,
	}

	priceHistory, err := s.queries.GetStockPriceHistoryByDate(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, priceHistory)
}

func (s *Service) CreatePriceHistory(c *gin.Context) {}
