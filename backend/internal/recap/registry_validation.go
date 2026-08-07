package recap

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

func (r Registry) IsZero() bool {
	return r.Verticals == nil && r.Categories == nil && r.Roles == nil && r.Styles == nil &&
				 r.Metrics == nil && r.MetricCardTemplates == nil && r.Achievements == nil &&
				 r.ThemeAccents == nil && r.PublicAvatarHosts == nil
}

func ValidateRegistry(r Registry) error {
	if err := validateVerticals(r); err != nil {
		return err
	}
	if err := validateCategories(r); err != nil {
		return err
	}
	if err := validateRoles(r); err != nil {
		return err
	}
	if err := validateStyles(r); err != nil {
		return err
	}
	if err := validateMetrics(r); err != nil {
		return err
	}
	if err := validateMetricCardTemplates(r); err != nil {
		return err
	}
	if err := validateAchievements(r); err != nil {
		return err
	}
	if err := validateTheme(r); err != nil {
		return err
	}
	for host := range r.PublicAvatarHosts {
		if host == "" || strings.ToLower(host) != host || strings.ContainsAny(host, "/:@?#") {
			return &ConfigError{Code: "invalid_public_avatar_host", Message: host}
		}
	}
	return nil
}

func validateVerticals(r Registry) error {
	for _, code := range []VerticalCode{VerticalGoods, VerticalTransport, VerticalRealEstate, VerticalJobs, VerticalServices} {
		value, ok := r.Verticals[code]
		if !ok || value.Code != code || !validText(value.Title, 1, 100) {
			return &ConfigError{Code: "invalid_vertical_registry", Message: string(code)}
		}
	}
	return nil
}

func validateCategories(r Registry) error {
	expected := map[CategoryCode]VerticalCode{
		CategoryElectronics: VerticalGoods, CategoryHomeAndGarden: VerticalGoods,
		CategoryClothingAndAccessories: VerticalGoods, CategoryHobbiesAndLeisure: VerticalGoods,
		CategoryCars: VerticalTransport, CategoryApartments: VerticalRealEstate,
		CategoryVacancies: VerticalJobs, CategoryPersonalServices: VerticalServices,
	}
	for code, vertical := range expected {
		value, ok := r.Categories[code]
		if !ok || value.Code != code || value.VerticalCode != vertical || !validText(value.Title, 1, 100) {
			return &ConfigError{Code: "invalid_category_registry", Message: string(code)}
		}
	}
	return nil
}

func validateRoles(r Registry) error {
	for _, code := range []ArchetypeRoleCode{RoleFindingsSeeker, RoleShowcaseOwner, RoleUniversalCitizen, RoleCityObserver} {
		value, ok := r.Roles[code]
		if !ok || value.Code != code || !validText(value.Title, 1, 100) {
			return &ConfigError{Code: "invalid_role_registry", Message: string(code)}
		}
	}
	return nil
}

func validateStyles(r Registry) error {
	for _, code := range []ArchetypeStyleCode{StyleThoughtful, StyleExplorer, StyleDistrictExpert, StyleRegular, StyleResultOriented, StyleCityLocal} {
		value, ok := r.Styles[code]
		if !ok || value.Code != code || !validText(value.Title, 1, 100) {
			return &ConfigError{Code: "invalid_style_registry", Message: string(code)}
		}
	}
	return nil
}

func validateMetrics(r Registry) error {
	for code, unit := range defaultMetricDefinitions() {
		value, ok := r.Metrics[code]
		if !ok || value.Code != code || value.Unit != unit.Unit {
			return &ConfigError{Code: "invalid_metric_registry", Message: string(code)}
		}
	}
	return nil
}

func validateMetricCardTemplates(r Registry) error {
	for _, code := range []MetricCode{MetricFavoritesCount, MetricViewsCount, MetricSalesCount, MetricPublishedListingsCount, MetricCompletedDealsCount, MetricActiveDays} {
		value, ok := r.MetricCardTemplates[code]
		if !ok || !validText(value.TitleOne, 1, 200) || !validText(value.TitleFew, 1, 200) || !validText(value.TitleMany, 1, 200) || !validText(value.Description, 1, 500) {
			return &ConfigError{Code: "invalid_metric_card_registry", Message: string(code)}
		}
	}
	return nil
}

func validateAchievements(r Registry) error {
	expected := map[AchievementCode]struct{}{
		AchievementDealMaster: {}, AchievementFindingsCollector: {}, AchievementCityNavigator: {}, AchievementFrequentGuest: {},
		AchievementOldTimer: {}, AchievementDoorstepDelivery: {}, AchievementOwnShowcase: {}, AchievementFindingsHunter: {}, AchievementCityRhythm: {},
	}
	seen := make(map[AchievementCode]struct{}, len(r.Achievements))
	for _, value := range r.Achievements {
		if _, ok := expected[value.Code]; !ok {
			return &ConfigError{Code: "unknown_achievement_registry_code", Message: string(value.Code)}
		}
		if _, duplicate := seen[value.Code]; duplicate {
			return &ConfigError{Code: "duplicate_achievement_registry_code", Message: string(value.Code)}
		}
		seen[value.Code] = struct{}{}
		if !validText(value.Title, 1, 100) || !validText(value.Description, 1, 300) || value.Group == "" {
			return &ConfigError{Code: "invalid_achievement_registry", Message: string(value.Code)}
		}
		if _, ok := r.Metrics[value.MetricCode]; !ok {
			return &ConfigError{Code: "achievement_metric_missing", Message: string(value.Code)}
		}
		for i, threshold := range value.Thresholds {
			if threshold <= 0 || (i > 0 && threshold <= value.Thresholds[i-1]) {
				return &ConfigError{Code: "invalid_achievement_thresholds", Message: string(value.Code)}
			}
		}
	}
	if len(seen) != len(expected) {
		return &ConfigError{Code: "missing_achievement_registry", Message: fmt.Sprintf("got=%d expected=%d", len(seen), len(expected))}
	}
	return nil
}

func validateTheme(r Registry) error {
	for _, vertical := range []VerticalCode{VerticalGoods, VerticalTransport, VerticalRealEstate, VerticalJobs, VerticalServices} {
		_, ok := r.ThemeAccents[vertical]
		if !ok {
			return &ConfigError{Code: "missing_theme_registry", Message: string(vertical)}
		}
	}
	return nil
}

func validText(value string, min, max int) bool {
	count := utf8.RuneCountInString(strings.TrimSpace(value))
	return count >= min && count <= max
}
