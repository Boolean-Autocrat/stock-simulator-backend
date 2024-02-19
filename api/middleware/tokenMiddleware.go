package middleware

import (
	"log"
	"net/http"
	"os"
	"strings"

	db "github.com/Boolean-Autocrat/stock-simulator-backend/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/admin/login" || c.FullPath() == "/courses" {
			c.Next()
			return
		} else if strings.Contains(c.FullPath(), "/admin") {
			cookie, err := c.Cookie("admin_auth")
			if err != nil {
				c.Redirect(http.StatusFound, "/admin/login")
				return
			}
			if cookie != os.Getenv("ADMIN_SECRET") {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				c.Abort()
				return
			}
			c.Next()
			return
		} else if c.FullPath() == "/auth/google/login" || c.FullPath() == "/auth/google/callback" {
			c.Next()
			return
		}
		accessToken := c.GetHeader("Authorization")
		if accessToken == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token provided"})
			c.Abort()
			return
		}
		tokenData, err := s.queries.GetTokenData(c, accessToken)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid access token"})
			c.Abort()
			return
		}
		c.Set("userID", tokenData)
		c.Next()
	}
}
