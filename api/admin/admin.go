package admin

import (
	"log"
	"net/http"
	"os"

	db "github.com/Boolean-Autocrat/stock-simulator-backend/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println(err.Error())
		log.Fatal("Error loading .env file")
	}
}

type Service struct {
	queries *db.Queries
}

type StockForm struct {
	Name        string  `form:"name" binding:"required"`
	Symbol      string  `form:"symbol" binding:"required"`
	Price       float32 `form:"price" binding:"required"`
	Quantity    int32   `form:"quantity" binding:"required"`
	StockCrypto string  `form:"stock-crypto" binding:"required"`
}

type NewsForm struct {
	Title   string `form:"title" binding:"required"`
	Author  string `form:"author" binding:"required"`
	Content string `form:"content" binding:"required"`
	Tag     string `form:"tag" binding:"required"`
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) RegisterHandlers(router *gin.RouterGroup) {
	router.GET("/login", s.adminLogin)
	router.POST("/login", s.adminLoginHandler)
	router.GET("/logout", s.adminLogout)
	router.GET("/dashboard", s.adminDashboard)
	router.POST("/stock", s.addStock)
	router.POST("/news", s.addNews)
	router.GET("/news/:id/delete", s.deleteNews)
	router.GET("/news/:id/edit", s.editNews)
	router.POST("/news/:id/edit", s.editNewsHandler)
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
		log.Println(err.Error())
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	if os.Getenv("ADMIN_USERNAME") != form.Username || os.Getenv("ADMIN_PASSWORD") != form.Password {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("admin_auth", os.Getenv("ADMIN_SECRET"), 3600*24*15, "", "", false, true)
	c.Redirect(http.StatusFound, "/admin/dashboard")
}

func (s *Service) adminLogout(c *gin.Context) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie("admin_auth", "", -1, "", "", false, true)
	c.Redirect(http.StatusFound, "/admin/login")
}

func (s *Service) adminDashboard(c *gin.Context) {
	stocks, err := s.queries.GetStocks(c)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}

	news, err := s.queries.GetArticles(c)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error"})
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
		log.Println(err.Error())
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	params := db.CreateStockParams{
		Name:        form.Name,
		Symbol:      form.Symbol,
		IpoQuantity: form.Quantity,
		Price:       form.Price,
		IsCrypto:    form.StockCrypto == "crypto",
		IsStock:     form.StockCrypto == "stock",
	}

	stock, err := s.queries.CreateStock(c, params)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	s.queries.CreatePriceHistory(c, db.CreatePriceHistoryParams{
		Stock: stock.ID,
		Price: stock.Price,
	})
	c.HTML(http.StatusOK, "stock_table.tmpl", gin.H{
		"stock": stock,
	})
}

func (s *Service) addNews(c *gin.Context) {
	var form NewsForm
	if err := c.ShouldBind(&form); err != nil {
		log.Println(err.Error())
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	params := db.AddArticleParams{
		Title:   form.Title,
		Author:  form.Author,
		Content: form.Content,
		Tag:     form.Tag,
	}

	article, err := s.queries.AddArticle(c, params)
	if err != nil {
		log.Println(err.Error())
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.HTML(http.StatusOK, "article_table.tmpl", gin.H{
		"article": article,
	})
}

func (s *Service) deleteNews(c *gin.Context) {
	articleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}
	delArticleErr := s.queries.DeleteArticle(c, articleID)
	if delArticleErr != nil {
		log.Println(delArticleErr.Error())
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.Redirect(302, "/admin/dashboard#news-table-body")
}

func (s *Service) editNews(c *gin.Context) {
	articleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}
	article, err := s.queries.GetArticle(c, articleID)
	if err != nil {
		log.Println(err.Error())
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}
	c.HTML(http.StatusOK, "article_edit.tmpl", gin.H{
		"article": article,
	})
}

func (s *Service) editNewsHandler(c *gin.Context) {
	articleID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		log.Println(err.Error())
		c.JSON(400, gin.H{"message": "Invalid request"})
		return
	}
	var form NewsForm
	if err := c.ShouldBind(&form); err != nil {
		log.Println(err.Error())
		c.JSON(400, gin.H{"error": "Invalid request"})
		return
	}

	params := db.UpdateArticleParams{
		Title:   form.Title,
		Author:  form.Author,
		Content: form.Content,
		Tag:     form.Tag,
		ID:      articleID,
	}
	articleErr := s.queries.UpdateArticle(c, params)
	if articleErr != nil {
		log.Println(articleErr.Error())
		c.JSON(500, gin.H{"error": "Internal server error"})
		return
	}
	c.Redirect(302, "/admin/dashboard#news-table-body")
}
