package pipeline

import (
	"context"
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"unicode/utf8"

	recap "recap-personalization/internal/recap"
	"recap-personalization/internal/recap/analytics"
	"recap-personalization/internal/recap/assembly"
	"recap-personalization/internal/recap/narrative"
	"recap-personalization/internal/recap/personalization"
)

var _ recap.Generator = (*Service)(nil)

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

type Service struct {
	Registry          recap.Registry
	Capabilities      recap.Capabilities
	NarrativeProvider narrative.Provider
}

func NewGenerator(provider narrative.Provider) *Service {
	return &Service{
		Registry: recap.DefaultRegistry(),
		Capabilities: recap.Capabilities{
			ShareAvailable: true, ExplanationAvailable: true, FeedbackAvailable: true,
		},
		NarrativeProvider: provider,
	}
}

func (s *Service) Generate(ctx context.Context, input recap.GenerateInput) (recap.GenerateOutput, error) {
	if s == nil {
		return recap.GenerateOutput{}, &recap.ConfigError{Code: "nil_generator"}
	}
	if err := validateInput(input); err != nil {
		return recap.GenerateOutput{}, err
	}

	registry := s.Registry
	if registry.IsZero() {
		registry = recap.DefaultRegistry()
	}
	if err := recap.ValidateRegistry(registry); err != nil {
		return recap.GenerateOutput{}, err
	}

	normalized, err := analytics.Normalize(input.Profile.ID, input.Year, input.Activities)
	if err != nil {
		return recap.GenerateOutput{}, err
	}
	metrics := analytics.CalculateMetrics(normalized)
	geography := analytics.CalculateGeography(normalized, &metrics)
	if metrics.MeaningfulEvents < 10 || metrics.ActiveDays < 7 || geography.MainVerticalCode == "" {
		return recap.GenerateOutput{}, recap.ErrInsufficientActivity
	}
	activityHash, err := analytics.ActivityHash(normalized)
	if err != nil {
		return recap.GenerateOutput{}, err
	}
	archetype, _, err := personalization.SelectArchetype(metrics, registry)
	if err != nil {
		return recap.GenerateOutput{}, err
	}
	allAchievements := personalization.CalculateAchievements(metrics, registry)
	topAchievements := personalization.SelectTopAchievements(allAchievements)
	if len(topAchievements) == 0 {
		return recap.GenerateOutput{}, &recap.ConfigError{Code: "no_achievements_after_sufficiency"}
	}

	if err := validateResolvedReferences(geography, archetype, topAchievements, registry); err != nil {
		return recap.GenerateOutput{}, err
	}
	theme, err := assembly.BuildTheme(geography.MainVerticalCode, registry)
	if err != nil {
		return recap.GenerateOutput{}, err
	}

	narrativeInput := buildNarrativeInput(input.Year, metrics, geography, archetype, topAchievements, registry)
	topAchievementTitle := topAchievements[0].Title
	narrativeResult := (narrative.Service{Provider: s.NarrativeProvider}).Build(ctx, narrativeInput, topAchievementTitle)

	cards, err := assembly.BuildCards(input.Year, metrics, geography, archetype, topAchievements, len(allAchievements), narrativeResult, s.Capabilities, registry)
	if err != nil {
		return recap.GenerateOutput{}, err
	}

	var explanation *recap.RecapExplanation
	if s.Capabilities.ExplanationAvailable {
		value := assembly.BuildExplanations(input.RecapID, activityHash, metrics, archetype, topAchievements)
		explanation = &value
	}
	var share *recap.ShareCard
	if s.Capabilities.ShareAvailable {
		value, buildErr := assembly.BuildShareCard(input.RecapID, input.Profile, input.Year, archetype, metrics, topAchievements, theme, registry)
		if buildErr != nil {
			return recap.GenerateOutput{}, buildErr
		}
		share = &value
	}

	output := recap.GenerateOutput{
		Recap: recap.RecapSnapshot{
			SchemaVersion: recap.SchemaVersion,
			ID:            input.RecapID,
			ProfileID:     input.Profile.ID,
			Year:          input.Year,
			Profile: recap.RecapProfile{
				ID: input.Profile.ID, Name: input.Profile.Name, AvatarURL: optionalString(input.Profile.AvatarURL),
			},
			Generation: recap.RecapGeneration{
				AlgorithmVersion: recap.AlgorithmVersion, FeatureSchemaVersion: recap.FeatureSchemaVersion,
				ActivityHash: activityHash, GeneratedAt: input.GeneratedAt.UTC(),
				Narrative: recap.NarrativeGeneration{
					Source: narrativeResult.Source, PromptVersion: narrativeResult.PromptVersion, Model: narrativeResult.Model,
				},
			},
			Theme: theme, Cards: cards, Capabilities: s.Capabilities,
		},
		Metrics: analytics.PublicMetrics(metrics, registry), Archetype: archetype, Achievements: topAchievements,
		Narrative: narrativeResult, Explanation: explanation, Share: share,
	}
	if err := ValidateOutput(output, geography, registry); err != nil {
		return recap.GenerateOutput{}, err
	}
	return output, nil
}

