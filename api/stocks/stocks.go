package stocks

import (
	"log"

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
	router.GET("/stocks", s.GetStocks)
	router.GET("/stocks/watchlist", s.GetWatchlist)
	router.POST("/stocks/watchlist", s.AddToWatchlist)
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
		c.JSON(400, gin.H{"error": "Invalid stock ID"})
		return
	}
	stock, err := s.queries.GetStockById(c, stockID)
	if err != nil {
		log.Print(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error."})
		return
	}

	c.JSON(200, stock)
}

func (s *Service) GetStocks(c *gin.Context) {
	stocks, err := s.queries.GetStocks(c)
	if err != nil {
		log.Print(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error."})
		return
	}
	c.JSON(200, stocks)
}

func (s *Service) GetWatchlist(c *gin.Context) {
	userID, _ := c.Get("userID")
	watchlist, err := s.queries.GetWatchlist(c, userID.(uuid.UUID))
	if err != nil {
		log.Print(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error."})
		return
	}
	c.JSON(200, gin.H{"watchlist": watchlist})
}

type AddToWatchlistRequest struct {
	StockID uuid.UUID `json:"stock"`
}

func (s *Service) AddToWatchlist(c *gin.Context) {
	userID, _ := c.Get("userID")
	var stock AddToWatchlistRequest
	if err := c.ShouldBindJSON(&stock); err != nil {
		log.Print(err.Error())
		c.JSON(400, gin.H{"error": "Invalid request."})
		return
	}
	err := s.queries.AddToWatchlist(c, db.AddToWatchlistParams{
		User:  userID.(uuid.UUID),
		Stock: stock.StockID,
	})
	if err != nil {
		log.Print(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error."})
		return
	}
	c.JSON(200, gin.H{"message": "Stock added to watchlist."})
}

func (s *Service) GetTrendingStocks(c *gin.Context) {
	stocks, err := s.queries.GetTrendingStocks(c)
	if err != nil {
		log.Print(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error."})
		return
	}
	c.JSON(200, stocks)
}

func (s *Service) SearchStocks(c *gin.Context) {
	query := c.Query("query")

	stocks, err := s.queries.SearchStocks(c, query)
	if err != nil {
		log.Print(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error."})
		return
	}

	c.JSON(200, stocks)
}

func (s *Service) GetStockPriceHistory(c *gin.Context) {
	idStr := c.Param("id")
	stockID, err := uuid.Parse(idStr)
	if err != nil {
		log.Print(err.Error())
		c.JSON(400, gin.H{"error": "Invalid stock ID"})
		return
	}

	priceHistory, err := s.queries.GetStockPriceHistory(c, stockID)
	if err != nil {
		log.Print(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error."})
		return
	}

	c.JSON(200, priceHistory)
}
