package coursecodes

import (
	"os"

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
	router.POST("/courses", s.courses)
}

func (s *Service) courses(c *gin.Context) {
	jsonFile, err := os.Open("courseCodes.json")
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	defer jsonFile.Close()
	c.JSON(200, jsonFile)
}
