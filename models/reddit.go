package models

import "time"

type RedditPost struct {
	Title     string `json:"title"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
}

type RedditResponse struct {
	Data struct {
		Children []struct {
			Data RedditPost `json:"data"`
		} `json:"children"`
	} `json:"data"`
}

type Summary struct {
	Title     string `json:"title"`
	Summary   string `json:"summary"`
	URL       string `json:"url"`
	Thumbnail string `json:"thumbnail"`
}

type FirestoreData struct {
	Summaries   []Summary `json:"summaries"`
	LastUpdated time.Time `json:"lastUpdated"`
}

type GeminiRequest struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}
