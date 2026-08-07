package assembly

import (
	"fmt"
	"strings"

	recap "recap-personalization/internal/recap"
)

func BuildTheme(verticalCode string, registry recap.Registry) (recap.RecapTheme, error) {
	vertical, ok := registry.Verticals[recap.VerticalCode(verticalCode)]
	if !ok {
		return recap.RecapTheme{}, &recap.ConfigError{Code: "missing_vertical_registry", Message: verticalCode}
	}
	accent, ok := registry.ThemeAccents[vertical.Code]
	if !ok {
		return recap.RecapTheme{}, &recap.ConfigError{Code: "missing_theme_registry", Message: verticalCode}
	}
	return recap.RecapTheme{Code: "city", MainDistrict: vertical, AccentToken: &accent}, nil
}

func BuildCards(year int, metrics recap.Metrics, geography recap.Geography, archetype recap.ArchetypeDecision,
								topAchievements []recap.AchievementDecision, totalAchievements int, narrative recap.Narrative,
								capabilities recap.Capabilities, registry recap.Registry) ([]recap.RecapCard, error) {
	vertical, ok := registry.Verticals[recap.VerticalCode(geography.MainVerticalCode)]
	if !ok {
		return nil, &recap.ConfigError{Code: "missing_vertical_registry", Message: geography.MainVerticalCode}
	}
	var category *recap.Category
	if geography.TopCategoryCode != "" {
		value, exists := registry.Categories[recap.CategoryCode(geography.TopCategoryCode)]
		if !exists {
			return nil, &recap.ConfigError{Code: "missing_category_registry", Message: geography.TopCategoryCode}
		}
		category = &value
	}

	metricCode, metricValue := roleMetric(archetype.Role.Code, metrics)
	metricTitle, metricDescription, err := metricText(registry, metricCode, metricValue)
	if err != nil {
		return nil, err
	}
	metricDefinition, ok := registry.Metrics[metricCode]
	if !ok {
		return nil, &recap.ConfigError{Code: "missing_metric_registry", Message: string(metricCode)}
	}

	districtDescription := "Именно здесь проходила большая часть твоего маршрута"
	districtAsset := geography.MainVerticalCode
	if category != nil {
		districtDescription = fmt.Sprintf("Чаще всего маршрут проходил по улице «%s»", category.Title)
		districtAsset += "-" + geography.TopCategoryCode
	}
	activeMonths := activeMonthsLabel(metrics.ActiveMonths)
	finalDescription := "Сохрани итоги или поделись городскими званиями"
	if !capabilities.ShareAvailable {
		finalDescription = "Сохрани итоги своего года в городе"
	}

	cards := []recap.RecapCard{
		recap.IntroCard{
			BaseRecapCard: baseCard(0, "intro", recap.CardTypeIntro, recap.VisibilityShareable, false,
				"Твой город за год готов", "Посмотрим, какие районы, маршруты и достижения запомнились сильнее всего",
				recap.VisualSkyline, "city-intro"),
			Data: recap.IntroCardData{Year: year},
		},
		recap.MetricCard{
			BaseRecapCard: baseCard(1, "active-days", recap.CardTypeMetric, recap.VisibilityShareable, false,
				activeDaysTitle(metrics.ActiveDays), "", recap.VisualCalendar, "active-days"),
			Data: recap.MetricValue{MetricCode: recap.MetricActiveDays, Value: float64(metrics.ActiveDays), Unit: recap.UnitDays, SecondaryLabel: &activeMonths},
		},
		recap.MetricCard{
			BaseRecapCard: baseCard(2, "role-highlight", recap.CardTypeMetric, recap.VisibilityShareable, false,
				metricTitle, metricDescription, recap.VisualChart, "metric-"+strings.ReplaceAll(string(metricCode), "_", "-")),
			Data: recap.MetricValue{MetricCode: metricCode, Value: float64(metricValue), Unit: metricDefinition.Unit},
		},
		recap.DistrictCard{
			BaseRecapCard: baseCard(3, "main-district", recap.CardTypeDistrict, recap.VisibilityShareable, false,
				fmt.Sprintf("Твой главный район - %s", vertical.Title), districtDescription, recap.VisualDistrict, districtAsset),
			Data: recap.DistrictCardData{Vertical: vertical, ActivityShare: metrics.TopVerticalShare, TopCategory: category},
		},
		recap.ArchetypeCard{
			BaseRecapCard: baseCardWithEyebrow(4, "archetype", recap.CardTypeArchetype, recap.VisibilityShareable, capabilities.ExplanationAvailable,
				"Твоя роль в городе", archetype.Role.Title, fmt.Sprintf("Твой стиль - %s", archetype.Style.Title),
				recap.VisualCharacter, strings.ReplaceAll(string(archetype.Style.Code)+"-"+string(archetype.Role.Code), "_", "-")),
			Data: recap.ArchetypeCardData{Role: archetype.Role, Style: archetype.Style},
		},
		recap.AchievementsCard{
			BaseRecapCard: baseCard(5, "achievements", recap.CardTypeAchievements, recap.VisibilityShareable, capabilities.ExplanationAvailable,
				"Твои городские звания", "Достижения, которые лучше всего описывают твой год", recap.VisualBadge, "top-achievements"),
			Data: recap.AchievementsCardData{Items: publicAchievements(topAchievements), TotalCount: totalAchievements},
		},
		recap.SummaryCard{
			BaseRecapCard: baseCard(6, "summary", recap.CardTypeSummary, recap.VisibilityPersonal, false,
				narrative.SummaryTitle, narrative.SummaryText, recap.VisualIllustration, "city-summary"),
			Data: recap.SummaryCardData{RoleCode: archetype.Role.Code, StyleCode: archetype.Style.Code, AchievementCodes: achievementCodes(topAchievements)},
		},
		recap.FinalCard{
			BaseRecapCard: baseCard(7, "final", recap.CardTypeFinal, recap.VisibilityShareable, false,
				"До встречи в городе", finalDescription, recap.VisualSkyline, "city-final"),
			Data: recap.FinalCardData{ShowShareButton: capabilities.ShareAvailable, ShowFeedback: capabilities.FeedbackAvailable},
		},
	}
	return cards, nil
}

