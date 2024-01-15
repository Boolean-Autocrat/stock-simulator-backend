package portfolio

import (
	"net/http"

	"github.com/Boolean-Autocrat/stock-simulator-backend/api/utils"
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
	router.GET("/portfolio", s.GetPortfolio)
	router.POST("/portfolio/sell", s.SellStock)
	router.POST("/portfolio/buy", s.BuyStock)
}

func (s *Service) GetPortfolio(c *gin.Context) {
	userId, _ := utils.GetUserfromContext(c)
	var params db.GetPortfolioParams
	params.UserID = userId
	if err := c.ShouldBindJSON(&params); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	portfolio, err := s.queries.GetPortfolio(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, portfolio)
}

func (s *Service) SellStock(c *gin.Context) {
	return
}

func (s *Service) BuyStock(c *gin.Context) {
	return
}