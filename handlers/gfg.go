package handlers

import (
	"myspace-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetGFGHeatmap(c *gin.Context) {
	userHandle := c.Param("userHandle")
	data, err := services.GetGFGHeatmap(userHandle)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, data)
}
