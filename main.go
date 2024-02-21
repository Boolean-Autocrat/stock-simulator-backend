package main

import (
	"log"
	"os"

	_ "time/tzdata"

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
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	amqp "github.com/rabbitmq/amqp091-go"
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

	orderBooks := make(map[string]*engine.OrderBook)

	amqpServerURL := os.Getenv("AMQP_SERVER_URL")
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}
	defer connectRabbitMQ.Close()
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}
	defer channelRabbitMQ.Close()
	_, err = channelRabbitMQ.QueueDeclare(
		"orders", // queue name
		true,     // durable
		false,    // auto delete
		false,    // exclusive
		false,    // no wait
		nil,      // arguments
	)
	if err != nil {
		panic(err)
	}
	messages, err := channelRabbitMQ.Consume(
		"orders", // queue name
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no local
		false,    // no wait
		nil,      // arguments
	)
	if err != nil {
		log.Println(err)
	}

	forever := make(chan bool)

	go func() {
		defer close(forever)
		for message := range messages {
			var order engine.Order
			order.FromJSON(message.Body)
			log.Printf("Received a message: %v", order)
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
			}
		}
	}()

	engine.AddPendingOrders(queries, channelRabbitMQ)

	adminService := admin.NewService(queries)
	authService := userAuth.NewService(queries)
	stockService := stocks.NewService(queries)
	newsService := news.NewService(queries)
	portfolioService := portfolio.NewService(queries)
	leaderboardService := leaderboard.NewService(queries)
	marketService := market.NewService(queries, channelRabbitMQ)
	courseService := coursecodes.NewService(queries)
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
	courseService.RegisterHandlers(usersGroup)
	ipoService.RegisterHandlers(usersGroup)

	router.Run()
}
