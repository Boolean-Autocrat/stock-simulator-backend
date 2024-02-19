package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Boolean-Autocrat/stock-simulator-backend/api/admin"
	coursecodes "github.com/Boolean-Autocrat/stock-simulator-backend/api/courseCodes"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/ipo"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/leaderboard"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/market"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/middleware"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/news"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/portfolio"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/stocks"
	"github.com/Boolean-Autocrat/stock-simulator-backend/api/userAuth"
	db "github.com/Boolean-Autocrat/stock-simulator-backend/db/sqlc"
	"github.com/Boolean-Autocrat/stock-simulator-backend/engine"
	"github.com/confluentinc/confluent-kafka-go/kafka"
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

	_, timezoneErr := postgres.DB.Exec("SET TIME ZONE 'Asia/Kolkata'")
	if timezoneErr != nil {
		log.Fatal(timezoneErr.Error())
	}

	queries := db.New(postgres.DB)

	consumer := engine.CreateConsumer()
	producer := engine.CreateProducer()

	orderBooks := make(map[string]*engine.OrderBook)
	done := make(chan bool)
	tradesTopic := "trades"

	go func() {
		fmt.Println("Starting trade processor")
		for {
			msg, err := consumer.ReadMessage(-1)
			if err != nil {
				log.Printf("Error reading message: %v\n", err)
				continue
			}

			var order engine.Order
			order.FromJSON(msg.Value)
			book, exists := orderBooks[order.ID.String()]
			if !exists {
				book = &engine.OrderBook{
					BuyOrders:  make([]engine.Order, 0, 1000),
					SellOrders: make([]engine.Order, 0, 1000),
				}
				orderBooks[order.ID.String()] = book
			}

			trades := book.Process(order)
			for _, trade := range trades {
				rawTrade := trade.ToJSON()
				producer.Produce(&kafka.Message{
					TopicPartition: kafka.TopicPartition{
						Topic:     &tradesTopic,
						Partition: kafka.PartitionAny,
					},
					Value: rawTrade,
				}, nil)
			}

			consumer.CommitMessage(msg)
		}
		done <- true
	}()

	adminService := admin.NewService(queries)
	authService := userAuth.NewService(queries)
	stockService := stocks.NewService(queries)
	newsService := news.NewService(queries)
	portfolioService := portfolio.NewService(queries)
	leaderboardService := leaderboard.NewService(queries)
	marketService := market.NewService(queries)
	courseService := coursecodes.NewService(queries)
	ipoService := ipo.NewService(queries)

	// gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Health OK!",
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
	courseService.RegisterHandlers(router)
	ipoService.RegisterHandlers(router)
	router.Run()
}
