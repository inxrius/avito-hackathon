package analytics

import (
	"errors"
	"strings"
	"testing"
	"time"

	recap "recap-personalization/internal/recap"
)

func TestNormalizeFiltersResolvesTaxonomyAndDeduplicates(t *testing.T) {
	at := time.Date(2026, 2, 3, 10, 0, 0, 0, time.FixedZone("plus3", 3*60*60))
	events := []recap.ActivityEvent{
		{EventID: "a", ProfileID: "p1", EventType: recap.EventListingViewed, VerticalCode: "transport", CategoryCode: "electronics", OccurredAt: at},
		{EventID: "a", ProfileID: "p1", EventType: recap.EventListingViewed, VerticalCode: "transport", CategoryCode: "electronics", OccurredAt: at},
		{EventID: "other-profile", ProfileID: "p2", EventType: recap.EventListingViewed, CategoryCode: "cars", OccurredAt: at},
		{EventID: "other-year", ProfileID: "p1", EventType: recap.EventListingViewed, CategoryCode: "cars", OccurredAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
		{EventID: "", ProfileID: "p1", EventType: recap.EventListingViewed, CategoryCode: "cars", OccurredAt: at},
	}

	got, err := Normalize("p1", 2026, events)
	if err != nil {
		t.Fatalf("Normalize() error = %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("len = %d, want 1", len(got))
	}
	if got[0].ResolvedVerticalCode != "goods" {
		t.Fatalf("vertical = %q, want goods", got[0].ResolvedVerticalCode)
	}
	if got[0].ResolvedCategoryCode != "electronics" {
		t.Fatalf("category = %q", got[0].ResolvedCategoryCode)
	}
	if got[0].OccurredAt.Location() != time.UTC {
		t.Fatalf("location = %v, want UTC", got[0].OccurredAt.Location())
	}
}

func TestNormalizeRejectsConflictingDuplicate(t *testing.T) {
	at := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := Normalize("p1", 2026, []recap.ActivityEvent{
		{EventID: "same", ProfileID: "p1", EventType: recap.EventListingViewed, CategoryCode: "electronics", OccurredAt: at},
		{EventID: "same", ProfileID: "p1", EventType: recap.EventFavoriteAdded, CategoryCode: "electronics", OccurredAt: at},
	})
	var inputErr *recap.InputError
	if !errors.As(err, &inputErr) || inputErr.Code != "conflicting_event_duplicate" {
		t.Fatalf("error = %#v, want conflicting_event_duplicate", err)
	}
}

func TestNormalizeTreatsRawVerticalAsSignificantForDuplicates(t *testing.T) {
	at := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := Normalize("p1", 2026, []recap.ActivityEvent{
		{EventID: "same", ProfileID: "p1", EventType: recap.EventListingViewed, VerticalCode: "goods", CategoryCode: "electronics", OccurredAt: at},
		{EventID: "same", ProfileID: "p1", EventType: recap.EventListingViewed, VerticalCode: "transport", CategoryCode: "electronics", OccurredAt: at},
	})
	var inputErr *recap.InputError
	if !errors.As(err, &inputErr) || inputErr.Code != "conflicting_event_duplicate" {
		t.Fatalf("error = %#v, want conflicting_event_duplicate", err)
	}
}

func TestNormalizeRejectsUnknownRelevantEvent(t *testing.T) {
	_, err := Normalize("p1", 2026, []recap.ActivityEvent{{
		EventID: "x", ProfileID: "p1", EventType: "future_event", CategoryCode: "electronics",
		OccurredAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
	}})
	var inputErr *recap.InputError
	if !errors.As(err, &inputErr) || inputErr.Code != "unknown_event_type" {
		t.Fatalf("error = %v", err)
	}
}

func TestCalculateMetricsUsesUTCAndBoundsRatios(t *testing.T) {
	events := []recap.NormalizedEvent{
		{EventID: "1", EventType: recap.EventFavoriteAdded, ResolvedVerticalCode: "goods", ResolvedCategoryCode: "electronics", OccurredAt: time.Date(2026, 1, 1, 23, 30, 0, 0, time.UTC)},
		{EventID: "2", EventType: recap.EventPurchaseCompleted, ResolvedVerticalCode: "goods", ResolvedCategoryCode: "electronics", OccurredAt: time.Date(2026, 1, 2, 0, 30, 0, 0, time.UTC)},
		{EventID: "3", EventType: recap.EventPurchaseCompleted, ResolvedVerticalCode: "transport", ResolvedCategoryCode: "cars", OccurredAt: time.Date(2026, 1, 3, 0, 0, 0, 0, time.UTC)},
		{EventID: "4", EventType: recap.EventDeliveryUsed, OccurredAt: time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC)},
	}
	m := CalculateMetrics(events)
	if m.ActiveDays != 4 || m.MaxActivityStreak != 3 {
		t.Fatalf("active_days=%d streak=%d", m.ActiveDays, m.MaxActivityStreak)
	}
	if m.FavoriteToPurchaseRate != 1 {
		t.Fatalf("favorite rate=%v, want 1", m.FavoriteToPurchaseRate)
	}
	if m.DeliveryUsageRate != 0.5 {
		t.Fatalf("delivery rate=%v, want .5", m.DeliveryUsageRate)
	}
	if m.UniqueCategories != 2 || m.UniqueVerticals != 2 {
		t.Fatalf("unique geography = %d/%d", m.UniqueCategories, m.UniqueVerticals)
	}
}

func TestCalculateGeographyTieBreaksByDealsThenCode(t *testing.T) {
	events := []recap.NormalizedEvent{
		{EventID: "1", EventType: recap.EventSaleCompleted, ResolvedVerticalCode: "transport", ResolvedCategoryCode: "cars"},
		{EventID: "2", EventType: recap.EventListingViewed, ResolvedVerticalCode: "goods", ResolvedCategoryCode: "electronics"},
		{EventID: "3", EventType: recap.EventListingViewed, ResolvedVerticalCode: "goods", ResolvedCategoryCode: "electronics"},
		{EventID: "4", EventType: recap.EventListingViewed, ResolvedVerticalCode: "goods", ResolvedCategoryCode: "electronics"},
		{EventID: "5", EventType: recap.EventListingViewed, ResolvedVerticalCode: "goods", ResolvedCategoryCode: "electronics"},
		{EventID: "6", EventType: recap.EventListingViewed, ResolvedVerticalCode: "goods", ResolvedCategoryCode: "electronics"},
	}
	m := CalculateMetrics(events)
	geo := CalculateGeography(events, &m)
	if geo.MainVerticalCode != "transport" {
		t.Fatalf("main vertical=%q, want transport due to deal tie-break", geo.MainVerticalCode)
	}
	if m.TopVerticalShare != 0.5 {
		t.Fatalf("share=%v, want .5", m.TopVerticalShare)
	}
}

func TestActivityHashIsOrderIndependentAndCanonical(t *testing.T) {
	a := recap.NormalizedEvent{EventID: "a", EventType: recap.EventListingViewed, ResolvedVerticalCode: "goods", ResolvedCategoryCode: "electronics", OccurredAt: time.Date(2026, 1, 1, 1, 2, 3, 4, time.UTC)}
	b := recap.NormalizedEvent{EventID: "b", EventType: recap.EventFavoriteAdded, ResolvedVerticalCode: "transport", ResolvedCategoryCode: "cars", OccurredAt: time.Date(2026, 2, 1, 1, 2, 3, 4, time.UTC)}
	h1, err := ActivityHash([]recap.NormalizedEvent{a, b})
	if err != nil {
		t.Fatal(err)
	}
	h2, err := ActivityHash([]recap.NormalizedEvent{b, a})
	if err != nil {
		t.Fatal(err)
	}
	if h1 != h2 {
		t.Fatalf("hashes differ: %s %s", h1, h2)
	}
	if !strings.HasPrefix(h1, "sha256:") || len(h1) != 71 {
		t.Fatalf("bad hash: %q", h1)
	}
}
