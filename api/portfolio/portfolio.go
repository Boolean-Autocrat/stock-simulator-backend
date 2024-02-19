package portfolio

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
	router.GET("/portfolio", s.GetPortfolio)
}

func (s *Service) GetPortfolio(c *gin.Context) {
	userId, _ := c.Get("userID")
	portfolio, err := s.queries.GetPortfolio(c, userId.(uuid.UUID))
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if portfolio == nil {
		c.JSON(200, gin.H{})
		return
	}
	c.JSON(http.StatusOK, portfolio)
}
