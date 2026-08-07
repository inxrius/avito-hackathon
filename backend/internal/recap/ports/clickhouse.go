package ports

// ActivityRepository — порт для работы с активностями в ClickHouse
type ActivityRepository interface {
	GetActivitiesByProfileIDAndYear(profileID string, year int) ([]ActivityEvent, error)
}

// InteractionRepository — порт для работы с взаимодействиями в ClickHouse
type InteractionRepository interface {
	SaveInteraction(event InteractionEvent) error
	GetInteractionsByRecapID(recapID string) ([]InteractionEvent, error)
}

// ActivityEvent — событие активности (для ClickHouse)
type ActivityEvent struct {
	EventID      string
	ProfileID    string
	EventType    string
	CategoryCode string
	OccurredAt   int64 // Unix timestamp
}

// InteractionEvent — событие взаимодействия (для ClickHouse)
type InteractionEvent struct {
	EventID    string
	RecapID    string
	SessionID  string
	EventName  string
	OccurredAt int64 // Unix timestamp
	Properties map[string]interface{}
}