package market

import (
	db "github.com/Boolean-Autocrat/stock-simulator-backend/db/sqlc"
	"github.com/gin-gonic/gin"
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

func (s *Service) sellAsset(c *gin.Context) {
}

func (s *Service) buyAsset(c *gin.Context) {
}

func (s *Service) GetOrderStats(c *gin.Context) {
}
