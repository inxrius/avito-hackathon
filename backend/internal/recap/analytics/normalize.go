package analytics

import (
	"fmt"
	"sort"
	"time"

	recap "recap-personalization/internal/recap"
)

var supportedEvents = map[recap.EventType]struct{}{
	recap.EventListingViewed: {}, recap.EventFavoriteAdded: {}, recap.EventSearchSaved: {},
	recap.EventChatStarted: {}, recap.EventListingPublished: {}, recap.EventSaleCompleted: {},
	recap.EventPurchaseCompleted: {}, recap.EventDeliveryUsed: {},
}

type rawSignificantFields struct {
	EventType    recap.EventType
	VerticalCode string
	CategoryCode string
	OccurredAt   time.Time
}

func Normalize(profileID string, year int, events []recap.ActivityEvent) ([]recap.NormalizedEvent, error) {
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	end := start.AddDate(1, 0, 0)
	byID := make(map[string]recap.NormalizedEvent)
	rawByID := make(map[string]rawSignificantFields)

	for _, event := range events {
		occurredAt := event.OccurredAt.UTC()
		if event.ProfileID != profileID || occurredAt.Before(start) || !occurredAt.Before(end) || event.EventID == "" {
			continue
		}
		if _, ok := supportedEvents[event.EventType]; !ok {
			return nil, &recap.InputError{Code: "unknown_event_type", Message: string(event.EventType)}
		}

		raw := rawSignificantFields{
			EventType: event.EventType, VerticalCode: event.VerticalCode,
			CategoryCode: event.CategoryCode, OccurredAt: occurredAt,
		}
		if previous, exists := rawByID[event.EventID]; exists {
			if !sameRawSignificantFields(previous, raw) {
				return nil, &recap.InputError{Code: "conflicting_event_duplicate", Message: fmt.Sprintf("event_id=%s", event.EventID)}
			}
			continue
		}

		vertical := ""
		category := ""
		if event.CategoryCode != "" {
			resolved, ok := VerticalForCategory(event.CategoryCode)
			if !ok {
				return nil, &recap.InputError{Code: "unknown_category_code", Message: event.CategoryCode}
			}
			category = event.CategoryCode
			vertical = resolved
		} else if event.VerticalCode != "" {
			if !IsAllowedVertical(event.VerticalCode) {
				return nil, &recap.InputError{Code: "unknown_vertical_code", Message: event.VerticalCode}
			}
			vertical = event.VerticalCode
		}

		rawByID[event.EventID] = raw
		byID[event.EventID] = recap.NormalizedEvent{
			EventID: event.EventID, EventType: event.EventType,
			ResolvedVerticalCode: vertical, ResolvedCategoryCode: category, OccurredAt: occurredAt,
		}
	}

	result := make([]recap.NormalizedEvent, 0, len(byID))
	for _, event := range byID {
		result = append(result, event)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].EventID < result[j].EventID })
	return result, nil
}

func sameRawSignificantFields(a, b rawSignificantFields) bool {
	return a.EventType == b.EventType && a.VerticalCode == b.VerticalCode &&
				 a.CategoryCode == b.CategoryCode && a.OccurredAt.Equal(b.OccurredAt)
}
