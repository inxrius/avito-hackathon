package pipeline

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	recap "recap-personalization/internal/recap"
	"recap-personalization/internal/recap/narrative"
)

const (
	fixtureProfileID = "19369d67-1db3-4e87-b8a8-9c82991ba173"
	fixtureRecapID   = "550e8400-e29b-41d4-a716-446655440000"
)

type eventBuilder struct {
	id      int
	profile string
	year    int
	events  []recap.ActivityEvent
	dateFn  func(int) time.Time
}

func newBuilder(profile string, year int, dateFn func(int) time.Time) *eventBuilder {
	if dateFn == nil {
		dateFn = func(i int) time.Time {
			return time.Date(year, 1, 1+(i%20), 12, 0, 0, 0, time.UTC)
		}
	}
	return &eventBuilder{profile: profile, year: year, dateFn: dateFn}
}

func (b *eventBuilder) add(eventType recap.EventType, count int, categories ...string) {
	for i := 0; i < count; i++ {
		category := ""
		if len(categories) > 0 {
			category = categories[i%len(categories)]
		}
		b.id++
		b.events = append(b.events, recap.ActivityEvent{
			EventID: fmt.Sprintf("event-%04d", b.id), ProfileID: b.profile, EventType: eventType,
			CategoryCode: category, OccurredAt: b.dateFn(b.id - 1),
		})
	}
}

func fixtureInput(name string) recap.GenerateInput {
	const year = 2026
	var b *eventBuilder
	switch name {
	case "thoughtful_findings_seeker":
		b = newBuilder(fixtureProfileID, year, nil)
		b.add(recap.EventListingViewed, 30, "electronics", "cars")
		b.add(recap.EventFavoriteAdded, 10, "electronics", "cars")
		b.add(recap.EventPurchaseCompleted, 2, "electronics", "cars")
	case "productive_showcase_owner":
		b = newBuilder(fixtureProfileID, year, nil)
		b.add(recap.EventListingPublished, 10, "electronics", "cars")
		b.add(recap.EventSaleCompleted, 4, "electronics", "cars")
	case "universal_city_local":
		dateFn := func(i int) time.Time {
			month := time.Month(1 + i%3)
			day := 1 + (i/3)%20
			return time.Date(year, month, day, 12, 0, 0, 0, time.UTC)
		}
		b = newBuilder(fixtureProfileID, year, dateFn)
		cats := []string{"electronics", "cars", "apartments"}
		b.add(recap.EventListingViewed, 50, cats...)
		b.add(recap.EventFavoriteAdded, 6, cats...)
		b.add(recap.EventSearchSaved, 5, cats...)
		b.add(recap.EventChatStarted, 10, cats...)
		b.add(recap.EventPurchaseCompleted, 2, cats...)
		b.add(recap.EventListingPublished, 7, cats...)
		b.add(recap.EventSaleCompleted, 2, cats...)
	case "city_observer":
		b = newBuilder(fixtureProfileID, year, nil)
		b.add(recap.EventListingViewed, 30, "electronics", "cars")
		b.add(recap.EventFavoriteAdded, 8, "electronics", "cars")
	case "district_expert":
		b = newBuilder(fixtureProfileID, year, nil)
		b.add(recap.EventListingViewed, 20, "electronics")
	case "city_explorer":
		b = newBuilder(fixtureProfileID, year, nil)
		b.add(recap.EventListingViewed, 12, "electronics", "home_and_garden", "cars", "apartments", "vacancies", "personal_services")
	case "delivery_only":
		dateFn := func(i int) time.Time {
			return time.Date(year, 1, 1+(i%7), 12, 0, 0, 0, time.UTC)
		}
		b = newBuilder(fixtureProfileID, year, dateFn)
		b.add(recap.EventDeliveryUsed, 10, "electronics")
	case "insufficient_activity":
		b = newBuilder(fixtureProfileID, year, nil)
		b.add(recap.EventListingViewed, 9, "electronics")
	default:
		panic("unknown fixture " + name)
	}
	return recap.GenerateInput{
		RecapID: fixtureRecapID,
		Profile: recap.ProfileSnapshot{ID: fixtureProfileID, Name: "Тестовый пользователь", AvatarURL: "https://cdn.example/avatar.png"},
		Year:    year, Activities: b.events, GeneratedAt: time.Date(2026, 8, 6, 10, 0, 0, 0, time.UTC),
	}
}

