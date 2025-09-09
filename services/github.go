package services

import (
	"fmt"
	"myspace-backend/utils"
)

func GetGitHubContributions(username, token string) (interface{}, error) {
	query := fmt.Sprintf(`{
		"query": "query { user(login: \"%s\") { contributionsCollection { contributionCalendar { weeks { contributionDays { date contributionCount } } } } } }"
	}`, username)

	var response interface{}
	_, err := utils.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("Authorization", "Bearer "+token).
		SetBody(query).
		SetResult(&response).
		Post("https://api.github.com/graphql")

	return response, err
}
