package handlers

import (
	"myspace-backend/services"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetLeetCodeRecentSubmissions(c *gin.Context) {
	username := c.Param("username")
	data, err := services.GetLeetCodeRecentSubmissions(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func GetLeetCodeProblemsSolved(c *gin.Context) {
	username := c.Param("username")
	data, err := services.GetLeetCodeProblemsSolved(username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func GetLeetCodeCalender(c *gin.Context) {
	username := c.Param("username")
	yearStr := c.Query("year")

	var year *int
	if yearStr != "" {
		if y, err := strconv.Atoi(yearStr); err == nil {
			year = &y
		}
	}

	data, err := services.GetLeetCodeCalender(username, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}

func GetLeetCodeContests(c *gin.Context) {
	data, err := services.GetLeetCodeContests()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
