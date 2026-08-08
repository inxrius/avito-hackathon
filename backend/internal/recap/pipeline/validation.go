package pipeline

import (
	"encoding/hex"
	"fmt"
	"net"
	"net/url"
	"strings"
	"unicode/utf8"

	recap "recap-personalization/internal/recap"
	"recap-personalization/internal/recap/analytics"
)

func ValidateOutput(output recap.GenerateOutput, geography recap.Geography, registry recap.Registry) error {
	if err := validateCategoryTaxonomy(geography); err != nil {
		return err
	}
	if output.Recap.SchemaVersion != recap.SchemaVersion || !schemaVersionValid(output.Recap.SchemaVersion) {
		return &recap.ConfigError{Code: "invalid_schema_version"}
	}
	if !uuidPattern.MatchString(output.Recap.ID) || !uuidPattern.MatchString(output.Recap.ProfileID) || output.Recap.ProfileID != output.Recap.Profile.ID {
		return &recap.ConfigError{Code: "invalid_recap_identity"}
	}
	if output.Recap.Year < 2000 || output.Recap.Year > 2100 || !validText(output.Recap.Profile.Name, 1, 100) {
		return &recap.ConfigError{Code: "invalid_recap_profile"}
	}
	if output.Recap.Profile.AvatarURL != nil {
		if _, err := url.ParseRequestURI(*output.Recap.Profile.AvatarURL); err != nil {
			return &recap.ConfigError{Code: "invalid_recap_avatar_url"}
		}
	}
	if output.Recap.Generation.AlgorithmVersion != recap.AlgorithmVersion || len(output.Recap.Generation.AlgorithmVersion) > 50 ||
		output.Recap.Generation.FeatureSchemaVersion != recap.FeatureSchemaVersion || len(output.Recap.Generation.FeatureSchemaVersion) > 50 ||
		output.Recap.Generation.GeneratedAt.IsZero() || !validActivityHash(output.Recap.Generation.ActivityHash) {
		return &recap.ConfigError{Code: "invalid_generation_metadata"}
	}
	if output.Recap.Generation.Narrative.Source != output.Narrative.Source ||
		output.Recap.Generation.Narrative.PromptVersion != output.Narrative.PromptVersion ||
		!sameOptionalString(output.Recap.Generation.Narrative.Model, output.Narrative.Model) {
		return &recap.ConfigError{Code: "narrative_metadata_mismatch"}
	}
	if err := validateNarrative(output.Narrative); err != nil {
		return err
	}
	if err := validateTheme(output.Recap.Theme, geography, registry); err != nil {
		return err
	}
	if err := validateCards(output.Recap.Cards, output.Recap.Year, output.Recap.Capabilities, registry); err != nil {
		return err
	}
	if err := validateMetrics(output.Metrics, registry); err != nil {
		return err
	}
	if err := validateAchievements(output.Achievements); err != nil {
		return err
	}
	if output.Archetype.Role.Code == "" || output.Archetype.Style.Code == "" {
		return &recap.ConfigError{Code: "invalid_archetype_result"}
	}
	if output.Recap.Capabilities.ExplanationAvailable {
		if output.Explanation == nil {
			return &recap.ConfigError{Code: "missing_explanation_projection"}
		}
		if err := validateExplanation(*output.Explanation, output, registry); err != nil {
			return err
		}
	} else if output.Explanation != nil {
		return &recap.ConfigError{Code: "disabled_explanation_was_built"}
	}
	if output.Recap.Capabilities.ShareAvailable {
		if output.Share == nil {
			return &recap.ConfigError{Code: "missing_share_projection"}
		}
		if err := validateShare(*output.Share, output, registry); err != nil {
			return err
		}
	} else if output.Share != nil {
		return &recap.ConfigError{Code: "disabled_share_was_built"}
	}
	return nil
}

