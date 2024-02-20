package coursecodes

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	db "github.com/Boolean-Autocrat/stock-simulator-backend/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Service struct {
	queries *db.Queries
}

func NewService(queries *db.Queries) *Service {
	return &Service{queries: queries}
}

func (s *Service) RegisterHandlers(router *gin.RouterGroup) {
	router.GET("/courses", s.courses)
}

func (s *Service) courses(c *gin.Context) {
	type Course struct {
		Department string `json:"department"`
		Year       string `json:"year"`
		CourseCode string `json:"courseCode"`
		CourseName string `json:"courseName"`
	}

	type CourseData struct {
		Courses []Course `json:"courses"`
	}
	filename := "courseCodes.json"
	data, err := os.Open(filename)
	if err != nil {
		log.Print(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer data.Close()
	byteValue, _ := ioutil.ReadAll(data)
	var courseData CourseData
	json.Unmarshal(byteValue, &courseData)
	c.JSON(200, courseData)
}
