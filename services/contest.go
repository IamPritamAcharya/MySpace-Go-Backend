package services

import (
	"context"
	"fmt"
	"myspace-backend/models"
	"sort"
	"sync"
	"time"
)

var (
	istLocation *time.Location
	initOnce    sync.Once
)

func init() {
	initOnce.Do(func() {
		var err error
		istLocation, err = time.LoadLocation("Asia/Kolkata")
		if err != nil {
			istLocation = time.UTC
		}
	})
}

func parseCodechefDate(dateStr string) (int64, error) {
	t, err := time.ParseInLocation("02 Jan 2006 15:04:05", dateStr, istLocation)
	if err != nil {
		return 0, err
	}
	return t.Unix(), nil
}


type ContestFetcher func(ctx context.Context, resultChan chan<- ContestResult)

type ContestResult struct {
	Contests []models.ContestDTO
	Error    error
	Platform string
}

func fetchCodeforcesContests(ctx context.Context, resultChan chan<- ContestResult) {
	defer func() {
		if r := recover(); r != nil {
			resultChan <- ContestResult{Error: fmt.Errorf("codeforces panic: %v", r), Platform: "Codeforces"}
		}
	}()

	select {
	case <-ctx.Done():
		resultChan <- ContestResult{Error: ctx.Err(), Platform: "Codeforces"}
		return
	default:
	}

	cfContests, err := GetCodeforcesContests()
	if err != nil {
		resultChan <- ContestResult{Error: err, Platform: "Codeforces"}
		return
	}

	contests := make([]models.ContestDTO, 0, len(cfContests))
	for _, contest := range cfContests {
		contests = append(contests, models.ContestDTO{
			Name:      contest.Name,
			Platform:  "Codeforces",
			StartTime: contest.StartTimeSeconds,
			Link:      fmt.Sprintf("https://codeforces.com/contest/%d", contest.ID),
		})
	}

	resultChan <- ContestResult{Contests: contests, Platform: "Codeforces"}
}

func fetchLeetCodeContests(ctx context.Context, resultChan chan<- ContestResult) {
	defer func() {
		if r := recover(); r != nil {
			resultChan <- ContestResult{Error: fmt.Errorf("leetcode panic: %v", r), Platform: "LeetCode"}
		}
	}()

	select {
	case <-ctx.Done():
		resultChan <- ContestResult{Error: ctx.Err(), Platform: "LeetCode"}
		return
	default:
	}

	lcContests, err := GetLeetCodeContests()
	if err != nil {
		resultChan <- ContestResult{Error: err, Platform: "LeetCode"}
		return
	}

	var contests []models.ContestDTO

	response, ok := lcContests.(models.LeetCodeResponse)
	if !ok {
		resultChan <- ContestResult{Error: fmt.Errorf("invalid leetcode response format"), Platform: "LeetCode"}
		return
	}

	data, ok := response.Data.(map[string]interface{})
	if !ok {
		resultChan <- ContestResult{Error: fmt.Errorf("invalid leetcode data format"), Platform: "LeetCode"}
		return
	}

	upcomingContests, ok := data["upcomingContests"].([]interface{})
	if !ok {
		resultChan <- ContestResult{Contests: contests, Platform: "LeetCode"} 
		return
	}

	contests = make([]models.ContestDTO, 0, len(upcomingContests))

	for _, contest := range upcomingContests {
		c, ok := contest.(map[string]interface{})
		if !ok {
			continue
		}

		title, titleOk := c["title"].(string)
		titleSlug, slugOk := c["titleSlug"].(string)
		startTime, timeOk := c["startTime"].(float64)

		if titleOk && slugOk && timeOk {
			contests = append(contests, models.ContestDTO{
				Name:      title,
				Platform:  "LeetCode",
				StartTime: int64(startTime),
				Link:      fmt.Sprintf("https://leetcode.com/contest/%s", titleSlug),
			})
		}
	}

	resultChan <- ContestResult{Contests: contests, Platform: "LeetCode"}
}

func fetchCodechefContests(ctx context.Context, resultChan chan<- ContestResult) {
	defer func() {
		if r := recover(); r != nil {
			resultChan <- ContestResult{Error: fmt.Errorf("codechef panic: %v", r), Platform: "CodeChef"}
		}
	}()

	select {
	case <-ctx.Done():
		resultChan <- ContestResult{Error: ctx.Err(), Platform: "CodeChef"}
		return
	default:
	}

	ccContests, err := FetchCodechefContests()
	if err != nil {
		resultChan <- ContestResult{Error: err, Platform: "CodeChef"}
		return
	}


	contests := make([]models.ContestDTO, 0, len(ccContests))

	for _, contest := range ccContests {
		startTime, parseErr := parseCodechefDate(contest.StartDate)
		if parseErr != nil {
			continue 
		}

		contests = append(contests, models.ContestDTO{
			Name:      contest.ContestName,
			Platform:  "CodeChef",
			StartTime: startTime,
			Link:      fmt.Sprintf("https://www.codechef.com/%s", contest.ContestCode),
		})
	}

	resultChan <- ContestResult{Contests: contests, Platform: "CodeChef"}
}

func GetAllContests() ([]models.ContestDTO, error) {
	return GetAllContestsWithTimeout(30 * time.Second)
}

func GetAllContestsWithTimeout(timeout time.Duration) ([]models.ContestDTO, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()


	resultChan := make(chan ContestResult, 3)

	go fetchCodeforcesContests(ctx, resultChan)
	go fetchLeetCodeContests(ctx, resultChan)
	go fetchCodechefContests(ctx, resultChan)

	var allContests []models.ContestDTO
	estimatedCapacity := 50 
	allContests = make([]models.ContestDTO, 0, estimatedCapacity)

	var errors []error
	resultsReceived := 0

	for resultsReceived < 3 {
		select {
		case result := <-resultChan:
			resultsReceived++
			if result.Error != nil {
				errors = append(errors, fmt.Errorf("%s: %w", result.Platform, result.Error))
			} else {
				allContests = append(allContests, result.Contests...)
			}
		case <-ctx.Done():
			return nil, fmt.Errorf("timeout while fetching contests: %w", ctx.Err())
		}
	}

	validContests := make([]models.ContestDTO, 0, len(allContests))
	for _, contest := range allContests {
		if contest.StartTime > 0 {
			validContests = append(validContests, contest)
		}
	}

	sort.Slice(validContests, func(i, j int) bool {
		return validContests[i].StartTime < validContests[j].StartTime
	})

	if len(validContests) == 0 && len(errors) > 0 {
		return nil, fmt.Errorf("all contest APIs failed: %v", errors)
	}

	return validContests, nil
}