func TestAcceptanceFixtures(t *testing.T) {
	tests := []struct {
		fixture     string
		role        recap.ArchetypeRoleCode
		style       recap.ArchetypeStyleCode
		achievement recap.AchievementCode
	}{
		{"thoughtful_findings_seeker", recap.RoleFindingsSeeker, recap.StyleThoughtful, recap.AchievementFindingsCollector},
		{"productive_showcase_owner", recap.RoleShowcaseOwner, recap.StyleResultOriented, recap.AchievementDealMaster},
		{"universal_city_local", recap.RoleUniversalCitizen, recap.StyleCityLocal, ""},
		{"city_observer", recap.RoleCityObserver, recap.StyleThoughtful, recap.AchievementFindingsCollector},
		{"district_expert", recap.RoleFindingsSeeker, recap.StyleDistrictExpert, ""},
		{"city_explorer", recap.RoleFindingsSeeker, recap.StyleExplorer, recap.AchievementCityNavigator},
	}
	generator := NewGenerator(nil)
	for _, tt := range tests {
		t.Run(tt.fixture, func(t *testing.T) {
			got, err := generator.Generate(context.Background(), fixtureInput(tt.fixture))
			if err != nil {
				t.Fatalf("Generate() error = %v", err)
			}
			if got.Archetype.Role.Code != tt.role || got.Archetype.Style.Code != tt.style {
				t.Fatalf("archetype=%s/%s, want %s/%s", got.Archetype.Role.Code, got.Archetype.Style.Code, tt.role, tt.style)
			}
			if len(got.Recap.Cards) != 8 {
				t.Fatalf("cards=%d", len(got.Recap.Cards))
			}
			if len(got.Achievements) < 1 || len(got.Achievements) > 3 {
				t.Fatalf("achievements=%d", len(got.Achievements))
			}
			if tt.achievement != "" && !hasAchievement(got.Achievements, tt.achievement) {
				t.Fatalf("missing achievement %q: %+v", tt.achievement, got.Achievements)
			}
			if got.Narrative.Source != "template" || got.Recap.Generation.Narrative.Source != "template" {
				t.Fatalf("narrative source=%q", got.Narrative.Source)
			}
			if got.Narrative.Model != nil || got.Recap.Generation.Narrative.Model != nil {
				t.Fatalf("template model must be nil")
			}
			if got.Recap.ProfileID != fixtureProfileID || got.Recap.Profile.ID != fixtureProfileID {
				t.Fatalf("OpenAPI-required profile IDs are missing")
			}
			if got.Explanation == nil || len(got.Explanation.Decisions) != 2+len(got.Achievements) {
				t.Fatalf("explanations=%v", got.Explanation)
			}
			if got.Share == nil || len(got.Share.Facts) != 3 {
				t.Fatalf("share=%+v", got.Share)
			}
			encodedShare, _ := json.Marshal(got.Share)
			if strings.Contains(string(encodedShare), fixtureProfileID) || strings.Contains(string(encodedShare), `"profile_id"`) {
				t.Fatalf("share leaks profile id: %s", encodedShare)
			}
		})
	}
}

func TestInsufficientActivitySkipsNarrativeAndReturnsDomainError(t *testing.T) {
	provider := &countingProvider{err: errors.New("must not be called")}
	generator := NewGenerator(provider)
	_, err := generator.Generate(context.Background(), fixtureInput("insufficient_activity"))
	if !errors.Is(err, recap.ErrInsufficientActivity) {
		t.Fatalf("error=%v", err)
	}
	if provider.calls != 0 {
		t.Fatalf("provider calls=%d, want 0", provider.calls)
	}
}

type countingProvider struct {
	content []byte
	model   string
	err     error
	calls   int
}

func (p *countingProvider) Generate(context.Context, narrative.Input) (narrative.ProviderResult, error) {
	p.calls++
	return narrative.ProviderResult{Content: p.content, Model: p.model}, p.err
}

