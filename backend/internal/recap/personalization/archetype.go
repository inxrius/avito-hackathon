package personalization

import (
	recap "recap-personalization/internal/recap"
)

type Scores struct {
	Buyer   float64
	Seller  float64
	Balance float64
}

func CalculateScores(m recap.Metrics) Scores {
	buyer := 0.15*norm(m.ViewsCount, 50) + 0.20*norm(m.FavoritesCount, 20) + 0.10*norm(m.SavedSearchesCount, 5) +
		0.15*norm(m.ChatsStartedCount, 10) + 0.40*norm(m.PurchasesCount, 5)
	seller := 0.35*norm(m.PublishedListingsCount, 10) + 0.65*norm(m.SalesCount, 5)
	balance := 0.0
	if buyer+seller != 0 {
		balance = 1 - abs(buyer-seller)/(buyer+seller)
	}
	return Scores{Buyer: buyer, Seller: seller, Balance: balance}
}

func SelectArchetype(m recap.Metrics, registry recap.Registry) (recap.ArchetypeDecision, Scores, error) {
	scores := CalculateScores(m)
	roleCode := selectRole(m, scores)
	styleCode := selectStyle(m, scores)
	role, ok := registry.Roles[roleCode]
	if !ok {
		return recap.ArchetypeDecision{}, scores, &recap.ConfigError{Code: "missing_role_registry", Message: string(roleCode)}
	}
	style, ok := registry.Styles[styleCode]
	if !ok {
		return recap.ArchetypeDecision{}, scores, &recap.ConfigError{Code: "missing_style_registry", Message: string(styleCode)}
	}
	return recap.ArchetypeDecision{Role: role, Style: style}, scores, nil
}

func selectRole(m recap.Metrics, s Scores) recap.ArchetypeRoleCode {
	if m.ViewsCount >= 20 && (m.FavoritesCount >= 5 || m.SavedSearchesCount >= 3) && m.CompletedDealsCount <= 1 && s.Seller < 0.25 {
		return recap.RoleCityObserver
	}
	if m.CompletedDealsCount >= 1 && s.Buyer >= 0.35 && s.Seller >= 0.35 && s.Balance >= 0.65 {
		return recap.RoleUniversalCitizen
	}
	if s.Seller > s.Buyer && m.PublishedListingsCount >= 2 {
		return recap.RoleShowcaseOwner
	}
	return recap.RoleFindingsSeeker
}

func selectStyle(m recap.Metrics, s Scores) recap.ArchetypeStyleCode {
	if m.MeaningfulEvents >= 20 && m.TopVerticalShare >= 0.65 {
		return recap.StyleDistrictExpert
	}
	if m.UniqueCategories >= 6 && m.UniqueVerticals >= 3 && m.TopVerticalShare < 0.65 {
		return recap.StyleExplorer
	}
	if m.CompletedDealsCount >= 4 && (m.FavoriteToPurchaseRate >= 0.35 || m.PublishToSaleRate >= 0.35) {
		return recap.StyleResultOriented
	}
	if m.ViewsCount >= 30 && m.FavoritesCount >= 8 && m.FavoriteToPurchaseRate <= 0.20 {
		return recap.StyleThoughtful
	}
	if s.Balance >= 0.65 && m.UniqueCategories >= 3 && m.ActiveMonths >= 3 {
		return recap.StyleCityLocal
	}
	if m.ActiveMonths >= 6 && m.ActiveDays >= 30 {
		return recap.StyleRegular
	}
	return recap.StyleCityLocal
}

func norm(value, target int) float64 {
	if value <= 0 {
		return 0
	}
	v := float64(value) / float64(target)
	if v > 1 {
		return 1
	}
	return v
}

func abs(value float64) float64 {
	if value < 0 {
		return -value
	}
	return value
}
