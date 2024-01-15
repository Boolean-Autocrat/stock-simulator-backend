package main

import (
	"log"
	"os"

	"github.com/Boolean-Autocrat/stock-simulator-backend/api/middleware"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/stocks"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/userAuth"
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

func main() {
	postgres, err := db.NewPostgres(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("DB_HOST"))
	if err != nil {
		log.Fatal(err.Error())
	}

	queries := db.New(postgres.DB)
	authService := userAuth.NewService(queries)
	stockService := stocks.NewService(queries)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
	router.Use(middleware.TokenMiddleware())
	authService.RegisterHandlers(router)
	stockService.RegisterHandlers(router)

	router.Run()
}
