package main

import (
	"context"
	"fmt"
	"log"
	"os"

	_ "time/tzdata"

	"github.com/Boolean-Autocrat/stock-simulator-backend/api/admin"
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

	adminClient, _ := kafka.NewAdminClient(&kafka.ConfigMap{
		"bootstrap.servers": os.Getenv("KAFKA_BOOTSTRAP_SERVERS"),
	})
	var topicsExist bool = false
	topics, _ := adminClient.GetMetadata(nil, true, 5000)
	for _, topic := range topics.Topics {
		if topic.Topic == "trades" || topic.Topic == "orders" {
			topicsExist = true
		}
	}
	if !topicsExist {
		var mapTopics = []kafka.TopicSpecification{}
		mapTopics = append(mapTopics, kafka.TopicSpecification{
			Topic:             "trades",
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
		mapTopics = append(mapTopics, kafka.TopicSpecification{
			Topic:             "orders",
			NumPartitions:     1,
			ReplicationFactor: 1,
		})
		results, err := adminClient.CreateTopics(context.Background(), mapTopics)
		if err != nil {
			log.Println("Error creating topics:", err)
		} else {
			for _, result := range results {
				fmt.Printf("Topic %s created: %v\n", result.Topic, result.Error)
			}
		}
	}
	adminClient.Close()
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
		log.Println("Starting trade processor")
		for {
			msg, err := consumer.ReadMessage(-1)
			if err != nil {
				log.Printf("Error reading message: %v\n", err)
				continue
			}

			var order engine.Order
			order.FromJSON(msg.Value)
			book, exists := orderBooks[order.Stock.String()]
			if !exists {
				book = &engine.OrderBook{
					BuyOrders:  make([]engine.Order, 0, 1000),
					SellOrders: make([]engine.Order, 0, 1000),
				}
				orderBooks[order.Stock.String()] = book
			}

			trades := book.Process(order)
			for _, trade := range trades {
				engine.RunTradeQueries(trade, queries)
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
	marketService := market.NewService(queries, producer)
	// courseService := coursecodes.NewService(queries)
	ipoService := ipo.NewService(queries)

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/*")
	router.Static("/assets", "./assets")
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Health OK!",
		})
	})
	adminGroup := router.Group("/admin")
	adminGroup.Use(middleware.NewService(queries).AdminMiddleware())
	adminService.RegisterHandlers(adminGroup)

	usersGroup := router.Group("")
	usersGroup.Use(middleware.NewService(queries).TokenMiddleware())
	authService.RegisterHandlers(usersGroup)
	stockService.RegisterHandlers(usersGroup)
	newsService.RegisterHandlers(usersGroup)
	portfolioService.RegisterHandlers(usersGroup)
	leaderboardService.RegisterHandlers(usersGroup)
	marketService.RegisterHandlers(usersGroup)
	// courseService.RegisterHandlers(usersGroup)
	ipoService.RegisterHandlers(usersGroup)

	router.Run()
}
