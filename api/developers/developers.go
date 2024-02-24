package developers

import (
	"log"

	db "github.com/Boolean-Autocrat/stock-simulator-backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) RegisterHandlers(router *gin.RouterGroup) {
	router.GET("/developers", s.developers)
}

func (s *Service) developers(c *gin.Context) {
	developers, err := s.queries.GetDevelopers(c)
	if err != nil {
		log.Print(err.Error())
		c.JSON(500, gin.H{"error": "internal server error"})
		return
	}
	c.JSON(200, gin.H{"developers": developers})
}
