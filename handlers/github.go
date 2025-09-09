package handlers

import (
	"myspace-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGitHubContributions(c *gin.Context) {
	username := c.Query("username")
	token := c.GetHeader("X-GitHub-Token")
	
	if username == "" || token == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and token required"})
		return
	}
	
	data, err := services.GetGitHubContributions(username, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}