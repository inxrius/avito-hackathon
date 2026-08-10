package assembly

import (
	"fmt"

	recap "recap-personalization/internal/recap"
	"recap-personalization/internal/recap/personalization"
)

func BuildExplanations(recapID, activityHash string, m recap.Metrics, archetype recap.ArchetypeDecision, achievements []recap.AchievementDecision) recap.RecapExplanation {
	scores := personalization.CalculateScores(m)
	decisions := []recap.DecisionExplanation{
		roleExplanation(m, archetype.Role),
		styleExplanation(m, scores, archetype.Style),
	}
	for _, achievement := range achievements {
		decisions = append(decisions, recap.DecisionExplanation{
			CardID: "achievements", Kind: recap.ExplanationAchievement, Code: string(achievement.Code),
			RuleVersion: recap.AchievementRuleVersion,
			Reason:      fmt.Sprintf("Достижение «%s» получено на уровне «%s».", achievement.Title, achievement.Level),
			Facts: []recap.RuleFact{
				fact(achievement.MetricCode, float64(achievement.Value), recap.OperatorGTE, float64(achievement.CurrentThreshold)),
			},
		})
	}
	return recap.RecapExplanation{
		RecapID: recapID, AlgorithmVersion: recap.AlgorithmVersion, ActivityHash: activityHash, Decisions: decisions,
	}
}

func roleExplanation(m recap.Metrics, role recap.ArchetypeRole) recap.DecisionExplanation {
	explanation := recap.DecisionExplanation{
		CardID: "archetype", Kind: recap.ExplanationArchetypeRole, Code: string(role.Code), RuleVersion: recap.ArchetypeRuleVersion,
	}
	switch role.Code {
	case recap.RoleCityObserver:
		explanation.Reason = "Ты много изучал предложения, сохранял интересные варианты и редко переходил к завершённым сделкам."
		explanation.Facts = []recap.RuleFact{
			fact(recap.MetricViewsCount, float64(m.ViewsCount), recap.OperatorGTE, 20),
			fact(recap.MetricCompletedDealsCount, float64(m.CompletedDealsCount), recap.OperatorLTE, 1),
		}
		if m.FavoritesCount >= 5 {
			explanation.Facts = append(explanation.Facts, fact(recap.MetricFavoritesCount, float64(m.FavoritesCount), recap.OperatorGTE, 5))
		} else {
			explanation.Facts = append(explanation.Facts, fact(recap.MetricSavedSearchesCount, float64(m.SavedSearchesCount), recap.OperatorGTE, 3))
		}
	case recap.RoleUniversalCitizen:
		explanation.Reason = "Наблюдаемые метрики поддерживают роль с заметной покупательской и продавцовской активностью."
		explanation.Facts = []recap.RuleFact{
			fact(recap.MetricBuyerActionsCount, float64(m.BuyerActionsCount), recap.OperatorGTE, 1),
			fact(recap.MetricSellerActionsCount, float64(m.SellerActionsCount), recap.OperatorGTE, 1),
			fact(recap.MetricCompletedDealsCount, float64(m.CompletedDealsCount), recap.OperatorGTE, 1),
		}
	case recap.RoleShowcaseOwner:
		explanation.Reason = "Продавцовская активность была заметнее, а собственная витрина регулярно пополнялась."
		explanation.Facts = []recap.RuleFact{
			fact(recap.MetricPublishedListingsCount, float64(m.PublishedListingsCount), recap.OperatorGTE, 2),
		}
		if m.SalesCount > 0 {
			explanation.Facts = append(explanation.Facts, fact(recap.MetricSalesCount, float64(m.SalesCount), recap.OperatorGTE, 1))
		}
	case recap.RoleFindingsSeeker:
		code, value := strongestBuyerMetric(m)
		if value > 0 {
			explanation.Reason = "Среди наблюдаемых действий заметнее всего была активность в покупательских сценариях."
			explanation.Facts = []recap.RuleFact{fact(code, float64(value), recap.OperatorGTE, 1)}
		} else {
			explanation.Reason = "Ни одно специализированное правило роли не совпало; использована нейтральная fallback-роль."
			explanation.Facts = []recap.RuleFact{
				fact(recap.MetricMeaningfulEvents, float64(m.MeaningfulEvents), recap.OperatorGTE, 10),
				fact(recap.MetricActiveDays, float64(m.ActiveDays), recap.OperatorGTE, 7),
			}
		}
	}
	return explanation
}

