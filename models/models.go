package models

type CodeforcesContest struct {
	ID               int    `json:"id"`
	Name             string `json:"name"`
	Phase            string `json:"phase"`
	StartTimeSeconds int64  `json:"startTimeSeconds"`
}

type RatingChange struct {
	ContestName              string `json:"contestName"`
	RatingUpdateTimeSeconds  int    `json:"ratingUpdateTimeSeconds"`
	NewRating                int    `json:"newRating"`
}


type ContestDTO struct {
	Name      string `json:"name"`
	Platform  string `json:"platform"`
	StartTime int64  `json:"startTime"`
	Link      string `json:"link"`
}

type News struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	ImageURL    string `json:"imageUrl"`
	PublishedAt string `json:"publishedAt"`
	Source      string `json:"source"`
}

type GFGHeatmapResponse struct {
	PageProps struct {
		HeatMapData struct {
			Result map[string]int `json:"result"`
		} `json:"heatMapData"`
	} `json:"pageProps"`
}

type CodeforcesResponse struct {
	Result interface{} `json:"result"`
}

type NewsResponse struct {
	Results []NewsArticle `json:"results"`
}

type NewsArticle struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Link        string `json:"link"`
	ImageURL    string `json:"image_url"`
	PubDate     string `json:"pubDate"`
	SourceID    string `json:"source_id"`
}


type LeetCodeResponse struct {
	Data interface{} `json:"data"`
}

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type CodechefContest struct {
	ContestCode string `json:"contest_code"`
	ContestName string `json:"contest_name"`
	StartDate   string `json:"contest_start_date"`
	EndDate     string `json:"contest_end_date"`
}

type CodechefAPIResponse struct {
	FutureContests  []CodechefContest `json:"future_contests"`
	PresentContests []CodechefContest `json:"present_contests"`
}
