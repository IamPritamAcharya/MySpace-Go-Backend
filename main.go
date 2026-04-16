package main

import (
	"log"
	"net/http"
	"os"

	"myspace-backend/handlers"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	r := gin.Default()

	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "*")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.Next()
	})

	api := r.Group("/api")
	{
		api.GET("/codeforces/rating/:handle", handlers.GetCodeforcesRating)
		api.GET("/codeforces/contests", handlers.GetCodeforcesContests)

		api.GET("/leetcode/submissions-recent/:username", handlers.GetLeetCodeRecentSubmissions)
		api.GET("/leetcode/problem-solved/:username", handlers.GetLeetCodeProblemsSolved)
		api.GET("/leetcode/submissions-calender/:username", handlers.GetLeetCodeCalender)
		api.GET("/leetcode/contests", handlers.GetLeetCodeContests)

		api.GET("/gfg/heatmap/:userHandle", handlers.GetGFGHeatmap)
		api.GET("/contests", handlers.GetAllContests)
		api.GET("/codechef/contests", handlers.GetCodechefContests)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Fatal(r.Run("0.0.0.0:" + port))
}
