package portfolio

import (
	"database/sql"
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
	router.GET("/portfolio", s.GetPortfolio)
}

type ReturnPortfolio struct {
	Stock            uuid.UUID `json:"stock"`
	Name             string    `json:"name"`
	Symbol           string    `json:"symbol"`
	Price            float32   `json:"price"`
	IsCrypto         bool      `json:"isCrypto"`
	IsStock          bool      `json:"isStock"`
	Trend            string    `json:"trend"`
	PercentageChange float32   `json:"percentageChange"`
	Quantity         int32     `json:"quantity"`
	Bookmarked       bool      `json:"bookmarked"`
	TotalValue       float32   `json:"totalValue"`
}

func (s *Service) GetPortfolio(c *gin.Context) {
	userId, _ := c.Get("userID")
	portfolio, err := s.queries.GetPortfolio(c, userId.(uuid.UUID))
	if err != nil {
		log.Print(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error."})
		return
	}
	var returnPortfolio []ReturnPortfolio
	for i, stock := range portfolio {
		returnPortfolio[i].Stock = stock.Stock
		returnPortfolio[i].Name = stock.Name
		returnPortfolio[i].Symbol = stock.Symbol
		returnPortfolio[i].Price = stock.Price
		returnPortfolio[i].IsCrypto = stock.IsCrypto
		returnPortfolio[i].IsStock = stock.IsStock
		returnPortfolio[i].Trend = stock.Trend
		returnPortfolio[i].PercentageChange = stock.PercentageChange
		returnPortfolio[i].Quantity = stock.Quantity
		_, err := s.queries.CheckWatchlist(c, db.CheckWatchlistParams{
			Stock: stock.Stock,
			User:  userId.(uuid.UUID),
		})
		if err != nil {
			if err != sql.ErrNoRows {
				returnPortfolio[i].Bookmarked = false
			} else {
				log.Print(err.Error())
				c.JSON(500, gin.H{"error": "Internal server error."})
				return
			}
		} else {
			returnPortfolio[i].Bookmarked = true
		}
		returnPortfolio[i].TotalValue = stock.Price * float32(stock.Quantity)
	}
	c.JSON(200, gin.H{"portfolio": returnPortfolio})
}
