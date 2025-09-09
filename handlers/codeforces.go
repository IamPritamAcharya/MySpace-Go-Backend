package handlers

import (
	"myspace-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetCodeforcesRating(c *gin.Context) {
	handle := c.Param("handle")
	ratings, err := services.GetCodeforcesRating(handle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, ratings)
}

func GetCodeforcesContests(c *gin.Context) {
	contests, err := services.GetCodeforcesContests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, contests)
}
