package services

import (
	"fmt"
	"myspace-backend/models"
	"myspace-backend/utils"
)

func GetCodeforcesRating(handle string) ([]models.RatingChange, error) {
	var response models.CodeforcesResponse
	_, err := utils.Client.R().
		SetResult(&response).
		Get(fmt.Sprintf("https://codeforces.com/api/user.rating?handle=%s", handle))

	if err != nil {
		return nil, err
	}


	ratings := []models.RatingChange{}
	if result, ok := response.Result.([]interface{}); ok {
		for _, item := range result {
			if rating, ok := item.(map[string]interface{}); ok {
				ratings = append(ratings, models.RatingChange{
					ContestName:             rating["contestName"].(string),
					RatingUpdateTimeSeconds: int(rating["ratingUpdateTimeSeconds"].(float64)),
					NewRating:               int(rating["newRating"].(float64)),
				})
			}
		}
	}
	return ratings, nil
}

func GetCodeforcesContests() ([]models.CodeforcesContest, error) {
	var response models.CodeforcesResponse
	_, err := utils.Client.R().
		SetResult(&response).
		Get("https://codeforces.com/api/contest.list?gym=false")

	if err != nil {
		return nil, err
	}

	contests := []models.CodeforcesContest{}
	if result, ok := response.Result.([]interface{}); ok {
		for _, item := range result {
			if contest, ok := item.(map[string]interface{}); ok {
				if contest["phase"].(string) == "BEFORE" {
					contests = append(contests, models.CodeforcesContest{
						ID:               int(contest["id"].(float64)),
						Name:             contest["name"].(string),
						Phase:            contest["phase"].(string),
						StartTimeSeconds: int64(contest["startTimeSeconds"].(float64)),
					})
				}
			}
		}
	}
	return contests, nil
}