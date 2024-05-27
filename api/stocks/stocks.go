package stocks

import (
	"encoding/json"
	"fmt"
	"log"
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

func (s *Service) RegisterHandlers(router *gin.RouterGroup) {
	router.GET("/stocks", s.GetStocks)
	router.GET("/stocks/watchlist", s.GetWatchlist)
	router.POST("/stocks/watchlist", s.AddToWatchlist)
	router.GET("/stocks/:id", s.GetStock)
	router.GET("/stocks/:id/stream", s.GetStockPriceStream)
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

func (s *Service) GetStockPriceStream(c *gin.Context) {
	idStr := c.Param("id")
	stockID, err := uuid.Parse(idStr)
	if err != nil {
		log.Print(err.Error())
		c.JSON(400, gin.H{"error": "Invalid stock ID"})
		return
	}

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(500, gin.H{"error": "Streaming unsupported!"})
		return
	}

	_, err = s.queries.GetStockById(c, stockID)
	if err != nil {
		log.Print(err.Error())
		c.JSON(400, gin.H{"error": "Invalid stock ID"})
		return
	}

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	flusher.Flush()

	log.Printf("Streaming stock price for stock %s", stockID.String())

	priceChan := make(chan float32)

	go s.StockPriceStream(c, priceChan, stockID)

	for price := range priceChan {
		event, err := json.Marshal(gin.H{"price": price})
		if err != nil {
			log.Print(err.Error())
			c.JSON(500, gin.H{"error": "Internal server error."})
			return
		}
		log.Print(event)
		fmt.Fprintf(c.Writer, "data: %s\n\n", event)
		flusher.Flush()
	}
}

func (s *Service) StockPriceStream(c *gin.Context, priceCh chan<- float32, stockID uuid.UUID) bool {
	ticker := time.NewTicker(2 * time.Second)

outerloop:
	for {
		select {
		case <-c.Done():
			break outerloop
		case <-ticker.C:
			price, err := s.queries.GetStockPrice(c, stockID)
			if err != nil {
				log.Print(err.Error())
				break outerloop
			}
			priceCh <- price
		}
	}

	ticker.Stop()

	close(priceCh)
	return true
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
		c.JSON(400, gin.H{"error": "Stock already in watchlist or invalid stock ID."})
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
