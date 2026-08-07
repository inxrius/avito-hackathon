package model

import (
	"time"

	"github.com/google/uuid"
)

// Activity — событие активности пользователя (таблица activities)
type Activity struct {
	ID          uuid.UUID `json:"id"`
	ProfileID   uuid.UUID `json:"profile_id"`
	Type        string    `json:"type"`        // "view", "favorite", "purchase", "sale", "chat", "search", "publish", "delivery"
	Category    string    `json:"category"`    // категория (например, "Недвижимость", "Электроника")
	Title       string    `json:"title"`       // заголовок объявления или действия
	Description string    `json:"description"` // описание
	Value       float64   `json:"value"`       // числовое значение (для метрик)
	Timestamp   time.Time `json:"timestamp"`   // время совершения действия
}
