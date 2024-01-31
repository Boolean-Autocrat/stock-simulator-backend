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
	router.GET("/market/order/:id/status", s.GetOrderStatus)
}

type Order struct {
	Stock    uuid.UUID `json:"stock"`
	Quantity int       `json:"quantity"`
	Price    float32   `json:"price"`
}

func (s *Service) sellAsset(c *gin.Context) {
	sellID, err := uuid.NewUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order added to orderbook successfully!", "orderID": sellID})
}

func (s *Service) buyAsset(c *gin.Context) {
	buyID, err := uuid.NewUUID()
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order added to orderbook successfully!", "orderID": buyID})
}

func (s *Service) GetOrderStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"message": "Processing!"})
}
