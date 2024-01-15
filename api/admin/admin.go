package admin

import (
	"fmt"
	"log"
	"net/http"
	"os"

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

type StockForm struct {
	Name        string `form:"name" binding:"required"`
	Symbol      string `form:"symbol" binding:"required"`
	Price       string `form:"price" binding:"required"`
	Quantity    int32  `form:"quantity" binding:"required"`
	StockCrypto string `form:"stock-crypto" binding:"required"`
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.GET("/admin/login", s.adminLogin)
	router.POST("/admin/login", s.adminLoginHandler)
	router.GET("/admin/dashboard", s.adminDashboard)
	router.POST("/admin/stock", s.addStock)
	router.POST("/admin/news", s.addNews)
}

func (s *Service) adminLogin(c *gin.Context) {
	c.HTML(http.StatusOK, "login.tmpl", gin.H{})
}

func (s *Service) adminLoginHandler(c *gin.Context) {
	var form struct {
		Username string `form:"username" binding:"required"`
		Password string `form:"password" binding:"required"`
	}

	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if os.Getenv("ADMIN_USERNAME") != form.Username || os.Getenv("ADMIN_PASSWORD") != form.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}
	c.SetCookie("admin_auth", os.Getenv("ADMIN_SECRET"), 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

func (s *Service) adminDashboard(c *gin.Context) {
	// params := db.GetStocksParams{
	// 	IsCrypto: true,
	// 	IsStock:  true,
	// }
	stocks, err := s.queries.GetStocks(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Print("Stocks are: ", stocks)

	news, err := s.queries.GetArticles(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"stocks": stocks,
		"news":   news,
	})
}

func (s *Service) addStock(c *gin.Context) {
	var form StockForm
	if err := c.ShouldBind(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	params := db.CreateStockParams{
		Name:     form.Name,
		Symbol:   form.Symbol,
		Quantity: form.Quantity,
		Price:    form.Price,
		IsCrypto: form.StockCrypto == "crypto",
		IsStock:  form.StockCrypto == "stock",
	}

	stock, err := s.queries.CreateStock(c, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.HTML(http.StatusOK, "stock_table.tmpl", gin.H{
		"stock": stock,
	})
}

func (s *Service) addNews(c *gin.Context) {

}