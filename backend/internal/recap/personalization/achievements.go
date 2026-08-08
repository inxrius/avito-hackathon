package personalization

import (
	"sort"
	"strings"

	recap "recap-personalization/internal/recap"
)

func CalculateAchievements(m recap.Metrics, registry recap.Registry) []recap.AchievementDecision {
	result := make([]recap.AchievementDecision, 0, len(registry.Achievements))
	for _, rule := range registry.Achievements {
		value := metricInt(m, rule.MetricCode)
		if value < rule.Thresholds[0] {
			continue
		}
		levelIndex := 0
		for i, threshold := range rule.Thresholds {
			if value >= threshold {
				levelIndex = i
			}
		}
		levels := [...]recap.AchievementLevel{recap.LevelNewcomer, recap.LevelLocal, recap.LevelExpert, recap.LevelGuru}
		level := levels[levelIndex]
		current := rule.Thresholds[levelIndex]
		var next *int
		progress := 1.0
		if levelIndex < len(rule.Thresholds)-1 {
			n := rule.Thresholds[levelIndex+1]
			next = &n
			progress = clamp(float64(value-current)/float64(n-current), 0, 1)
		}
		result = append(result, recap.AchievementDecision{
			Code: rule.Code, Title: rule.Title, Description: rule.Description, MetricCode: rule.MetricCode,
			Value: value, Level: level, CurrentThreshold: current, NextLevelThreshold: next,
			Group: rule.Group, Priority: rule.Priority, Progress: progress,
			Icon: strings.ReplaceAll(string(rule.Code), "_", "-") + "-" + string(level),
		})
	}
	return result
}

func SelectTopAchievements(all []recap.AchievementDecision) []recap.AchievementDecision {
	bestByGroup := map[string]recap.AchievementDecision{}
	for _, achievement := range all {
		current, exists := bestByGroup[achievement.Group]
		if !exists || achievementLess(achievement, current) {
			bestByGroup[achievement.Group] = achievement
		}
	}
	selected := make([]recap.AchievementDecision, 0, len(bestByGroup))
	for _, achievement := range bestByGroup {
		selected = append(selected, achievement)
	}
	sort.Slice(selected, func(i, j int) bool { return achievementLess(selected[i], selected[j]) })
	if len(selected) > 3 {
		selected = selected[:3]
	}
	return selected
}

func achievementLess(a, b recap.AchievementDecision) bool {
	ar, br := levelRank(a.Level), levelRank(b.Level)
	if ar != br {
		return ar > br
	}
	if a.Progress != b.Progress {
		return a.Progress > b.Progress
	}
	if a.Priority != b.Priority {
		return a.Priority > b.Priority
	}
	return a.Code < b.Code
}

func levelRank(level recap.AchievementLevel) int {
	switch level {
	case recap.LevelNewcomer:
		return 1
	case recap.LevelLocal:
		return 2
	case recap.LevelExpert:
		return 3
	case recap.LevelGuru:
		return 4
	default:
		return 0
	}
}

func metricInt(m recap.Metrics, code recap.MetricCode) int {
	switch code {
	case recap.MetricSalesCount:
		return m.SalesCount
	case recap.MetricFavoritesCount:
		return m.FavoritesCount
	case recap.MetricUniqueCategories:
		return m.UniqueCategories
	case recap.MetricActiveDays:
		return m.ActiveDays
	case recap.MetricActiveMonths:
		return m.ActiveMonths
	case recap.MetricDeliveryCount:
		return m.DeliveryCount
	case recap.MetricPublishedListingsCount:
		return m.PublishedListingsCount
	case recap.MetricPurchasesCount:
		return m.PurchasesCount
	case recap.MetricMaxActivityStreak:
		return m.MaxActivityStreak
	default:
		return 0
	}
}

func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