func TestMistralSuccessAndFallbackAreSingleCall(t *testing.T) {
	input := fixtureInput("thoughtful_findings_seeker")
	provider := &countingProvider{content: []byte(`{"summary_title":"Городской маршрут","summary_text":"Ты был активен 20 дней.","unexpected":true}`), model: "model-x"}
	generator := NewGenerator(provider)
	got, err := generator.Generate(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if provider.calls != 1 || got.Narrative.Source != "template" {
		t.Fatalf("calls=%d narrative=%+v", provider.calls, got.Narrative)
	}

	provider.content = []byte(`{"summary_title":"Итоги 2026","summary_text":"Ты был активен 20 дней."}`)
	provider.calls = 0
	got, err = generator.Generate(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if provider.calls != 1 || got.Narrative.Source != "mistral" || got.Narrative.Model == nil || *got.Narrative.Model != "model-x" {
		t.Fatalf("narrative=%+v calls=%d", got.Narrative, provider.calls)
	}
}

func TestRegistryFailureIsDetectedBeforeNarrativeCall(t *testing.T) {
	provider := &countingProvider{content: []byte(`{"summary_title":"x","summary_text":"y"}`), model: "m"}
	generator := NewGenerator(provider)
	delete(generator.Registry.Categories, recap.CategoryElectronics)
	_, err := generator.Generate(context.Background(), fixtureInput("district_expert"))
	if err == nil {
		t.Fatal("expected registry error")
	}
	if provider.calls != 0 {
		t.Fatalf("provider calls=%d", provider.calls)
	}
}

func TestDeliveryOnlyFallbackUsesNeutralMetricAndExplanation(t *testing.T) {
	got, err := NewGenerator(nil).Generate(context.Background(), fixtureInput("delivery_only"))
	if err != nil {
		t.Fatal(err)
	}
	if got.Archetype.Role.Code != recap.RoleFindingsSeeker {
		t.Fatalf("role=%s", got.Archetype.Role.Code)
	}
	metric, ok := got.Recap.Cards[2].(recap.MetricCard)
	if !ok {
		t.Fatalf("card type=%T", got.Recap.Cards[2])
	}
	if metric.Data.MetricCode != recap.MetricActiveDays || metric.Data.Value <= 0 || strings.Contains(metric.Title, "0 ") {
		t.Fatalf("metric=%+v", metric)
	}
	if metric.Description == nil || strings.Contains(*metric.Description, "исследовал") {
		t.Fatalf("description=%v", metric.Description)
	}
	if got.Explanation == nil || !strings.Contains(got.Explanation.Decisions[0].Reason, "fallback") {
		t.Fatalf("role explanation=%+v", got.Explanation)
	}
}

func TestAchievementsDescriptionIsCountNeutral(t *testing.T) {
	got, err := NewGenerator(nil).Generate(context.Background(), fixtureInput("district_expert"))
	if err != nil {
		t.Fatal(err)
	}
	card, ok := got.Recap.Cards[5].(recap.AchievementsCard)
	if !ok || card.Description == nil {
		t.Fatalf("card=%T", got.Recap.Cards[5])
	}
	if *card.Description != "Достижения, которые лучше всего описывают твой год" {
		t.Fatalf("description=%q", *card.Description)
	}
}

func TestShareIsNotBuiltWhenCapabilityIsDisabled(t *testing.T) {
	generator := NewGenerator(nil)
	generator.Capabilities.ShareAvailable = false
	got, err := generator.Generate(context.Background(), fixtureInput("district_expert"))
	if err != nil {
		t.Fatal(err)
	}
	if got.Share != nil {
		t.Fatalf("share must be nil: %+v", got.Share)
	}
	final := got.Recap.Cards[7].(recap.FinalCard)
	if final.Data.ShowShareButton {
		t.Fatal("share button must be hidden")
	}
}

func TestShareAvatarRequiresAllowlistedPublicCDNURL(t *testing.T) {
	generator := NewGenerator(nil)
	generator.Registry.PublicAvatarHosts["cdn.example"] = struct{}{}
	input := fixtureInput("district_expert")
	got, err := generator.Generate(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if got.Share == nil || got.Share.AvatarURL == nil || *got.Share.AvatarURL != input.Profile.AvatarURL {
		t.Fatalf("avatar=%v", got.Share)
	}

	input.Profile.AvatarURL = "https://cdn.example/avatar.png?token=secret"
	got, err = generator.Generate(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if got.Share == nil || got.Share.AvatarURL != nil {
		t.Fatalf("signed avatar leaked: %+v", got.Share)
	}
	if got.Recap.Profile.AvatarURL == nil {
		t.Fatal("personal recap may retain the profile snapshot avatar")
	}

	input.Profile.AvatarURL = "https://cdn.example/avatars/" + fixtureProfileID + ".png"
	got, err = generator.Generate(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if got.Share == nil || got.Share.AvatarURL != nil {
		t.Fatalf("identifier-bearing avatar leaked: %+v", got.Share)
	}
}

func TestSameInputProducesSameDeterministicOutput(t *testing.T) {
	generator := NewGenerator(nil)
	input := fixtureInput("productive_showcase_owner")
	first, err := generator.Generate(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	second, err := generator.Generate(context.Background(), input)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(first, second) {
		t.Fatal("outputs differ for identical input")
	}
}

func TestValidActivityHashRequiresLowercaseSHA256(t *testing.T) {
	valid := "sha256:" + strings.Repeat("a", 64)
	if !validActivityHash(valid) {
		t.Fatal("valid hash rejected")
	}
	if validActivityHash("sha256:" + strings.Repeat("A", 64)) {
		t.Fatal("uppercase hash accepted")
	}
	if validActivityHash("sha256:" + strings.Repeat("z", 64)) {
		t.Fatal("non-hex hash accepted")
	}
}

func hasAchievement(values []recap.AchievementDecision, code recap.AchievementCode) bool {
	for _, value := range values {
		if value.Code == code {
			return true
		}
	}
	return false
}
