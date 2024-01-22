package news

import (
	"fmt"
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

func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.GET("/news", s.getAllNews)
	router.GET("/news/:id", s.getNews)
	router.POST("/news/:id/:type", s.addNewsSentiment)
}

func (s *Service) getAllNews(c *gin.Context) {
	type NewsItem struct {
		ID        uuid.UUID `json:"id"`
		Title     string    `json:"title"`
		Author    string    `json:"author"`
		Content   string    `json:"content"`
		Tag       string    `json:"tag"`
		CreatedAt time.Time `json:"createdAt"`
		TimeAgo   string    `json:"timeAgo"`
	}
	news, err := s.queries.GetArticles(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var newsItems []NewsItem
	for i := range news {
		item := NewsItem{
			ID:        news[i].ID,
			Title:     news[i].Title,
			Author:    news[i].Author,
			Content:   news[i].Content,
			Tag:       news[i].Tag,
			CreatedAt: news[i].CreatedAt,
			TimeAgo:   calculateTimeAgo(news[i].CreatedAt),
		}
		newsItems = append(newsItems, item)
	}
	c.JSON(http.StatusOK, newsItems)
}

func calculateTimeAgo(timestamp time.Time) string {
	now := time.Now().UTC()
	diff := now.Sub(timestamp)

	switch {
	case diff.Hours() >= 24*7:
		return fmt.Sprintf("%.0f weeks ago", diff.Hours()/(24*7))
	case diff.Hours() >= 24:
		return fmt.Sprintf("%.0f days ago", diff.Hours()/24)
	case diff.Hours() >= 1:
		return fmt.Sprintf("%.0f hours ago", diff.Hours())
	case diff.Minutes() >= 1:
		return fmt.Sprintf("%.0f minutes ago", diff.Minutes())
	default:
		return fmt.Sprintf("%.0f seconds ago", diff.Seconds())
	}
}

func (s *Service) getNews(c *gin.Context) {
	userId, _ := c.Get("userID")
	newsIdStr := c.Param("id")
	newsID, err := uuid.Parse(newsIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	newsItem, err := s.queries.GetArticle(c, newsID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	articleSentiment, _ := s.queries.GetArticleSentiment(c, newsID)
	userSentiment, _ := s.queries.GetUserSentiment(c, db.GetUserSentimentParams{
		UserID:    userId.(uuid.UUID),
		ArticleID: newsID,
	})
	var userSentimentStr string
	if userSentiment.Like {
		userSentimentStr = "like"
	} else if userSentiment.Dislike {
		userSentimentStr = "dislike"
	} else {
		userSentimentStr = "none"
	}
	c.JSON(http.StatusOK, gin.H{
		"article": newsItem,
		"sentiment": gin.H{
			"article": articleSentiment,
			"user":    userSentimentStr,
		},
	})
}

func (s *Service) addNewsSentiment(c *gin.Context) {
	userId, _ := c.Get("userID")
	newsIdStr := c.Param("id")
	newsID, err := uuid.Parse(newsIdStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	typeStr := c.Param("type")
	if typeStr != "like" && typeStr != "dislike" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid type parameter"})
		return
	}
	err = s.queries.AddArticleSentiment(c, db.AddArticleSentimentParams{
		UserID:    userId.(uuid.UUID),
		ArticleID: newsID,
		Like:      typeStr == "like",
		Dislike:   typeStr == "dislike",
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	articleSentiment, _ := s.queries.GetArticleSentiment(c, newsID)
	c.JSON(http.StatusOK, articleSentiment)
}
