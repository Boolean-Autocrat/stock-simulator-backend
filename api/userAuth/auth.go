package userAuth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	db "github.com/Boolean-Autocrat/stock-simulator-backend/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  os.Getenv("REDIRECT_URL"),
		ClientID:     os.Getenv("CLIENT_ID"),
		ClientSecret: os.Getenv("CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
}

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

type UserInfo struct {
	Sub           string `json:"sub"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Profile       string `json:"profile"`
	Picture       string `json:"picture"`
	Email         string `json:"email"`
	EmailVerified bool   `json:"email_verified"`
}

func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.GET("/auth/google/login", s.GoogleLogin)
	router.GET("/auth/google/callback", s.GoogleCallback)
	router.GET("/auth/userinfo", s.GetUserInfo)
}

func (s *Service) GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("state")
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func (s *Service) GoogleCallback(c *gin.Context) {
	code := c.Query("code")

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to exchange token"})
		return
	}
	idToken, ok := extractIDToken(token)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to extract ID token"})
		return
	}

	valid, err := verifyIDToken(idToken)
	if err != nil || !valid {
		fmt.Print(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid ID token"})
		return
	}

	userInfo, err := getGoogleUserInfo(token)
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user info"})
		return
	}
	user, err := s.queries.CreateUser(c, db.CreateUserParams{
		FullName: userInfo.Name,
		Email:    userInfo.Email,
		Picture:  userInfo.Picture,
	})
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
		return
	}

	_, err = s.queries.CreateAccessToken(c, db.CreateAccessTokenParams{
		UserID:    user.ID,
		Token:     token.AccessToken,
		ExpiresAt: token.Expiry.AddDate(0, 0, 365),
	})
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create access token"})
		return
	}

	_, err = s.queries.CreateRefreshToken(c, db.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     token.RefreshToken,
		ExpiresAt: token.Expiry,
	})
	if err != nil {
		fmt.Print(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create refresh token"})
		return
	}

	returnParams := gin.H{
		"accessToken": token.AccessToken,
		"user": gin.H{
			"id":       user.ID,
			"fullName": user.FullName,
			"email":    user.Email,
			"picture":  user.Picture,
		},
	}

	c.JSON(http.StatusOK, returnParams)
}

func (s *Service) GetUserInfo(c *gin.Context) {
	userID, _ := c.Get("userID")
	user, err := s.queries.GetUser(c, userID.(uuid.UUID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func getGoogleUserInfo(token *oauth2.Token) (*UserInfo, error) {
	client := googleOauthConfig.Client(context.Background(), token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	defer response.Body.Close()

	var userInfo UserInfo
	err = json.NewDecoder(response.Body).Decode(&userInfo)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}

	return &userInfo, nil
}

func extractIDToken(response *oauth2.Token) (string, bool) {
	idToken, ok := response.Extra("id_token").(string)
	return idToken, ok
}
