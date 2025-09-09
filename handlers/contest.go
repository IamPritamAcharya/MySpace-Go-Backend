package handlers

import (
	"myspace-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetAllContests(c *gin.Context) {
	contests, err := services.GetAllContests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, contests)
}