func validateNarrative(value recap.Narrative) error {
	if value.Source != "mistral" && value.Source != "template" {
		return &recap.ConfigError{Code: "invalid_narrative_source"}
	}
	if value.PromptVersion != recap.PromptVersion || len(value.PromptVersion) > 50 || !validText(value.SummaryTitle, 1, 120) || !validText(value.SummaryText, 1, 500) {
		return &recap.ConfigError{Code: "invalid_narrative"}
	}
	if value.Source == "mistral" {
		if value.Model == nil || !validText(*value.Model, 1, 100) {
			return &recap.ConfigError{Code: "invalid_mistral_model"}
		}
	} else if value.Model != nil {
		return &recap.ConfigError{Code: "template_model_must_be_null"}
	}
	return nil
}

func validateTheme(theme recap.RecapTheme, geography recap.Geography, registry recap.Registry) error {
	if theme.Code != "city" || theme.MainDistrict.Code != recap.VerticalCode(geography.MainVerticalCode) || !validText(theme.MainDistrict.Title, 1, 100) {
		return &recap.ConfigError{Code: "invalid_theme"}
	}
	registered, ok := registry.Verticals[theme.MainDistrict.Code]
	if !ok || registered != theme.MainDistrict || theme.AccentToken == nil {
		return &recap.ConfigError{Code: "invalid_theme_registry_reference"}
	}
	if expected := registry.ThemeAccents[theme.MainDistrict.Code]; expected != *theme.AccentToken {
		return &recap.ConfigError{Code: "invalid_theme_accent"}
	}
	return nil
}

func validateCards(cards []recap.RecapCard, year int, capabilities recap.Capabilities, registry recap.Registry) error {
	if len(cards) != 8 {
		return &recap.ConfigError{Code: "invalid_card_count", Message: fmt.Sprint(len(cards))}
	}
	expectedIDs := []string{"intro", "active-days", "role-highlight", "main-district", "archetype", "achievements", "summary", "final"}
	expectedTypes := []recap.CardType{recap.CardTypeIntro, recap.CardTypeMetric, recap.CardTypeMetric, recap.CardTypeDistrict, recap.CardTypeArchetype, recap.CardTypeAchievements, recap.CardTypeSummary, recap.CardTypeFinal}
	for i, card := range cards {
		if card == nil {
			return &recap.ConfigError{Code: "nil_card", Message: fmt.Sprint(i)}
		}
		base := card.Base()
		if base.Position != i || base.ID != expectedIDs[i] || base.Type != expectedTypes[i] {
			return &recap.ConfigError{Code: "invalid_card_order", Message: base.ID}
		}
		if !validText(base.ID, 1, 100) || !validText(base.Title, 1, 120) || (base.Description != nil && !validText(*base.Description, 0, 500)) ||
			(base.Eyebrow != nil && !validText(*base.Eyebrow, 0, 100)) || !validVisibility(base.Visibility) {
			return &recap.ConfigError{Code: "card_text_out_of_bounds", Message: base.ID}
		}
		if base.Visual == nil || !validVisualKind(base.Visual.Kind) || (base.Visual.AssetCode != nil && !validText(*base.Visual.AssetCode, 0, 100)) {
			return &recap.ConfigError{Code: "invalid_card_visual", Message: base.ID}
		}
		if base.Explainable && base.Type != recap.CardTypeArchetype && base.Type != recap.CardTypeAchievements {
			return &recap.ConfigError{Code: "invalid_explainable_card", Message: base.ID}
		}
		switch typed := card.(type) {
		case recap.IntroCard:
			if typed.Data.Year != year {
				return &recap.ConfigError{Code: "invalid_intro_data"}
			}
		case recap.MetricCard:
			definition, ok := registry.Metrics[typed.Data.MetricCode]
			if !ok || typed.Data.Unit != definition.Unit || (typed.Data.SecondaryLabel != nil && !validText(*typed.Data.SecondaryLabel, 0, 200)) {
				return &recap.ConfigError{Code: "invalid_metric_card_data", Message: string(typed.Data.MetricCode)}
			}
		case recap.DistrictCard:
			if typed.Data.ActivityShare < 0 || typed.Data.ActivityShare > 1 {
				return &recap.ConfigError{Code: "invalid_district_share"}
			}
			vertical, ok := registry.Verticals[typed.Data.Vertical.Code]
			if !ok || vertical != typed.Data.Vertical {
				return &recap.ConfigError{Code: "invalid_district_vertical"}
			}
			if typed.Data.TopCategory != nil {
				category, ok := registry.Categories[typed.Data.TopCategory.Code]
				if !ok || category != *typed.Data.TopCategory || category.VerticalCode != typed.Data.Vertical.Code {
					return &recap.ConfigError{Code: "invalid_district_category"}
				}
			}
		case recap.ArchetypeCard:
			if _, ok := registry.Roles[typed.Data.Role.Code]; !ok {
				return &recap.ConfigError{Code: "invalid_archetype_card_role"}
			}
			if _, ok := registry.Styles[typed.Data.Style.Code]; !ok {
				return &recap.ConfigError{Code: "invalid_archetype_card_style"}
			}
		case recap.AchievementsCard:
			if len(typed.Data.Items) < 1 || len(typed.Data.Items) > 3 || typed.Data.TotalCount < 1 || typed.Data.TotalCount < len(typed.Data.Items) {
				return &recap.ConfigError{Code: "invalid_achievements_card_data"}
			}
			for _, achievement := range typed.Data.Items {
				if !validPublicAchievement(achievement, registry) {
					return &recap.ConfigError{Code: "invalid_achievement_card_item", Message: string(achievement.Code)}
				}
			}
		case recap.SummaryCard:
			if len(typed.Data.AchievementCodes) < 1 || len(typed.Data.AchievementCodes) > 3 {
				return &recap.ConfigError{Code: "invalid_summary_card_data"}
			}
			if _, ok := registry.Roles[typed.Data.RoleCode]; !ok {
				return &recap.ConfigError{Code: "invalid_summary_role"}
			}
			if _, ok := registry.Styles[typed.Data.StyleCode]; !ok {
				return &recap.ConfigError{Code: "invalid_summary_style"}
			}
		case recap.FinalCard:
			if typed.Data.ShowShareButton != capabilities.ShareAvailable || typed.Data.ShowFeedback != capabilities.FeedbackAvailable {
				return &recap.ConfigError{Code: "final_capability_mismatch"}
			}
		default:
			return &recap.ConfigError{Code: "unknown_card_concrete_type", Message: base.ID}
		}
	}
	return nil
}

