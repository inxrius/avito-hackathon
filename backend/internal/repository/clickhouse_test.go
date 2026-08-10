package repository

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"recap-personalization/internal/recap/ports"
	"recap-personalization/pkg/database"
)

func TestGetActivitiesByProfileIDAndYearDecodesStringMillis(t *testing.T) {
	var query string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		query = string(body)
		_, _ = io.WriteString(w, `{"event_id":"event-1","profile_id":"profile-1","event_type":"listing_viewed","vertical_code":"goods","category_code":"electronics","occurred_at_ms":"1767225600123"}`+"\n")
	}))
	defer server.Close()

	repo := NewClickHouseRepository(&database.ClickHouseHTTP{Endpoint: server.URL, Database: "recap", Client: server.Client()})
	events, err := repo.GetActivitiesByProfileIDAndYear(context.Background(), "profile-1", 2026)
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("events=%d", len(events))
	}
	if events[0].OccurredAt.UnixMilli() != 1767225600123 {
		t.Fatalf("occurred_at=%d", events[0].OccurredAt.UnixMilli())
	}
	if !strings.Contains(query, "WHERE toString(profile_id) = 'profile-1'") {
		t.Fatalf("profile filter missing: %s", query)
	}
	if !strings.Contains(query, "toString(toUnixTimestamp64Milli(occurred_at)) AS occurred_at_ms") {
		t.Fatalf("string millis conversion missing: %s", query)
	}
}

func TestGetInteractionsByRecapIDUsesRecapFilterAndDecodesProperties(t *testing.T) {
	var query string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		query = string(body)
		_, _ = io.WriteString(w, `{"event_id":"event-1","recap_id":"recap-1","session_id":"session-1","event_name":"recap_opened","occurred_at_ms":"1767225600123","properties":"{\"source\":\"test\"}"}`+"\n")
	}))
	defer server.Close()

	repo := NewClickHouseRepository(&database.ClickHouseHTTP{Endpoint: server.URL, Database: "recap", Client: server.Client()})
	events, err := repo.GetInteractionsByRecapID(context.Background(), "recap-1")
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 {
		t.Fatalf("events=%d", len(events))
	}
	if events[0].RecapID != "recap-1" {
		t.Fatalf("recap_id=%q", events[0].RecapID)
	}
	if events[0].OccurredAt.UnixMilli() != 1767225600123 {
		t.Fatalf("occurred_at=%d", events[0].OccurredAt.UnixMilli())
	}
	if events[0].Properties["source"] != "test" {
		t.Fatalf("properties=%#v", events[0].Properties)
	}
	if !strings.Contains(query, "WHERE toString(recap_id) = 'recap-1'") {
		t.Fatalf("recap filter missing: %s", query)
	}
}

func TestSaveInteractionSkipsExistingEvent(t *testing.T) {
	calls := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		body, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(body), "SELECT count() AS count") {
			t.Fatalf("unexpected query: %s", body)
		}
		_ = json.NewEncoder(w).Encode(map[string]string{"count": "1"})
	}))
	defer server.Close()

	repo := NewClickHouseRepository(&database.ClickHouseHTTP{Endpoint: server.URL, Database: "recap", Client: server.Client()})
	err := repo.SaveInteraction(context.Background(), ports.InteractionEvent{
		EventID:    "event-1",
		RecapID:    "recap-1",
		SessionID:  "session-1",
		EventName:  "recap_opened",
		OccurredAt: time.Date(2026, time.January, 1, 0, 0, 0, 0, time.UTC),
		Properties: map[string]interface{}{"source": "test"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if calls != 1 {
		t.Fatalf("calls=%d", calls)
	}
}
