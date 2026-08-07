package analytics

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"sort"
	"strings"
	"time"

	recap "recap-personalization/internal/recap"
)

type canonicalEvent struct {
	EventID              string `json:"event_id"`
	EventType            string `json:"event_type"`
	ResolvedVerticalCode string `json:"resolved_vertical_code"`
	ResolvedCategoryCode string `json:"resolved_category_code"`
	OccurredAt           string `json:"occurred_at"`
}

func ActivityHash(events []recap.NormalizedEvent) (string, error) {
	ordered := append([]recap.NormalizedEvent(nil), events...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].EventID < ordered[j].EventID })
	lines := make([]string, 0, len(ordered))
	for _, event := range ordered {
		payload := canonicalEvent{
			EventID:              event.EventID,
			EventType:            string(event.EventType),
			ResolvedVerticalCode: event.ResolvedVerticalCode,
			ResolvedCategoryCode: event.ResolvedCategoryCode,
			OccurredAt:           event.OccurredAt.UTC().Format(time.RFC3339Nano),
		}
		encoded, err := json.Marshal(payload)
		if err != nil {
			return "", err
		}
		lines = append(lines, string(encoded))
	}
	sum := sha256.Sum256([]byte(strings.Join(lines, "\n")))
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}