func styleExplanation(m recap.Metrics, scores personalization.Scores, style recap.ArchetypeStyle) recap.DecisionExplanation {
	explanation := recap.DecisionExplanation{
		CardID: "archetype", Kind: recap.ExplanationArchetypeStyle, Code: string(style.Code), RuleVersion: recap.ArchetypeRuleVersion,
	}
	switch style.Code {
	case recap.StyleDistrictExpert:
		explanation.Reason = "Большая часть активности была сосредоточена в одном главном районе."
		explanation.Facts = []recap.RuleFact{
			fact(recap.MetricMeaningfulEvents, float64(m.MeaningfulEvents), recap.OperatorGTE, 20),
			fact(recap.MetricTopVerticalShare, m.TopVerticalShare, recap.OperatorGTE, 0.65),
		}
	case recap.StyleExplorer:
		explanation.Reason = "Ты исследовал много категорий и несколько районов без одного доминирующего направления."
		explanation.Facts = []recap.RuleFact{
			fact(recap.MetricUniqueCategories, float64(m.UniqueCategories), recap.OperatorGTE, 6),
			fact(recap.MetricUniqueVerticals, float64(m.UniqueVerticals), recap.OperatorGTE, 3),
			fact(recap.MetricTopVerticalShare, m.TopVerticalShare, recap.OperatorLT, 0.65),
		}
	case recap.StyleResultOriented:
		explanation.Reason = "Активность часто приводила к завершённым сделкам."
		explanation.Facts = []recap.RuleFact{
			fact(recap.MetricCompletedDealsCount, float64(m.CompletedDealsCount), recap.OperatorGTE, 4),
		}
		if m.FavoriteToPurchaseRate >= 0.35 {
			explanation.Facts = append(explanation.Facts, fact(recap.MetricFavoriteToPurchaseRate, m.FavoriteToPurchaseRate, recap.OperatorGTE, 0.35))
		}
		if m.PublishToSaleRate >= 0.35 {
			explanation.Facts = append(explanation.Facts, fact(recap.MetricPublishToSaleRate, m.PublishToSaleRate, recap.OperatorGTE, 0.35))
		}
	case recap.StyleThoughtful:
		explanation.Reason = "Ты внимательно изучал и сохранял предложения, не торопясь переходить к покупке."
		explanation.Facts = []recap.RuleFact{
			fact(recap.MetricViewsCount, float64(m.ViewsCount), recap.OperatorGTE, 30),
			fact(recap.MetricFavoritesCount, float64(m.FavoritesCount), recap.OperatorGTE, 8),
			fact(recap.MetricFavoriteToPurchaseRate, m.FavoriteToPurchaseRate, recap.OperatorLTE, 0.20),
		}
	case recap.StyleRegular:
		explanation.Reason = "Ты возвращался в город регулярно на протяжении значительной части года."
		explanation.Facts = []recap.RuleFact{
			fact(recap.MetricActiveMonths, float64(m.ActiveMonths), recap.OperatorGTE, 6),
			fact(recap.MetricActiveDays, float64(m.ActiveDays), recap.OperatorGTE, 30),
		}
	case recap.StyleCityLocal:
		matchedRule := scores.Balance >= 0.65 && m.UniqueCategories >= 3 && m.ActiveMonths >= 3
		if matchedRule {
			explanation.Reason = "Наблюдаемые метрики поддерживают сбалансированный городской стиль в нескольких категориях."
			explanation.Facts = []recap.RuleFact{
				fact(recap.MetricBuyerActionsCount, float64(m.BuyerActionsCount), recap.OperatorGTE, 1),
				fact(recap.MetricSellerActionsCount, float64(m.SellerActionsCount), recap.OperatorGTE, 1),
				fact(recap.MetricUniqueCategories, float64(m.UniqueCategories), recap.OperatorGTE, 3),
				fact(recap.MetricActiveMonths, float64(m.ActiveMonths), recap.OperatorGTE, 3),
			}
		} else {
			explanation.Reason = "Ни один более выраженный стиль не доминировал, поэтому выбран нейтральный городской стиль."
			explanation.Facts = []recap.RuleFact{
				fact(recap.MetricMeaningfulEvents, float64(m.MeaningfulEvents), recap.OperatorGTE, 10),
				fact(recap.MetricActiveDays, float64(m.ActiveDays), recap.OperatorGTE, 7),
			}
		}
	}
	return explanation
}

func strongestBuyerMetric(m recap.Metrics) (recap.MetricCode, int) {
	values := []struct {
		code  recap.MetricCode
		value int
	}{
		{recap.MetricViewsCount, m.ViewsCount},
		{recap.MetricFavoritesCount, m.FavoritesCount},
		{recap.MetricSavedSearchesCount, m.SavedSearchesCount},
		{recap.MetricChatsStartedCount, m.ChatsStartedCount},
		{recap.MetricPurchasesCount, m.PurchasesCount},
	}
	best := values[0]
	for _, candidate := range values[1:] {
		if candidate.value > best.value {
			best = candidate
		}
	}
	return best.code, best.value
}

func fact(code recap.MetricCode, actual float64, operator recap.RuleOperator, threshold float64) recap.RuleFact {
	matched := false
	switch operator {
	case recap.OperatorEQ:
		matched = actual == threshold
	case recap.OperatorNEQ:
		matched = actual != threshold
	case recap.OperatorGTE:
		matched = actual >= threshold
	case recap.OperatorLTE:
		matched = actual <= threshold
	case recap.OperatorLT:
		matched = actual < threshold
	case recap.OperatorGT:
		matched = actual > threshold
	}
	return recap.RuleFact{MetricCode: code, Actual: actual, Operator: operator, Threshold: threshold, Matched: matched}
}
