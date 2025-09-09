package handlers

import (
	"myspace-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetStackShare(c *gin.Context) {
	company := c.Query("company")
	if company == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing 'company' query parameter"})
		return
	}

	data, err := services.FetchStack(company)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, data)
}