func validateMetrics(metrics []recap.MetricValue, registry recap.Registry) error {
	if len(metrics) != len(registry.Metrics) {
		return &recap.ConfigError{Code: "invalid_public_metric_count", Message: fmt.Sprint(len(metrics))}
	}
	seen := map[recap.MetricCode]struct{}{}
	for _, metric := range metrics {
		definition, ok := registry.Metrics[metric.MetricCode]
		if !ok || definition.Unit != metric.Unit {
			return &recap.ConfigError{Code: "invalid_public_metric", Message: string(metric.MetricCode)}
		}
		if _, duplicate := seen[metric.MetricCode]; duplicate {
			return &recap.ConfigError{Code: "duplicate_public_metric", Message: string(metric.MetricCode)}
		}
		seen[metric.MetricCode] = struct{}{}
	}
	return nil
}

func validateAchievements(values []recap.AchievementDecision) error {
	if len(values) < 1 || len(values) > 3 {
		return &recap.ConfigError{Code: "invalid_achievement_count", Message: fmt.Sprint(len(values))}
	}
	groups := map[string]struct{}{}
	for _, achievement := range values {
		if achievement.Value < achievement.CurrentThreshold || achievement.Group == "" {
			return &recap.ConfigError{Code: "achievement_below_threshold", Message: string(achievement.Code)}
		}
		if _, exists := groups[achievement.Group]; exists {
			return &recap.ConfigError{Code: "duplicate_achievement_group", Message: achievement.Group}
		}
		groups[achievement.Group] = struct{}{}
	}
	return nil
}

