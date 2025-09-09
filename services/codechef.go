package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"myspace-backend/models"
	"net/http"
)

func FetchCodechefContests() ([]models.CodechefContest, error) {
	url := "https://www.codechef.com/api/list/contests/all?sort_by=START&sorting_order=asc&offset=0&mode=all"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte{}))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to fetch codechef contests")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsed models.CodechefAPIResponse
	if err := json.Unmarshal(body, &parsed); err != nil {
		return nil, err
	}

	return parsed.FutureContests, nil
}