package models

type Stack struct {
	Title      string              `json:"title"`
	ToolCount  int                 `json:"toolCount"`
	Categories map[string][]string `json:"categories"`
}

type CompanyStack struct {
	Company string  `json:"company"`
	Stacks  []Stack `json:"stacks"`
}
