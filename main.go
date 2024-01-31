package main

import (
	"log"
	"os"

	"github.com/Boolean-Autocrat/stock-simulator-backend/api/admin"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/leaderboard"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/market"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/middleware"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/news"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/portfolio"
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
	// brokerAddrs := []string{"localhost:9092"}
	// config := sarama.NewConfig()
	// saramaAdmin, err := sarama.NewClusterAdmin(brokerAddrs, config)
	// if err != nil {
	// 	log.Fatal("Error while creating cluster admin: ", err.Error())
	// }
	// defer func() { _ = saramaAdmin.Close() }()
	// err = saramaAdmin.("topic.test.1", &sarama.TopicDetail{
	// 	NumPartitions:     1,
	// 	ReplicationFactor: 1,
	// }, false)
	// if err != nil {
	// 	log.Fatal("Error while creating topic: ", err.Error())
	// }
	postgres, err := db.NewPostgres(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"), os.Getenv("DB_HOST"))
	if err != nil {
		log.Fatal(err.Error())
	}

	queries := db.New(postgres.DB)
	adminService := admin.NewService(queries)
	authService := userAuth.NewService(queries)
	stockService := stocks.NewService(queries)
	newsService := news.NewService(queries)
	portfolioService := portfolio.NewService(queries)
	leaderboardService := leaderboard.NewService(queries)
	marketService := market.NewService(queries)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "OK",
		})
	})
	router.Use(middleware.NewService(queries).TokenMiddleware())
	adminService.RegisterHandlers(router)
	authService.RegisterHandlers(router)
	stockService.RegisterHandlers(router)
	newsService.RegisterHandlers(router)
	portfolioService.RegisterHandlers(router)
	leaderboardService.RegisterHandlers(router)
	marketService.RegisterHandlers(router)

	router.Run()
}
