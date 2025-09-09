package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"myspace-backend/models"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/api/option"
)

var firestoreClient *firestore.Client

func InitFirestore() error {
	ctx := context.Background()
	credPath := os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")
	if credPath == "" {
		credPath = "firebase/serviceAccountKey.json"
	}

	app, err := firebase.NewApp(ctx, nil, option.WithCredentialsFile(credPath))
	if err != nil {
		return err
	}

	firestoreClient, err = app.Firestore(ctx)
	if err != nil {
		return err
	}

	go func() {
		if shouldUpdate() {
			log.Println("Data is stale, updating...")
			updateData()
		} else {
			log.Println("Data is fresh, skipping update")
		}
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for range ticker.C {
			log.Println("24-hour timer triggered, updating data...")
			updateData()
		}
	}()

	return nil
}

func shouldUpdate() bool {
	log.Println("Checking if update is needed...")
	doc, err := firestoreClient.Collection("reddit-summaries").Doc("daily-data").Get(context.Background())
	if err != nil {
		log.Println("No existing data found, will update")
		return true
	}
	var data models.FirestoreData
	if err := doc.DataTo(&data); err != nil {
		log.Printf("Error reading Firestore data: %v", err)
		return true
	}
	timeSince := time.Since(data.LastUpdated)
	if timeSince > 24*time.Hour {
		log.Printf("Data is %v old, will update", timeSince)
		return true
	}
	log.Printf("Data is only %v old, skipping update", timeSince)
	return false
}

func updateData() {
	log.Println("Updating Firestore data...")
	posts, err := fetchPosts()
	if err != nil {
		log.Printf("Error fetching posts: %v", err)
		return
	}
	log.Printf("Found %d Reddit posts", len(posts))

	var summaries []models.Summary
	for i, post := range posts {
		log.Printf("Processing post %d/%d: %s", i+1, len(posts), post.Title)
		content := post.Title
		thumbnailURL := post.Thumbnail

		if strings.HasPrefix(post.URL, "http") && !strings.Contains(post.URL, "reddit.com") {
			log.Printf("Fetching article content from: %s", post.URL)
			if article := fetchArticle(post.URL); article != "" {
				content = post.Title + " " + article
				log.Printf("Article content fetched, length: %d characters", len(article))

				if thumbnailURL == "" || thumbnailURL == "self" || thumbnailURL == "default" || thumbnailURL == "nsfw" {
					log.Printf("Extracting thumbnail from article content...")
					thumbnailURL = extractThumbnailFromHTML(article)
					if thumbnailURL != "" {
						log.Printf("Found thumbnail in article: %s", thumbnailURL)
					}
				}
			} else {
				log.Printf("Failed to fetch article content")
			}
		}

		if thumbnailURL == "" || thumbnailURL == "self" || thumbnailURL == "default" || thumbnailURL == "nsfw" {
			thumbnailURL = generatePlaceholderThumbnail(post.Title)
			log.Printf("Using placeholder thumbnail: %s", thumbnailURL)
		}

		if len(content) > 20000 {
			content = content[:20000]
			log.Printf("Truncated content to 20000 characters")
		}

		log.Printf("Attempting to summarize...")
		summary := summarize(content)
		if summary != "" && len(summary) > 20 {
			summaries = append(summaries, models.Summary{
				Title:     post.Title,
				Summary:   summary,
				URL:       post.URL,
				Thumbnail: thumbnailURL,
			})
			log.Printf("Successfully processed post: %s", post.Title)
		} else {
			log.Printf("Failed to generate valid summary for: %s", post.Title)
		}
	}

	if len(summaries) > 0 {
		_, err := firestoreClient.Collection("reddit-summaries").Doc("daily-data").Set(context.Background(), models.FirestoreData{
			Summaries:   summaries,
			LastUpdated: time.Now(),
		})
		if err != nil {
			log.Printf("Error updating Firestore: %v", err)
		} else {
			log.Printf("Firestore data updated successfully with %d summaries", len(summaries))
		}
	} else {
		log.Println("No summaries generated, skipping Firestore update")
	}
}

func GetSummariesFromFirestore() ([]models.Summary, error) {
	doc, err := firestoreClient.Collection("reddit-summaries").Doc("daily-data").Get(context.Background())
	if err != nil {
		return nil, err
	}
	var data models.FirestoreData
	if err := doc.DataTo(&data); err != nil {
		return nil, err
	}
	return data.Summaries, nil
}

