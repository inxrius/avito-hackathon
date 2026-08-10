package narrative

import "context"

type SafeFact struct {
	MetricCode string `json:"metric_code"`
	Value      int    `json:"value"`
	Unit       string `json:"unit"`
}

type NamedCode struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

type AchievementFact struct {
	Code  string `json:"code"`
	Title string `json:"title"`
	Level string `json:"level"`
}

type Input struct {
	Locale       string            `json:"locale"`
	Year         int               `json:"year"`
	Theme        string            `json:"theme"`
	Role         NamedCode         `json:"role"`
	Style        NamedCode         `json:"style"`
	MainDistrict NamedCode         `json:"main_district"`
	TopCategory  *NamedCode        `json:"top_category,omitempty"`
	Achievements []AchievementFact `json:"achievements"`
	SafeFacts    []SafeFact        `json:"safe_facts"`
}

type ProviderResult struct {
	Content []byte
	Model   string
}

type Provider interface {
	Generate(ctx context.Context, input Input) (ProviderResult, error)
}