func validateInput(input recap.GenerateInput) error {
	if !uuidPattern.MatchString(input.RecapID) {
		return &recap.InputError{Code: "invalid_recap_id", Message: input.RecapID}
	}
	if !uuidPattern.MatchString(input.Profile.ID) {
		return &recap.InputError{Code: "invalid_profile_id", Message: input.Profile.ID}
	}
	if input.Year < 2000 || input.Year > 2100 {
		return &recap.InputError{Code: "invalid_year", Message: fmt.Sprint(input.Year)}
	}
	name := strings.TrimSpace(input.Profile.Name)
	if utf8.RuneCountInString(name) < 1 || utf8.RuneCountInString(name) > 100 {
		return &recap.InputError{Code: "invalid_profile_name"}
	}
	if input.GeneratedAt.IsZero() {
		return &recap.InputError{Code: "empty_generated_at"}
	}
	if input.Profile.AvatarURL != "" {
		if _, err := url.ParseRequestURI(input.Profile.AvatarURL); err != nil {
			return &recap.InputError{Code: "invalid_avatar_url"}
		}
	}
	return nil
}

func validateResolvedReferences(geography recap.Geography, archetype recap.ArchetypeDecision, achievements []recap.AchievementDecision, registry recap.Registry) error {
	vertical, ok := registry.Verticals[recap.VerticalCode(geography.MainVerticalCode)]
	if !ok {
		return &recap.ConfigError{Code: "missing_vertical_registry", Message: geography.MainVerticalCode}
	}
	if geography.TopCategoryCode != "" {
		category, exists := registry.Categories[recap.CategoryCode(geography.TopCategoryCode)]
		if !exists {
			return &recap.ConfigError{Code: "missing_category_registry", Message: geography.TopCategoryCode}
		}
		if category.VerticalCode != vertical.Code {
			return &recap.ConfigError{Code: "category_vertical_mismatch", Message: geography.TopCategoryCode}
		}
	}
	if _, ok := registry.Roles[archetype.Role.Code]; !ok {
		return &recap.ConfigError{Code: "missing_role_registry", Message: string(archetype.Role.Code)}
	}
	if _, ok := registry.Styles[archetype.Style.Code]; !ok {
		return &recap.ConfigError{Code: "missing_style_registry", Message: string(archetype.Style.Code)}
	}
	for _, achievement := range achievements {
		found := false
		for _, definition := range registry.Achievements {
			if definition.Code == achievement.Code {
				found = true
				break
			}
		}
		if !found {
			return &recap.ConfigError{Code: "missing_achievement_registry", Message: string(achievement.Code)}
		}
	}
	return nil
}

func buildNarrativeInput(year int, metrics recap.Metrics, geography recap.Geography, archetype recap.ArchetypeDecision, achievements []recap.AchievementDecision, registry recap.Registry) narrative.Input {
	vertical := registry.Verticals[recap.VerticalCode(geography.MainVerticalCode)]
	input := narrative.Input{
		Locale: "ru-RU", Year: year, Theme: "city",
		Role:         narrative.NamedCode{Code: string(archetype.Role.Code), Title: archetype.Role.Title},
		Style:        narrative.NamedCode{Code: string(archetype.Style.Code), Title: archetype.Style.Title},
		MainDistrict: narrative.NamedCode{Code: string(vertical.Code), Title: vertical.Title},
		Achievements: make([]narrative.AchievementFact, 0, len(achievements)),
	}
	if geography.TopCategoryCode != "" {
		category := registry.Categories[recap.CategoryCode(geography.TopCategoryCode)]
		input.TopCategory = &narrative.NamedCode{Code: string(category.Code), Title: category.Title}
	}
	for _, achievement := range achievements {
		input.Achievements = append(input.Achievements, narrative.AchievementFact{
			Code: string(achievement.Code), Title: achievement.Title, Level: string(achievement.Level),
		})
	}
	input.SafeFacts = safeNarrativeFacts(metrics, achievements, registry)
	return input
}

func safeNarrativeFacts(metrics recap.Metrics, achievements []recap.AchievementDecision, registry recap.Registry) []narrative.SafeFact {
	facts := []struct {
		code  recap.MetricCode
		value int
	}{
		{recap.MetricActiveDays, metrics.ActiveDays},
		{recap.MetricActiveMonths, metrics.ActiveMonths},
		{recap.MetricCompletedDealsCount, metrics.CompletedDealsCount},
		{recap.MetricUniqueCategories, metrics.UniqueCategories},
	}
	result := make([]narrative.SafeFact, 0, len(facts)+len(achievements))
	seen := map[recap.MetricCode]struct{}{}
	for _, item := range facts {
		definition := registry.Metrics[item.code]
		result = append(result, narrative.SafeFact{MetricCode: string(item.code), Value: item.value, Unit: string(definition.Unit)})
		seen[item.code] = struct{}{}
	}
	for _, achievement := range achievements {
		if _, exists := seen[achievement.MetricCode]; exists {
			continue
		}
		definition := registry.Metrics[achievement.MetricCode]
		result = append(result, narrative.SafeFact{
			MetricCode: string(achievement.MetricCode), Value: achievement.Value, Unit: string(definition.Unit),
		})
		seen[achievement.MetricCode] = struct{}{}
	}
	return result
}

func optionalString(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}
