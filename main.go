package main

import (
	"log"
	"myspace-backend/handlers"
	"myspace-backend/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)

	err := services.InitFirestore()
	if err != nil {
		log.Fatal("Failed to initialize Firestore:", err)
	}

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

		api.GET("/github/contributions", handlers.GetGitHubContributions)

		api.GET("/gfg/heatmap/:userHandle", handlers.GetGFGHeatmap)

		api.GET("/news/tech", handlers.GetTechNews)

		api.GET("/contests", handlers.GetAllContests)

		api.GET("/codechef/contests", handlers.GetCodechefContests)

		api.GET("/stackshare", handlers.GetStackShare)

		api.GET("/summaries", handlers.GetSummariesHandler)

	}

	log.Println("Server starting on :8080")
	r.Run(":8080")
}