func fetchPosts() ([]models.RedditPost, error) {
	req, err := http.NewRequest("GET", "https://www.reddit.com/r/programming/top.json?t=week&limit=20", nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", "MySpaceBot/1.0")

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("reddit API returned status %d", resp.StatusCode)
	}

	var redditResp models.RedditResponse
	if err := json.NewDecoder(resp.Body).Decode(&redditResp); err != nil {
		return nil, err
	}

	posts := make([]models.RedditPost, len(redditResp.Data.Children))
	for i, child := range redditResp.Data.Children {
		posts[i] = child.Data
	}
	return posts, nil
}

func fetchArticle(url string) string {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return ""
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	client := &http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ""
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

func extractThumbnailFromHTML(html string) string {

	patterns := []string{

		`<meta\s+property="og:image"\s+content="([^"]+)"`,
		`<meta\s+content="([^"]+)"\s+property="og:image"`,

		`<meta\s+name="twitter:image"\s+content="([^"]+)"`,
		`<meta\s+content="([^"]+)"\s+name="twitter:image"`,

		`<meta\s+name="image"\s+content="([^"]+)"`,
		`<meta\s+content="([^"]+)"\s+name="image"`,

		`<link\s+rel="image_src"\s+href="([^"]+)"`,

		`<img[^>]+src="([^"]+)"[^>]*(?:width|height)`,

		`<img[^>]+src="([^"]+)"`,
	}

	for _, pattern := range patterns {
		re := regexp.MustCompile(`(?i)` + pattern)
		matches := re.FindStringSubmatch(html)
		if len(matches) > 1 {
			imageURL := strings.TrimSpace(matches[1])
			if isValidImageURL(imageURL) {

				if strings.HasPrefix(imageURL, "//") {
					imageURL = "https:" + imageURL
				} else if strings.HasPrefix(imageURL, "/") {

					continue
				}
				return imageURL
			}
		}
	}

	return ""
}

func isValidImageURL(url string) bool {
	if url == "" {
		return false
	}

	lowerURL := strings.ToLower(url)
	imageExtensions := []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".svg", ".bmp"}

	for _, ext := range imageExtensions {
		if strings.Contains(lowerURL, ext) {
			return true
		}
	}

	imageIndicators := []string{
		"image", "img", "photo", "pic", "thumbnail", "thumb", "avatar", "logo",
		"imgur.com", "i.redd.it", "github.com", "cdn", "static",
	}

	for _, indicator := range imageIndicators {
		if strings.Contains(lowerURL, indicator) {
			return true
		}
	}

	return false
}

func generatePlaceholderThumbnail(title string) string {

	hash := 0
	for _, char := range title {
		hash = int(char) + ((hash << 5) - hash)
	}

	if hash < 0 {
		hash = -hash
	}

	colors := []string{
		"3498db", "e74c3c", "2ecc71", "f39c12", "9b59b6",
		"1abc9c", "34495e", "e67e22", "95a5a6", "d35400",
	}

	color := colors[hash%len(colors)]

	encodedTitle := strings.ReplaceAll(title, " ", "+")
	if len(encodedTitle) > 50 {
		encodedTitle = encodedTitle[:50]
	}

	return "https://via.placeholder.com/400x200/" + color + "/ffffff?text=" + encodedTitle
}

func summarize(text string) string {
	apiKey := "AIzaSyBavhit4LJmccDPfFsh3Lk1HQjYywnO0Kg"
	if apiKey == "" {
		log.Println("GEMINI_API_KEY environment variable not set")
		return ""
	}

	prompt := "Summarize this content in exactly one paragraph (3-4 sentences max). Write clearly and concisely without formatting or special characters:\n" + text

	reqBody := models.GeminiRequest{
		Contents: []struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		}{
			{
				Parts: []struct {
					Text string `json:"text"`
				}{{Text: prompt}},
			},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		log.Printf("Error marshaling request: %v", err)
		return ""
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	resp, err := client.Post("https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key="+apiKey, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		log.Printf("Error calling Gemini API: %v", err)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Gemini API returned status %d", resp.StatusCode)
		return ""
	}

	var gemResp models.GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&gemResp); err != nil {
		log.Printf("Error decoding Gemini response: %v", err)
		return ""
	}

	if len(gemResp.Candidates) > 0 && len(gemResp.Candidates[0].Content.Parts) > 0 {
		summary := strings.TrimSpace(gemResp.Candidates[0].Content.Parts[0].Text)
		summary = strings.ReplaceAll(summary, "*", "")
		summary = strings.ReplaceAll(summary, "#", "")
		return summary
	}
	return ""
}
