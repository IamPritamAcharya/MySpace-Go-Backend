package handlers

import (
	"myspace-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetTechNews(c *gin.Context) {
	news, err := services.GetTechNews()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, news)
}