package assembly

import (
	"fmt"

	recap "recap-personalization/internal/recap"
)

func pluralIndex(value int) int {
	mod100 := value % 100
	if mod100 >= 11 && mod100 <= 14 {
		return 2
	}
	switch value % 10 {
	case 1:
		return 0
	case 2, 3, 4:
		return 1
	default:
		return 2
	}
}

func plural(value int, one, few, many string) string {
	return []string{one, few, many}[pluralIndex(value)]
}

func activeDaysTitle(value int) string {
	return fmt.Sprintf("%d %s в ритме города", value, plural(value, "день", "дня", "дней"))
}

func activeMonthsLabel(value int) string {
	return fmt.Sprintf("Активность была заметна в %d %s", value, plural(value, "месяце", "месяцах", "месяцах"))
}

func metricText(registry recap.Registry, code recap.MetricCode, value int) (string, string, error) {
	template, ok := registry.MetricCardTemplates[code]
	if !ok {
		return "", "", &recap.ConfigError{Code: "missing_metric_card_registry", Message: string(code)}
	}
	formats := []string{template.TitleOne, template.TitleFew, template.TitleMany}
	return fmt.Sprintf(formats[pluralIndex(value)], value), template.Description, nil
}