func validateExplanation(explanation recap.RecapExplanation, output recap.GenerateOutput, registry recap.Registry) error {
	if explanation.RecapID != output.Recap.ID || explanation.AlgorithmVersion != output.Recap.Generation.AlgorithmVersion || explanation.ActivityHash != output.Recap.Generation.ActivityHash {
		return &recap.ConfigError{Code: "explanation_metadata_mismatch"}
	}
	expected := 2 + len(output.Achievements)
	if len(explanation.Decisions) != expected {
		return &recap.ConfigError{Code: "missing_explanations", Message: fmt.Sprintf("got=%d expected=%d", len(explanation.Decisions), expected)}
	}
	for _, decision := range explanation.Decisions {
		if !validText(decision.CardID, 1, 100) || !validText(decision.Reason, 1, 500) || !validText(decision.RuleVersion, 1, 50) || len(decision.Facts) < 1 {
			return &recap.ConfigError{Code: "invalid_explanation_decision", Message: decision.Code}
		}
		if !validDecisionCode(decision, registry) {
			return &recap.ConfigError{Code: "invalid_explanation_code", Message: decision.Code}
		}
		for _, fact := range decision.Facts {
			if _, ok := registry.Metrics[fact.MetricCode]; !ok || !validOperator(fact.Operator) || !fact.Matched {
				return &recap.ConfigError{Code: "invalid_explanation_fact", Message: string(fact.MetricCode)}
			}
		}
	}
	return nil
}

func validateShare(share recap.ShareCard, output recap.GenerateOutput, registry recap.Registry) error {
	if share.SchemaVersion != recap.SchemaVersion || share.RecapID != output.Recap.ID || !uuidPattern.MatchString(share.RecapID) ||
		share.ProfileName != output.Recap.Profile.Name || share.Year != output.Recap.Year || !validText(share.ProfileName, 1, 100) ||
		!validText(share.Title, 1, 120) || !validText(share.Subtitle, 1, 200) {
		return &recap.ConfigError{Code: "invalid_share_header"}
	}
	if share.MainDistrict != output.Recap.Theme.MainDistrict || len(share.Facts) < 1 || len(share.Facts) > 5 || len(share.Achievements) < 1 || len(share.Achievements) > 3 {
		return &recap.ConfigError{Code: "invalid_share_content"}
	}
	if share.AvatarURL != nil && !validPublicAvatar(*share.AvatarURL, output.Recap.ProfileID, registry.PublicAvatarHosts) {
		return &recap.ConfigError{Code: "unsafe_share_avatar"}
	}
	for _, fact := range share.Facts {
		if !validShareFactKind(fact.Kind) || !validText(fact.Label, 1, 100) || !validText(fact.Value, 1, 200) {
			return &recap.ConfigError{Code: "invalid_share_fact"}
		}
	}
	for _, achievement := range share.Achievements {
		if !knownAchievement(achievement.Code, registry) || !validText(achievement.Title, 1, 100) || !validText(achievement.Icon, 1, 100) || !validAchievementLevel(achievement.Level) {
			return &recap.ConfigError{Code: "invalid_share_achievement", Message: string(achievement.Code)}
		}
	}
	if share.Visual.Theme != "city" {
		return &recap.ConfigError{Code: "invalid_share_visual"}
	}
	return nil
}

func validPublicAchievement(value recap.Achievement, registry recap.Registry) bool {
	return knownAchievement(value.Code, registry) &&
		validText(value.Title, 1, 100) &&
		validText(value.Description, 1, 300) &&
		validText(value.Icon, 1, 100) &&
		validAchievementLevel(value.Level) &&
		value.CurrentValue != nil &&
		value.MetricCode != ""
}

func knownAchievement(code recap.AchievementCode, registry recap.Registry) bool {
	for _, definition := range registry.Achievements {
		if definition.Code == code {
			return true
		}
	}
	return false
}

