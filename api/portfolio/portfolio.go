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
	var params db.GetPortfolioParams
	params.UserID = userId.(uuid.UUID)
	if err := c.ShouldBindJSON(&params); err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	portfolio, err := s.queries.GetPortfolio(c, params)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, portfolio)
}
