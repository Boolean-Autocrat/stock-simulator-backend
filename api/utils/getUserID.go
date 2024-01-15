package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func GetUserfromContext(c *gin.Context) (uuid.UUID, error) {
	idStr, _ := c.Get("userID")
	id, err := uuid.Parse(idStr.(string))
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}
