package ports

import (
	"context"
	"time"
)

type ActivityRepository interface {
	GetActivitiesByProfileIDAndYear(ctx context.Context, profileID string, year int) ([]ActivityEvent, error)
}

type InteractionRepository interface {
	SaveInteraction(ctx context.Context, event InteractionEvent) error
	GetInteractionsByRecapID(ctx context.Context, recapID string) ([]InteractionEvent, error)
}

type ActivityEvent struct {
	EventID      string
	ProfileID    string
	EventType    string
	VerticalCode string
	CategoryCode string
	OccurredAt   time.Time
}

type InteractionEvent struct {
	EventID    string
	RecapID    string
	SessionID  string
	EventName  string
	OccurredAt time.Time
	Properties map[string]interface{}
}
