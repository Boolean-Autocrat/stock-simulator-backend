package middleware

import (
	"net/http"
	"time"

	db "github.com/Boolean-Autocrat/stock-simulator-backend/db/sqlc"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) TokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.FullPath() == "/auth/google/login" || c.FullPath() == "/auth/google/callback" {
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
		if tokenData.ExpiresAt.Before(time.Now()) {
			refreshToken, _ := s.queries.GetRefreshToken(c, tokenData.UserID)
			if refreshToken.ExpiresAt.Before(time.Now()) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token expired"})
				c.Abort()
				return
			}
			refToken := &oauth2.Token{
				RefreshToken: refreshToken.Token,
			}
			newToken, err := refreshAccessToken(refToken)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to refresh access token"})
				c.Abort()
				return
			}
			params := db.UpdateAccessTokenParams{
				Token:     newToken.AccessToken,
				ExpiresAt: newToken.Expiry,
				UserID:    tokenData.UserID,
			}
			_, err = s.queries.UpdateAccessToken(c, params)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update access token"})
				c.Abort()
				return
			}
			c.Set("userID", tokenData.UserID)
			c.Header("Authorization", newToken.AccessToken)
			c.JSON(http.StatusOK, gin.H{"freshToken": "true"})
		} else {
			c.Set("userID", tokenData.UserID)
			c.Next()
		}
	}
}
