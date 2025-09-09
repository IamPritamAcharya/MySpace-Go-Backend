package services

import (
	"fmt"
	"myspace-backend/models"
	"myspace-backend/utils"
	"regexp"
	"sort"
	"strings"
)

const newsAPIKey = "pub_9c114068d5204b62b63f7f0e2087bd59"

var techKeywords = map[string]int{
	"javascript": 10, "python": 10, "java": 10, "typescript": 10,
	"react": 10, "angular": 9, "vue": 9, "nodejs": 9,
	"aws": 9, "azure": 9, "docker": 9, "kubernetes": 9,
	"machine learning": 10, "artificial intelligence": 10, "openai": 10,
	"programming": 8, "coding": 8, "software development": 9,
}

var excludeKeywords = []string{"sports", "entertainment", "politics", "weather"}

func GetTechNews() ([]models.News, error) {
	queries := []string{
		"programming OR coding OR software",
		"artificial intelligence OR machine learning",
		"web development OR mobile development",
	}

	var allArticles []models.NewsArticle
	for _, query := range queries {
		articles, _ := fetchNewsWithQuery(query)
		allArticles = append(allArticles, articles...)
	}


	uniqueArticles := removeDuplicates(allArticles)
	filteredArticles := filterTechArticles(uniqueArticles)
	

	sort.Slice(filteredArticles, func(i, j int) bool {
		return calculateRelevance(filteredArticles[i]) > calculateRelevance(filteredArticles[j])
	})
	
	if len(filteredArticles) > 15 {
		filteredArticles = filteredArticles[:15]
	}

	var news []models.News
	for _, article := range filteredArticles {
		summary, _ := SummarizeText(article.Title + ". " + article.Description)
		news = append(news, models.News{
			Title:       article.Title,
			Description: summary,
			URL:         article.Link,
			ImageURL:    article.ImageURL,
			PublishedAt: article.PubDate,
			Source:      article.SourceID,
		})
	}

	return news, nil
}

func fetchNewsWithQuery(query string) ([]models.NewsArticle, error) {
	var response models.NewsResponse
	_, err := utils.Client.R().
		SetResult(&response).
		Get(fmt.Sprintf("https://newsdata.io/api/1/news?apikey=%s&language=en&category=technology&q=%s&size=10", 
			newsAPIKey, query))

	if err != nil {
		return nil, err
	}
	return response.Results, nil
}

func removeDuplicates(articles []models.NewsArticle) []models.NewsArticle {
	seen := make(map[string]bool)
	var unique []models.NewsArticle
	
	for _, article := range articles {
		if !seen[article.Title] {
			seen[article.Title] = true
			unique = append(unique, article)
		}
	}
	return unique
}

func filterTechArticles(articles []models.NewsArticle) []models.NewsArticle {
	var filtered []models.NewsArticle
	for _, article := range articles {
		if calculateRelevance(article) >= 15 {
			filtered = append(filtered, article)
		}
	}
	return filtered
}

func calculateRelevance(article models.NewsArticle) int {
	content := strings.ToLower(article.Title + " " + article.Description)
	score := 0

	for keyword, weight := range techKeywords {
		if strings.Contains(content, keyword) {
			score += weight
			if strings.Contains(strings.ToLower(article.Title), keyword) {
				score += weight / 2
			}
		}
	}

	techCompanies := regexp.MustCompile(`\b(google|microsoft|apple|amazon|meta|facebook|nvidia|openai)\b`)
	if techCompanies.MatchString(content) {
		score += 5
	}

	for _, exclude := range excludeKeywords {
		if strings.Contains(content, exclude) {
			score -= 20
		}
	}

	if strings.Contains(content, "ai") || strings.Contains(content, "chatgpt") {
		score += 10
	}

	if score < 0 {
		score = 0
	}
	return score
}
