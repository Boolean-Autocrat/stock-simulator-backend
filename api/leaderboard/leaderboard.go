package leaderboard

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

type LeaderboardItem struct {
	FullName string  `json:"fullName"`
	Balance  float32 `json:"balance"`
	Position int     `json:"position"`
	Picture  string  `json:"picture"`
	IsYou    bool    `json:"isYou"`
}

func (s *Service) RegisterHandlers(router *gin.RouterGroup) {
	router.GET("/leaderboard", s.getLeaderboard)
}

func (s *Service) getLeaderboard(c *gin.Context) {
	userId, _ := c.Get("userID")
	leaderboard, err := s.queries.GetLeaderboard(c)
	if err != nil {
		log.Print(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error."})
		return
	}
	counter := 0
	userFlag := false
	var leaderboardItems []LeaderboardItem
	for i := range leaderboard {
		counter++
		if leaderboard[i].ID == userId.(uuid.UUID) {
			userFlag = true
		}
		item := LeaderboardItem{
			FullName: leaderboard[i].FullName,
			Balance:  leaderboard[i].Balance,
			Position: counter,
			Picture:  leaderboard[i].Picture,
			IsYou:    leaderboard[i].ID == userId.(uuid.UUID),
		}
		leaderboardItems = append(leaderboardItems, item)
	}
	if !userFlag {
		userPos, err := s.queries.GetUserPosition(c, userId.(uuid.UUID))
		if err != nil {
			log.Print(err.Error())
			c.JSON(500, gin.H{"error": "Internal server error."})
			return
		}
		userItem := LeaderboardItem{
			FullName: userPos.FullName,
			Balance:  userPos.Balance,
			Position: int(userPos.Position),
			Picture:  userPos.Picture,
			IsYou:    true,
		}
		leaderboardItems = append(leaderboardItems, userItem)
	}
	c.JSON(http.StatusOK, leaderboardItems)
}
