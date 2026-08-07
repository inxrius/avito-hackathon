package analytics

var categoryToVertical = map[string]string{
	"electronics":              "goods",
	"home_and_garden":          "goods",
	"clothing_and_accessories": "goods",
	"hobbies_and_leisure":      "goods",
	"cars":                     "transport",
	"apartments":               "real_estate",
	"vacancies":                "jobs",
	"personal_services":        "services",
}

var allowedVerticals = map[string]struct{}{
	"goods": {}, "transport": {}, "real_estate": {}, "jobs": {}, "services": {},
}

func VerticalForCategory(category string) (string, bool) {
	vertical, ok := categoryToVertical[category]
	return vertical, ok
}

func IsAllowedVertical(vertical string) bool {
	_, ok := allowedVerticals[vertical]
	return ok
}
