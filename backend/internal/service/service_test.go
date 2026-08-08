package service

import (
	"testing"

	recap "recap-personalization/internal/recap"
)

func TestConvertEventType(t *testing.T) {
	tests := []struct {
		input    string
		expected recap.EventType
	}{
		{input: "listing_viewed", expected: recap.EventListingViewed},
		{input: "favorite_added", expected: recap.EventFavoriteAdded},
		{input: "search_saved", expected: recap.EventSearchSaved},
		{input: "chat_started", expected: recap.EventChatStarted},
		{input: "listing_published", expected: recap.EventListingPublished},
		{input: "sale_completed", expected: recap.EventSaleCompleted},
		{input: "purchase_completed", expected: recap.EventPurchaseCompleted},
		{input: "delivery_used", expected: recap.EventDeliveryUsed},
	}

	for _, test := range tests {
		t.Run(test.input, func(t *testing.T) {
			actual, err := convertEventType(test.input)
			if err != nil {
				t.Fatalf("convertEventType returned error: %v", err)
			}
			if actual != test.expected {
				t.Fatalf("expected %q, got %q", test.expected, actual)
			}
		})
	}
}

func TestConvertEventTypeRejectsUnknownValue(t *testing.T) {
	if _, err := convertEventType("unknown"); err == nil {
		t.Fatal("expected unsupported event type error")
	}
}

func TestContainsYear(t *testing.T) {
	if !containsYear([]int{2025, 2026}, 2026) {
		t.Fatal("expected year to be available")
	}
	if containsYear([]int{2025, 2026}, 2024) {
		t.Fatal("expected year to be unavailable")
	}
}