func validDecisionCode(value recap.DecisionExplanation, registry recap.Registry) bool {
	switch value.Kind {
	case recap.ExplanationArchetypeRole:
		_, ok := registry.Roles[recap.ArchetypeRoleCode(value.Code)]
		return ok
	case recap.ExplanationArchetypeStyle:
		_, ok := registry.Styles[recap.ArchetypeStyleCode(value.Code)]
		return ok
	case recap.ExplanationAchievement:
		return knownAchievement(recap.AchievementCode(value.Code), registry)
	default:
		return false
	}
}

func validPublicAvatar(raw, profileID string, allowedHosts map[string]struct{}) bool {
	parsed, err := url.ParseRequestURI(raw)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return false
	}
	if parsed.Port() != "" && parsed.Port() != "443" {
		return false
	}
	if profileID != "" && strings.Contains(strings.ToLower(parsed.EscapedPath()), strings.ToLower(profileID)) {
		return false
	}
	if uuidPattern.FindString(parsed.EscapedPath()) != "" {
		return false
	}
	host := strings.ToLower(parsed.Hostname())
	if host == "" || net.ParseIP(host) != nil || host == "localhost" {
		return false
	}
	for allowed := range allowedHosts {
		if host == allowed || strings.HasSuffix(host, "."+allowed) {
			return true
		}
	}
	return false
}

func validVisibility(value recap.CardVisibility) bool {
	return value == recap.VisibilityPersonal || value == recap.VisibilityShareable
}

func validVisualKind(value recap.CardVisualKind) bool {
	switch value {
	case recap.VisualIllustration, recap.VisualDistrict, recap.VisualStreet, recap.VisualCalendar, recap.VisualBadge, recap.VisualChart, recap.VisualCharacter, recap.VisualSkyline:
		return true
	default:
		return false
	}
}

func validOperator(value recap.RuleOperator) bool {
	switch value {
	case recap.OperatorEQ, recap.OperatorNEQ, recap.OperatorGT, recap.OperatorGTE, recap.OperatorLT, recap.OperatorLTE:
		return true
	default:
		return false
	}
}

func validShareFactKind(value recap.ShareFactKind) bool {
	return value == recap.ShareFactMainDistrict || value == recap.ShareFactActiveDays || value == recap.ShareFactTopAchievement
}

func validAchievementLevel(value recap.AchievementLevel) bool {
	return value == recap.LevelNewcomer || value == recap.LevelLocal || value == recap.LevelExpert || value == recap.LevelGuru
}

func schemaVersionValid(value string) bool {
	parts := strings.Split(value, ".")
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return false
	}
	for _, part := range parts {
		for _, r := range part {
			if r < '0' || r > '9' {
				return false
			}
		}
	}
	return true
}

func validActivityHash(value string) bool {
	const prefix = "sha256:"
	if len(value) != len(prefix)+64 || !strings.HasPrefix(value, prefix) {
		return false
	}
	hexPart := value[len(prefix):]
	if hexPart != strings.ToLower(hexPart) {
		return false
	}
	decoded, err := hex.DecodeString(hexPart)
	return err == nil && len(decoded) == 32
}

func validText(value string, min, max int) bool {
	count := utf8.RuneCountInString(strings.TrimSpace(value))
	return count >= min && count <= max
}

func sameOptionalString(a, b *string) bool {
	if a == nil || b == nil {
		return a == nil && b == nil
	}
	return *a == *b
}

func hexColorValid(value string) bool {
	if len(value) != 7 || value[0] != '#' {
		return false
	}
	_, err := hex.DecodeString(value[1:])
	return err == nil
}

func validateCategoryTaxonomy(geography recap.Geography) error {
	if geography.TopCategoryCode == "" {
		return nil
	}
	vertical, ok := analytics.VerticalForCategory(geography.TopCategoryCode)
	if !ok || vertical != geography.MainVerticalCode {
		return &recap.ConfigError{Code: "category_vertical_mismatch", Message: geography.TopCategoryCode}
	}
	return nil
}
