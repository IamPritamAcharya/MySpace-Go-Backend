
package services

import (
	"myspace-backend/models"
	"myspace-backend/utils"
)

func GetLeetCodeRecentSubmissions(username string) (interface{}, error) {
	query := `query recentSubmissions($username: String!) {
		recentSubmissionList(username: $username) {
			title
			titleSlug
			timestamp
			statusDisplay
			lang
		}
	}`

	request := models.GraphQLRequest{
		Query:     query,
		Variables: map[string]interface{}{"username": username},
	}

	var response models.LeetCodeResponse

	_, err := utils.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		SetHeader("Referer", "https://leetcode.com/").
		SetBody(request).
		SetResult(&response).
		Post("https://leetcode.com/graphql")

	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetLeetCodeProblemsSolved(username string) (interface{}, error) {
	query := `query userProblemsSolved($username: String!) {
		allQuestionsCount {
			difficulty
			count
		}
		matchedUser(username: $username) {
			problemsSolvedBeatsStats {
				difficulty
				percentage
			}
			submitStatsGlobal {
				acSubmissionNum {
					difficulty
					count
				}
			}
		}
	}`

	request := models.GraphQLRequest{
		Query:     query,
		Variables: map[string]interface{}{"username": username},
	}

	var response models.LeetCodeResponse

	_, err := utils.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		SetHeader("Referer", "https://leetcode.com/").
		SetBody(request).
		SetResult(&response).
		Post("https://leetcode.com/graphql")

	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetLeetCodeCalender(username string, year *int) (interface{}, error) {
	query := `query userProfileCalendar($username: String!, $year: Int) {
		matchedUser(username: $username) {
			userCalendar(year: $year) {
				activeYears
				streak
				totalActiveDays
				dccBadges {
					timestamp
					badge {
						name
						icon
					}
				}
				submissionCalendar
			}
		}
	}`

	variables := map[string]interface{}{"username": username}
	if year != nil {
		variables["year"] = *year
	}

	request := models.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	var response models.LeetCodeResponse

	_, err := utils.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		SetHeader("Referer", "https://leetcode.com/").
		SetBody(request).
		SetResult(&response).
		Post("https://leetcode.com/graphql")

	if err != nil {
		return nil, err
	}

	return response, nil
}

func GetLeetCodeContests() (interface{}, error) {
	query := `query {
		upcomingContests {
			title
			titleSlug
			startTime
			duration
		}
	}`

	request := models.GraphQLRequest{
		Query:     query,
		Variables: map[string]interface{}{},
	}

	var response models.LeetCodeResponse

	_, err := utils.Client.R().
		SetHeader("Content-Type", "application/json").
		SetHeader("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36").
		SetHeader("Referer", "https://leetcode.com/").
		SetBody(request).
		SetResult(&response).
		Post("https://leetcode.com/graphql")

	if err != nil {
		return nil, err
	}

	return response, nil
}