func roleMetric(role recap.ArchetypeRoleCode, m recap.Metrics) (recap.MetricCode, int) {
	switch role {
	case recap.RoleFindingsSeeker:
		if m.FavoritesCount > 0 {
			return recap.MetricFavoritesCount, m.FavoritesCount
		}
		if m.ViewsCount > 0 {
			return recap.MetricViewsCount, m.ViewsCount
		}
		return recap.MetricActiveDays, m.ActiveDays
	case recap.RoleCityObserver:
		return recap.MetricViewsCount, m.ViewsCount
	case recap.RoleShowcaseOwner:
		if m.SalesCount > 0 {
			return recap.MetricSalesCount, m.SalesCount
		}
		return recap.MetricPublishedListingsCount, m.PublishedListingsCount
	case recap.RoleUniversalCitizen:
		return recap.MetricCompletedDealsCount, m.CompletedDealsCount
	default:
		return recap.MetricActiveDays, m.ActiveDays
	}
}

func publicAchievements(values []recap.AchievementDecision) []recap.Achievement {
	result := make([]recap.Achievement, 0, len(values))
	for _, value := range values {
		currentValue := float64(value.Value)
		var next *float64
		if value.NextLevelThreshold != nil {
			converted := float64(*value.NextLevelThreshold)
			next = &converted
		}
		result = append(result, recap.Achievement{
			Code: value.Code, Title: value.Title, Description: value.Description, Level: value.Level, Icon: value.Icon,
			MetricCode: value.MetricCode, CurrentValue: &currentValue, NextLevelThreshold: next,
		})
	}
	return result
}

func achievementCodes(values []recap.AchievementDecision) []recap.AchievementCode {
	result := make([]recap.AchievementCode, 0, len(values))
	for _, value := range values {
		result = append(result, value.Code)
	}
	return result
}

func baseCard(position int, id string, cardType recap.CardType, visibility recap.CardVisibility, explainable bool, title, description string, visualKind recap.CardVisualKind, assetCode string) recap.BaseRecapCard {
	card := recap.BaseRecapCard{ID: id, Type: cardType, Position: position, Visibility: visibility, Title: title, Explainable: explainable}
	if description != "" {
		card.Description = &description
	}
	card.Visual = &recap.CardVisual{Kind: visualKind, AssetCode: stringPointer(assetCode)}
	return card
}

func baseCardWithEyebrow(position int, id string, cardType recap.CardType, visibility recap.CardVisibility, explainable bool, eyebrow, title, description string, visualKind recap.CardVisualKind, assetCode string) recap.BaseRecapCard {
	card := baseCard(position, id, cardType, visibility, explainable, title, description, visualKind, assetCode)
	card.Eyebrow = &eyebrow
	return card
}

func stringPointer(value string) *string {
	return &value
}
