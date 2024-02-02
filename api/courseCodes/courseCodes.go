package coursecodes

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	db "github.com/Boolean-Autocrat/stock-simulator-backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) RegisterHandlers(router *gin.Engine) {
	router.GET("/courses", s.courses)
}

type Course struct {
	Department string `json:"department"`
	Year       string `json:"year"`
	CourseCode string `json:"courseCode"`
	CourseName string `json:"courseName"`
}

type CourseData struct {
	Courses []Course `json:"courses"`
}

func (s *Service) courses(c *gin.Context) {
	filename := "courseCodes.json"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read the JSON file"})
		return
	}

	var courseData CourseData
	err = json.Unmarshal(data, &courseData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse JSON data"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"courses": courseData.Courses})
}
