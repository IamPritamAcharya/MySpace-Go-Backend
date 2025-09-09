package services

import (
	"fmt"
	"myspace-backend/models"
	"myspace-backend/utils"
	"regexp"
)

func GetGFGHeatmap(userHandle string) (map[string]int, error) {

	htmlResp, err := utils.Client.R().
		Get(fmt.Sprintf("https://www.geeksforgeeks.org/user/%s", userHandle))
	if err != nil {
		return nil, err
	}

	re := regexp.MustCompile(`/_next/static/([a-zA-Z0-9_-]+)/_ssgManifest\.js`)
	matches := re.FindStringSubmatch(htmlResp.String())
	if len(matches) < 2 {
		return nil, fmt.Errorf("build ID not found")
	}
	buildID := matches[1]

	var response models.GFGHeatmapResponse
	_, err = utils.Client.R().
		SetResult(&response).
		Get(fmt.Sprintf("https://www.geeksforgeeks.org/gfg-assets/_next/data/%s/user/%s.json", buildID, userHandle))

	if err != nil {
		return nil, err
	}

	return response.PageProps.HeatMapData.Result, nil
}
