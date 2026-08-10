package assembly

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"strings"

	recap "recap-personalization/internal/recap"
)

var uuidInAvatarPathPattern = regexp.MustCompile(`(?i)[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`)

func BuildShareCard(recapID string, profile recap.ProfileSnapshot, year int, archetype recap.ArchetypeDecision, metrics recap.Metrics, achievements []recap.AchievementDecision, theme recap.RecapTheme, registry recap.Registry) (recap.ShareCard, error) {
	facts := []recap.ShareFact{
		{Kind: recap.ShareFactMainDistrict, Label: "Главный район", Value: theme.MainDistrict.Title},
		{Kind: recap.ShareFactActiveDays, Label: "Дней в ритме города", Value: fmt.Sprint(metrics.ActiveDays)},
	}
	if len(achievements) > 0 {
		facts = append(facts, recap.ShareFact{Kind: recap.ShareFactTopAchievement, Label: "Главное городское звание", Value: achievements[0].Title})
	}
	shareAchievements := make([]recap.ShareAchievement, 0, len(achievements))
	for _, achievement := range achievements {
		shareAchievements = append(shareAchievements, recap.ShareAchievement{Code: achievement.Code, Title: achievement.Title, Level: achievement.Level, Icon: achievement.Icon})
	}
	return recap.ShareCard{
		SchemaVersion: recap.SchemaVersion,
		RecapID:       recapID,
		ProfileName:   profile.Name,
		AvatarURL:     publicAvatarURL(profile.AvatarURL, profile.ID, registry.PublicAvatarHosts),
		Year:          year,
		Title:         "Мой год в городе Авито",
		Subtitle:      fmt.Sprintf("%s · %s", archetype.Style.Title, archetype.Role.Title),
		MainDistrict:  theme.MainDistrict,
		Facts:         facts,
		Achievements:  shareAchievements,
		Visual:        recap.ShareVisual{Theme: "city"},
	}, nil
}

func publicAvatarURL(raw, profileID string, allowedHosts map[string]struct{}) *string {
	if strings.TrimSpace(raw) == "" || len(allowedHosts) == 0 {
		return nil
	}
	parsed, err := url.ParseRequestURI(raw)
	if err != nil || parsed.Scheme != "https" || parsed.Host == "" || parsed.User != nil || parsed.RawQuery != "" || parsed.Fragment != "" {
		return nil
	}
	if parsed.Port() != "" && parsed.Port() != "443" {
		return nil
	}
	if profileID != "" && strings.Contains(strings.ToLower(parsed.EscapedPath()), strings.ToLower(profileID)) {
		return nil
	}
	if uuidInAvatarPathPattern.MatchString(parsed.EscapedPath()) {
		return nil
	}
	host := strings.ToLower(parsed.Hostname())
	if host == "" || net.ParseIP(host) != nil || host == "localhost" {
		return nil
	}
	for allowed := range allowedHosts {
		if host == allowed || strings.HasSuffix(host, "."+allowed) {
			value := parsed.String()
			return &value
		}
	}
	return nil
}
