package personalization

import (
	"testing"

	recap "recap-personalization/internal/recap"
)

func TestRoleRulesInPriorityOrder(t *testing.T) {
	tests := []struct {
		name string
		m    recap.Metrics
		want recap.ArchetypeRoleCode
	}{
		{"observer", recap.Metrics{ViewsCount: 30, FavoritesCount: 8, CompletedDealsCount: 1}, recap.RoleCityObserver},
		{"universal", recap.Metrics{ViewsCount: 50, FavoritesCount: 6, SavedSearchesCount: 5, ChatsStartedCount: 10, PurchasesCount: 2, PublishedListingsCount: 7, SalesCount: 2, CompletedDealsCount: 4}, recap.RoleUniversalCitizen},
		{"showcase", recap.Metrics{PublishedListingsCount: 10, SalesCount: 4}, recap.RoleShowcaseOwner},
		{"fallback", recap.Metrics{ViewsCount: 10}, recap.RoleFindingsSeeker},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := SelectArchetype(tt.m, recap.DefaultRegistry())
			if err != nil {
				t.Fatal(err)
			}
			if got.Role.Code != tt.want {
				t.Fatalf("role=%q, want %q", got.Role.Code, tt.want)
			}
		})
	}
}

func TestUniversalCitizenRequiresCompletedDeal(t *testing.T) {
	m := recap.Metrics{
		ViewsCount: 50, FavoritesCount: 20, SavedSearchesCount: 5, ChatsStartedCount: 10,
		PublishedListingsCount: 10, CompletedDealsCount: 0,
	}
	got, _, err := SelectArchetype(m, recap.DefaultRegistry())
	if err != nil {
		t.Fatal(err)
	}
	if got.Role.Code == recap.RoleUniversalCitizen {
		t.Fatalf("zero-deal profile must not become %s", recap.RoleUniversalCitizen)
	}
}

func TestStyleRulesInPriorityOrder(t *testing.T) {
	tests := []struct {
		name string
		m    recap.Metrics
		want recap.ArchetypeStyleCode
	}{
		{"district", recap.Metrics{MeaningfulEvents: 20, TopVerticalShare: .65, ActiveDays: 7}, recap.StyleDistrictExpert},
		{"explorer", recap.Metrics{MeaningfulEvents: 15, UniqueCategories: 6, UniqueVerticals: 3, TopVerticalShare: .4, ActiveDays: 7}, recap.StyleExplorer},
		{"result", recap.Metrics{MeaningfulEvents: 15, CompletedDealsCount: 4, FavoriteToPurchaseRate: .35, ActiveDays: 7}, recap.StyleResultOriented},
		{"thoughtful", recap.Metrics{MeaningfulEvents: 15, ViewsCount: 30, FavoritesCount: 8, FavoriteToPurchaseRate: .2, ActiveDays: 7}, recap.StyleThoughtful},
		{"local", recap.Metrics{MeaningfulEvents: 15, ViewsCount: 50, FavoritesCount: 6, SavedSearchesCount: 5, ChatsStartedCount: 10, PurchasesCount: 2, PublishedListingsCount: 7, SalesCount: 2, UniqueCategories: 3, ActiveMonths: 3, ActiveDays: 7}, recap.StyleCityLocal},
		{"regular", recap.Metrics{MeaningfulEvents: 40, ActiveMonths: 6, ActiveDays: 30, TopVerticalShare: .5}, recap.StyleRegular},
		{"fallback", recap.Metrics{MeaningfulEvents: 10, ActiveDays: 7}, recap.StyleCityLocal},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _, err := SelectArchetype(tt.m, recap.DefaultRegistry())
			if err != nil {
				t.Fatal(err)
			}
			if got.Style.Code != tt.want {
				t.Fatalf("style=%q, want %q", got.Style.Code, tt.want)
			}
		})
	}
}

func TestAchievementsChooseMaximumLevelAndOnePerGroup(t *testing.T) {
	m := recap.Metrics{
		SalesCount: 7, PurchasesCount: 15, FavoritesCount: 40, UniqueCategories: 8,
		ActiveDays: 90, ActiveMonths: 9, DeliveryCount: 7, PublishedListingsCount: 10, MaxActivityStreak: 14,
	}
	all := CalculateAchievements(m, recap.DefaultRegistry())
	if len(all) != 9 {
		t.Fatalf("earned=%d, want 9", len(all))
	}
	top := SelectTopAchievements(all)
	if len(top) != 3 {
		t.Fatalf("top=%d, want 3", len(top))
	}
	groups := map[string]bool{}
	for _, a := range top {
		if groups[a.Group] {
			t.Fatalf("duplicate group %q", a.Group)
		}
		groups[a.Group] = true
	}
	if top[0].Level != recap.LevelGuru || top[0].Code != recap.AchievementFindingsHunter {
		t.Fatalf("first=%s/%s, want findings_hunter/guru", top[0].Code, top[0].Level)
	}
}

func TestSufficiencyThresholdAlwaysEarnsFrequentGuest(t *testing.T) {
	all := CalculateAchievements(recap.Metrics{ActiveDays: 7}, recap.DefaultRegistry())
	if len(all) != 1 || all[0].Code != recap.AchievementFrequentGuest {
		t.Fatalf("achievements=%+v", all)
	}
}
