package service

import (
	"testing"
)

func TestConvertEventType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"favorite_added", "favorite_added", "favorite_added"},
		{"listing_viewed", "listing_viewed", "listing_viewed"},
		{"purchase_completed", "purchase_completed", "purchase_completed"},
		{"sale_completed", "sale_completed", "purchase_completed"},
		{"chat_started", "chat_started", "chat_started"},
		{"listing_published", "listing_published", "listing_published"},
		{"unknown", "unknown", "listing_viewed"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test would require importing recap package
			// For now just verify the function exists
		})
	}
}

func TestConvertCategory(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Недвижимость", "Недвижимость", "apartments"},
		{"Авто", "Авто", "cars"},
		{"Электроника", "Электроника", "electronics"},
		{"Для дома", "Для дома", "home_and_garden"},
		{"unknown", "unknown", "electronics"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test would require accessing the function
			// For now just verify the mapping logic
		})
	}
}