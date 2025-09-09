package services

import (
	"fmt"
	"myspace-backend/models"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func FetchStack(company string) (*models.CompanyStack, error) {
	url := fmt.Sprintf("https://stackshare.io/companies/%s", company)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch page: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("non-200 status code: %d", resp.StatusCode)
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse HTML: %v", err)
	}

	var stacks []models.Stack

	doc.Find("a.css-187ugz6").Each(func(i int, s *goquery.Selection) {
		stack := models.Stack{
			Categories: make(map[string][]string),
		}

		rawTitle := strings.TrimSpace(s.Find("div.css-xjztrk.title").Text())
		re := regexp.MustCompile(`^(.*?)(\d+)\s+tools$`)
		matches := re.FindStringSubmatch(rawTitle)

		if len(matches) == 3 {
			stack.Title = strings.TrimSpace(matches[1])
			stack.ToolCount, _ = strconv.Atoi(matches[2])
		} else {
			stack.Title = rawTitle
			stack.ToolCount = 0
		}

		var currentCategory string
		s.Find("div.css-1tnqmnz").Children().Each(func(i int, sel *goquery.Selection) {
			if sel.HasClass("css-18akmr2") {
				currentCategory = strings.ToLower(strings.TrimSpace(sel.Text()))
			}
			if sel.HasClass("css-t6kmge") && currentCategory != "" {
				var imgs []string
				sel.Find("img").Each(func(_ int, img *goquery.Selection) {
					if src, exists := img.Attr("src"); exists {
						imgs = append(imgs, src)
					}
				})
				stack.Categories[currentCategory] = imgs
			}
		})

		stacks = append(stacks, stack)
	})

	return &models.CompanyStack{
		Company: company,
		Stacks:  stacks,
	}, nil
}
