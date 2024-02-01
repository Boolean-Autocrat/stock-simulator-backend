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
	// p, err := kafka.NewProducer(&kafka.ConfigMap{
	// 	"bootstrap.servers": "host1:9092",
	// 	"client.id":         "stock-simulator-exchange",
	// 	"acks":              "all",
	// })
	// if err != nil {
	// 	fmt.Printf("Failed to create producer: %s\n", err)
	// 	os.Exit(1)
	// }
	// consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
	// 	"bootstrap.servers": "host1:9092,host2:9092",
	// 	"group.id":          "foo",
	// 	"auto.offset.reset": "smallest"})
	// if err != nil {
	// 	fmt.Printf("Failed to create consumer: %s\n", err)
	// 	os.Exit(1)
	// }
	// book := engine.OrderBook{
	// 	BuyOrders:  make([]engine.Order, 0, 1000),
	// 	SellOrders: make([]engine.Order, 0, 1000),
	// }

	// done := make(chan bool)
	// tradesTopic := "trades"
	// go func() {
	// 	for {
	// 		msg, _ := consumer.ReadMessage(-1)
	// 		var order engine.Order
	// 		order.FromJSON(msg.Value)
	// 		trades := book.Process(order)
	// 		for _, trade := range trades {
	// 			rawTrade := trade.ToJSON()
	// 			p.Produce(&kafka.Message{
	// 				TopicPartition: kafka.TopicPartition{
	// 					Topic: &tradesTopic,
	// 					// Partition: kafka.PartitionAny, // TODO: Add partitioning by stock symbol
	// 				},
	// 				Value: rawTrade,
	// 			}, nil)
	// 		}
	// 		consumer.CommitMessage(msg)
	// 	}
	// 	done <- true
	// }()
	router.Run()
}
