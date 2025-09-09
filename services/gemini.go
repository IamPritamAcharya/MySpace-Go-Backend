package services

import (
	"myspace-backend/utils"
)

const geminiAPIKey = "AIzaSyBavhit4LJmccDPfFsh3Lk1HQjYywnO0Kg"
const geminiURL = "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent"

func SummarizeText(text string) (string, error) {
	if text == "" {
		return "No content to summarize", nil
	}

	requestBody := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
					{
						"text": "Summarize this tech news in 2-3 sentences for mobile app: " + text,
					},
				},
			},
		},
	}

	var response map[string]interface{}
	_, err := utils.Client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(requestBody).
		SetResult(&response).
		Post(geminiURL + "?key=" + geminiAPIKey)

	if err != nil {
		return "Summary not available", nil
	}

	if candidates, ok := response["candidates"].([]interface{}); ok && len(candidates) > 0 {
		if candidate, ok := candidates[0].(map[string]interface{}); ok {
			if content, ok := candidate["content"].(map[string]interface{}); ok {
				if parts, ok := content["parts"].([]interface{}); ok && len(parts) > 0 {
					if part, ok := parts[0].(map[string]interface{}); ok {
						if text, ok := part["text"].(string); ok {
							return text, nil
						}
					}
				}
			}
		}
	}

	return "Summary not available", nil
}
